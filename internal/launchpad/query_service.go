package launchpad

import (
	"context"
	"fmt"
	"math/big"

	"github.com/mousecake-go/mousecake-go/internal/launchpad/domain"
)

// UserQueryService 用户查询用例，提供销售列表、Tier 查询、白名单检查、申购记录、结算信息等只读功能。
type UserQueryService struct {
	saleRepo      domain.SaleRepository
	poolRepo      domain.PoolRepository
	metaRepo      domain.SaleMetaRepository
	tierLimitRepo domain.TierLimitRepository
	whitelistRepo domain.WhitelistRepository
	depositRepo   domain.DepositRepository
	userPoolRepo  domain.UserPoolStateRepository
	creditRepo    domain.UserCreditRepository
	harvestRepo   domain.HarvestRepository
	vestingRepo   domain.VestingScheduleRepository
	releaseRepo   domain.VestingReleaseRepository
	chain         domain.ChainReader
}

// NewUserQueryService 创建 UserQueryService。
func NewUserQueryService(
	saleRepo domain.SaleRepository,
	poolRepo domain.PoolRepository,
	metaRepo domain.SaleMetaRepository,
	tierLimitRepo domain.TierLimitRepository,
	whitelistRepo domain.WhitelistRepository,
	depositRepo domain.DepositRepository,
	userPoolRepo domain.UserPoolStateRepository,
	creditRepo domain.UserCreditRepository,
	harvestRepo domain.HarvestRepository,
	vestingRepo domain.VestingScheduleRepository,
	releaseRepo domain.VestingReleaseRepository,
	chain domain.ChainReader,
) *UserQueryService {
	return &UserQueryService{
		saleRepo:      saleRepo,
		poolRepo:      poolRepo,
		metaRepo:      metaRepo,
		tierLimitRepo: tierLimitRepo,
		whitelistRepo: whitelistRepo,
		depositRepo:   depositRepo,
		userPoolRepo:  userPoolRepo,
		creditRepo:    creditRepo,
		harvestRepo:   harvestRepo,
		vestingRepo:   vestingRepo,
		releaseRepo:   releaseRepo,
		chain:         chain,
	}
}

// SaleListResult 销售列表查询结果。
type SaleListResult struct {
	Sales []SaleWithMeta `json:"sales"`
	Total int64          `json:"total"`
}

// SaleWithMeta 销售及其元信息。
type SaleWithMeta struct {
	Sale *domain.Sale     `json:"sale"`
	Meta *domain.SaleMeta `json:"meta,omitempty"`
}

// ListPublicSales 分页查询公开可见的销售列表。
func (s *UserQueryService) ListPublicSales(ctx context.Context, page, pageSize int) (*SaleListResult, error) {
	sales, total, err := s.saleRepo.FindPublicList(ctx, page, pageSize)
	if err != nil {
		return nil, fmt.Errorf("查询销售列表: %w", err)
	}

	result := make([]SaleWithMeta, len(sales))
	for i, sale := range sales {
		result[i] = SaleWithMeta{Sale: sale}
	}

	return &SaleListResult{Sales: result, Total: total}, nil
}

// GetSaleDetail 查询销售详情。
func (s *UserQueryService) GetSaleDetail(ctx context.Context, saleID int64) (*SaleDetail, error) {
	sale, err := s.saleRepo.FindByID(ctx, saleID)
	if err != nil {
		return nil, fmt.Errorf("查询销售 %d: %w", saleID, err)
	}

	meta, _ := s.metaRepo.FindBySaleID(ctx, saleID)
	pools, _ := s.poolRepo.FindBySaleID(ctx, saleID)

	return &SaleDetail{
		Sale:  sale,
		Meta:  meta,
		Pools: pools,
	}, nil
}

// SaleDetail 销售详情。
type SaleDetail struct {
	Sale  *domain.Sale     `json:"sale"`
	Meta  *domain.SaleMeta `json:"meta,omitempty"`
	Pools []*domain.Pool   `json:"pools"`
}

