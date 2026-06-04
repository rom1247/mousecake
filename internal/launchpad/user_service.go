package launchpad

import (
	"context"
	"fmt"
	"math/big"

	"github.com/mousecake-go/mousecake-go/internal/launchpad/domain"
)

// UserService 用户链上操作用例，处理 deposit、harvest、release 操作的资格校验和 calldata 生成。
type UserService struct {
	prepareSvc    *PrepareService
	querySvc      *UserQueryService
	encoder       EncoderInterface
	contract      string
	chain         domain.ChainReader
	saleRepo      domain.SaleRepository
	poolRepo      domain.PoolRepository
	tierLimitRepo domain.TierLimitRepository
	whitelistRepo domain.WhitelistRepository
	userPoolRepo  domain.UserPoolStateRepository
	creditRepo    domain.UserCreditRepository
	// findVestingScheduleByID 根据 ID 查询 VestingSchedule，Release 操作使用。
	findVestingScheduleByID func(ctx context.Context, id int64) (*domain.VestingSchedule, error)
}

// NewUserService 创建 UserService。
func NewUserService(
	prepareSvc *PrepareService,
	querySvc *UserQueryService,
	encoder EncoderInterface,
	contract string,
	chain domain.ChainReader,
	saleRepo domain.SaleRepository,
	poolRepo domain.PoolRepository,
	tierLimitRepo domain.TierLimitRepository,
	whitelistRepo domain.WhitelistRepository,
	userPoolRepo domain.UserPoolStateRepository,
	creditRepo domain.UserCreditRepository,
	vestingScheduleRepo *VestingScheduleRepository,
) *UserService {
	return &UserService{
		prepareSvc:              prepareSvc,
		querySvc:                querySvc,
		encoder:                 encoder,
		contract:                contract,
		chain:                   chain,
		saleRepo:                saleRepo,
		poolRepo:                poolRepo,
		tierLimitRepo:           tierLimitRepo,
		whitelistRepo:           whitelistRepo,
		userPoolRepo:            userPoolRepo,
		creditRepo:              creditRepo,
		findVestingScheduleByID: vestingScheduleRepo.FindByID,
	}
}

// DepositInput 用户申购的输入参数。
type DepositInput struct {
	CallerAddress string `json:"caller_address" binding:"required"`
	SaleID        int64  `json:"sale_id" binding:"required"`
	PoolIndex     int64  `json:"pool_index" binding:"required"`
	Amount        string `json:"amount" binding:"required"`
}

// Deposit 用户申购：校验资格 + 生成 calldata + 创建 Prepare 交易。
func (s *UserService) Deposit(ctx context.Context, input DepositInput) (*domain.PrepareTx, error) {
	sale, err := s.saleRepo.FindByID(ctx, input.SaleID)
	if err != nil {
		return nil, fmt.Errorf("查询 sale %d: %w", input.SaleID, err)
	}
	if sale.Status != domain.SaleDeployed {
		return nil, fmt.Errorf("sale %d 尚未部署: %w", input.SaleID, domain.ErrSaleNotDeployed)
	}

	pool, err := s.poolRepo.FindBySaleAndPool(ctx, input.SaleID, int(input.PoolIndex))
	if err != nil {
		return nil, fmt.Errorf("查询池 %d/%d: %w", input.SaleID, input.PoolIndex, err)
	}

	amount, ok := new(big.Int).SetString(input.Amount, 10)
	if !ok {
		return nil, fmt.Errorf("无效金额: %s", input.Amount)
	}

	// 普通池校验 Tier 额度
	if !pool.IsSpecialSale {
		tierLimit, err := s.tierLimitRepo.FindBySaleAndTier(ctx, input.SaleID, 0)
		if err == nil && tierLimit != nil && tierLimit.CreditLimit != nil && tierLimit.CreditLimit.Sign() > 0 {
			state, _ := s.userPoolRepo.FindByUserAndPool(ctx, input.CallerAddress, input.SaleID, int(input.PoolIndex))
			currentAmount := big.NewInt(0)
			if state != nil && state.TotalDeposited != nil {
				currentAmount = state.TotalDeposited
			}
			if new(big.Int).Add(currentAmount, amount).Cmp(tierLimit.CreditLimit) > 0 {
				return nil, fmt.Errorf("超出 Tier 额度限制: %w", domain.ErrTierLimitExceeded)
			}
		}
	}

	// 特殊池校验白名单
	if pool.IsSpecialSale {
		whitelisted, err := s.whitelistRepo.IsWhitelisted(ctx, input.SaleID, input.CallerAddress)
		if err != nil {
			return nil, fmt.Errorf("检查白名单: %w", err)
		}
		if !whitelisted {
			return nil, fmt.Errorf("用户不在白名单中: %w", domain.ErrNotWhitelisted)
		}
	}

	// 编码 calldata
	calldata, err := s.encoder.EncodeCall(s.contract, "deposit",
		amount, big.NewInt(input.PoolIndex),
	)
	if err != nil {
		return nil, fmt.Errorf("编码 deposit: %w", err)
	}

	saleID := input.SaleID
	poolIdx := input.PoolIndex
	return s.prepareSvc.Create(ctx, CreatePrepareInput{
		OperationType: string(domain.OpDeposit),
		CallerAddress: input.CallerAddress,
		SaleID:        &saleID,
		PoolIndex:     &poolIdx,
		Calldata:      calldata,
		TargetAddress: sale.ContractAddress,
		Value:         "0",
	})
}

// HarvestInput 用户结算的输入参数。
type HarvestInput struct {
	CallerAddress string `json:"caller_address" binding:"required"`
	SaleID        int64  `json:"sale_id" binding:"required"`
	PoolIndex     int64  `json:"pool_index" binding:"required"`
}

// Harvest 用户结算：生成 harvest calldata + 创建 Prepare 交易。
func (s *UserService) Harvest(ctx context.Context, input HarvestInput) (*domain.PrepareTx, error) {
	sale, err := s.saleRepo.FindByID(ctx, input.SaleID)
	if err != nil {
		return nil, fmt.Errorf("查询 sale %d: %w", input.SaleID, err)
	}
	if sale.Status != domain.SaleDeployed {
		return nil, fmt.Errorf("sale %d 尚未部署: %w", input.SaleID, domain.ErrSaleNotDeployed)
	}

	calldata, err := s.encoder.EncodeCall(s.contract, "harvest",
		big.NewInt(input.PoolIndex),
	)
	if err != nil {
		return nil, fmt.Errorf("编码 harvest: %w", err)
	}

	saleID := input.SaleID
	poolIdx := input.PoolIndex
	return s.prepareSvc.Create(ctx, CreatePrepareInput{
		OperationType: string(domain.OpHarvest),
		CallerAddress: input.CallerAddress,
		SaleID:        &saleID,
		PoolIndex:     &poolIdx,
		Calldata:      calldata,
		TargetAddress: sale.ContractAddress,
		Value:         "0",
	})
}

// ReleaseInput 用户释放 vesting 的输入参数。
type ReleaseInput struct {
	CallerAddress string `json:"caller_address" binding:"required"`
	ScheduleID    int64  `json:"schedule_id" binding:"required"`
}

// Release 用户释放 vesting：生成 release calldata + 创建 Prepare 交易。
func (s *UserService) Release(ctx context.Context, input ReleaseInput) (*domain.PrepareTx, error) {
	// 通过 ScheduleID 反查 vesting_schedule 获取 sale_id，再查 sale 表
	schedule, err := s.findVestingScheduleByID(ctx, input.ScheduleID)
	if err != nil {
		return nil, fmt.Errorf("查询 vesting schedule %d: %w", input.ScheduleID, err)
	}

	sale, err := s.saleRepo.FindByID(ctx, schedule.SaleID)
	if err != nil {
		return nil, fmt.Errorf("查询 sale %d: %w", schedule.SaleID, err)
	}
	if sale.Status != domain.SaleDeployed {
		return nil, fmt.Errorf("sale %d 尚未部署: %w", schedule.SaleID, domain.ErrSaleNotDeployed)
	}

	calldata, err := s.encoder.EncodeCall(s.contract, "release",
		big.NewInt(input.ScheduleID),
	)
	if err != nil {
		return nil, fmt.Errorf("编码 release: %w", err)
	}

	saleID := schedule.SaleID
	return s.prepareSvc.Create(ctx, CreatePrepareInput{
		OperationType: string(domain.OpRelease),
		CallerAddress: input.CallerAddress,
		SaleID:        &saleID,
		Calldata:      calldata,
		TargetAddress: sale.ContractAddress,
		Value:         "0",
	})
}