// GetUserTier 查询用户在 MouseTier 合约的实时 Tier。
func (s *UserQueryService) GetUserTier(ctx context.Context, mouseTierAddress, userAddress string) (int, error) {
	tier, err := s.chain.GetUserTier(ctx, mouseTierAddress, userAddress)
	if err != nil {
		return 0, fmt.Errorf("查询用户 Tier: %w", err)
	}
	return tier, nil
}

// CheckWhitelist 检查用户是否在白名单中。
func (s *UserQueryService) CheckWhitelist(ctx context.Context, saleID int64, userAddress string) (bool, error) {
	return s.whitelistRepo.IsWhitelisted(ctx, saleID, userAddress)
}

// GetUserDeposits 查询用户的申购记录。
func (s *UserQueryService) GetUserDeposits(ctx context.Context, userAddress string, saleID int64) ([]*domain.Deposit, error) {
	deposits, err := s.depositRepo.FindByUserAndSale(ctx, userAddress, saleID)
	if err != nil {
		return nil, fmt.Errorf("查询申购记录: %w", err)
	}
	return deposits, nil
}

// GetUserHarvest 查询用户的结算信息。
func (s *UserQueryService) GetUserHarvest(ctx context.Context, userAddress string, saleID int64) ([]*domain.Harvest, error) {
	harvests, err := s.harvestRepo.FindByUserAndSale(ctx, userAddress, saleID)
	if err != nil {
		return nil, fmt.Errorf("查询结算信息: %w", err)
	}
	return harvests, nil
}

// GetUserVesting 查询用户的 vesting 计划。
func (s *UserQueryService) GetUserVesting(ctx context.Context, beneficiary string) ([]*domain.VestingSchedule, error) {
	schedules, err := s.vestingRepo.FindByBeneficiary(ctx, beneficiary)
	if err != nil {
		return nil, fmt.Errorf("查询 vesting 计划: %w", err)
	}
	return schedules, nil
}

// CalculateReleasable 计算 vesting 可释放量。
// 基于链上已同步的数据：可释放量 = 锁仓总量 - 已释放量。
// 详细的 cliff/slice 计算由链上合约执行，链下直接使用同步的结果。
func CalculateReleasable(schedule *domain.VestingSchedule) *big.Int {
	return schedule.Remaining()
}

// EstimateAllocationInput 预估配售计算的输入。
type EstimateAllocationInput struct {
	SaleID    int64  `json:"sale_id" binding:"required"`
	PoolIndex int64  `json:"pool_index" binding:"required"`
	UserTotal string `json:"user_total" binding:"required"`
}

// EstimateAllocation 预估配售计算。
func (s *UserQueryService) EstimateAllocation(ctx context.Context, input EstimateAllocationInput) (*AllocationEstimate, error) {
	pool, err := s.poolRepo.FindBySaleAndPool(ctx, input.SaleID, int(input.PoolIndex))
	if err != nil {
		return nil, fmt.Errorf("查询池 %d/%d: %w", input.SaleID, input.PoolIndex, err)
	}

	userTotal, ok := new(big.Int).SetString(input.UserTotal, 10)
	if !ok {
		return nil, fmt.Errorf("无效金额: %s", input.UserTotal)
	}

	raising := pool.RaisingAmount
	if raising == nil || raising.Sign() == 0 {
		return nil, fmt.Errorf("池尚未配置")
	}

	// 计算预估配售：用户投入 * (发售量 / 募资量)
	offering := pool.OfferingAmount
	if offering == nil || offering.Sign() == 0 {
		return nil, fmt.Errorf("池尚未配置")
	}

	allocation := new(big.Int).Mul(userTotal, offering)
	allocation.Div(allocation, raising)

	return &AllocationEstimate{
		UserTotal:        userTotal,
		Allocation:       allocation,
		IsOversubscribed: false, // 需要实际募资总额来确定
	}, nil
}

// AllocationEstimate 配售预估结果。
type AllocationEstimate struct {
	UserTotal        *big.Int `json:"user_total"`
	Allocation       *big.Int `json:"allocation"`
	IsOversubscribed bool     `json:"is_oversubscribed"`
}
