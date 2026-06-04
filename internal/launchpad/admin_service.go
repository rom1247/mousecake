package launchpad

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/mousecake-go/mousecake-go/internal/launchpad/domain"
)

// AdminService 管理员操作用例，生成各种管理员操作的 calldata 并创建 Prepare 交易。
type AdminService struct {
	prepareSvc      *PrepareService
	encoder         EncoderInterface
	deployerName    string // MousePadByTierDeployer 合约名（用于 createSale）
	padContractName string // MousePadByTier 合约名（用于已部署的合约操作）
	// findSaleByID 根据 ID 查询 Sale。
	findSaleByID func(ctx context.Context, saleID int64) (*domain.Sale, error)
	// createSale 创建 Sale 记录。
	createSale func(ctx context.Context, sale *domain.Sale) error
	// deployerAddr Deployer 工厂合约地址，createSale 操作的 to 地址。
	deployerAddr string
}

// NewAdminService 创建 AdminService。
func NewAdminService(
	prepareSvc *PrepareService,
	encoder EncoderInterface,
	deployerName, padContractName string,
	saleRepo *SaleRepository,
	deployerAddr string,
) *AdminService {
	return &AdminService{
		prepareSvc:      prepareSvc,
		encoder:         encoder,
		deployerName:    deployerName,
		padContractName: padContractName,
		findSaleByID:    saleRepo.FindByID,
		createSale:      saleRepo.Create,
		deployerAddr:    deployerAddr,
	}
}

// getSaleContractAddress 查询 sale 表获取 ContractAddress，验证 status=deployed。
func (s *AdminService) getSaleContractAddress(ctx context.Context, saleID int64) (string, error) {
	sale, err := s.findSaleByID(ctx, saleID)
	if err != nil {
		return "", fmt.Errorf("查询 sale %d: %w", saleID, err)
	}
	if sale.Status != domain.SaleDeployed {
		return "", fmt.Errorf("sale %d 尚未部署: %w", saleID, domain.ErrSaleNotDeployed)
	}
	return sale.ContractAddress, nil
}

// CreateSaleInput 创建 IDO 销售的输入参数。
type CreateSaleInput struct {
	// CallerAddress 管理员地址。
	CallerAddress string `json:"caller_address" binding:"required"`
	// RaisingToken 募资代币地址。
	RaisingToken string `json:"raising_token" binding:"required"`
	// OfferingToken 发售代币地址。
	OfferingToken string `json:"offering_token" binding:"required"`
	// Admin 管理员地址（合约参数）。
	Admin string `json:"admin" binding:"required"`
	// MouseTier MouseTier 合约地址。
	MouseTier string `json:"mouse_tier" binding:"required"`
	// StartBlock 销售开始区块号。
	StartBlock int64 `json:"start_block" binding:"required"`
	// EndBlock 销售结束区块号。
	EndBlock int64 `json:"end_block" binding:"required"`
	// ChainID 链 ID。
	ChainID int `json:"chain_id"`
}

// CreateSale 创建 IDO 销售合约（Draft 模式）。
// 立即创建 sale 表记录（status=deploying），同时创建 PrepareTx。
func (s *AdminService) CreateSale(ctx context.Context, input CreateSaleInput) (*domain.PrepareTx, error) {
	calldata, err := s.encoder.EncodeCall(s.deployerName, "createSale",
		common.HexToAddress(input.RaisingToken),
		common.HexToAddress(input.OfferingToken),
		common.HexToAddress(input.Admin),
		common.HexToAddress(input.MouseTier),
		big.NewInt(input.StartBlock),
		big.NewInt(input.EndBlock),
	)
	if err != nil {
		return nil, fmt.Errorf("编码 createSale: %w", err)
	}

	// 创建 draft sale 记录（status=deploying，ContractAddress 为空）
	draftSale := &domain.Sale{
		Status:               domain.SaleDeploying,
		ChainID:              input.ChainID,
		DeployerAddress:      s.deployerAddr,
		OwnerAddress:         input.Admin,
		RaiseTokenAddress:    input.RaisingToken,
		OfferingTokenAddress: input.OfferingToken,
		MouseTierAddress:     input.MouseTier,
		StartBlock:           input.StartBlock,
		EndBlock:             input.EndBlock,
	}
	if err := s.createSale(ctx, draftSale); err != nil {
		return nil, fmt.Errorf("创建 draft sale: %w", err)
	}

	// 创建 PrepareTx，关联 draft sale，target_address 为 Deployer 合约地址
	saleID := draftSale.ID
	return s.prepareSvc.Create(ctx, CreatePrepareInput{
		OperationType: string(domain.OpCreateSale),
		CallerAddress: input.CallerAddress,
		SaleID:        &saleID,
		Calldata:      calldata,
		TargetAddress: s.deployerAddr,
		Value:         "0",
	})
}

// SetPoolInput 设置池参数的输入。
type SetPoolInput struct {
	CallerAddress      string `json:"caller_address" binding:"required"`
	SaleID             int64  `json:"sale_id" binding:"required"`
	PoolIndex          int64  `json:"pool_index" binding:"required"`
	OfferingAmount     string `json:"offering_amount" binding:"required"`
	RaisingAmount      string `json:"raising_amount" binding:"required"`
	LimitPerUser       string `json:"limit_per_user" binding:"required"`
	IsSpecialSale      bool   `json:"is_special_sale" binding:"required"`
	HasTax             bool   `json:"has_tax" binding:"required"`
	VestingPercentage  int64  `json:"vesting_percentage" binding:"required"`
	VestingCliff       int64  `json:"vesting_cliff" binding:"required"`
	VestingSlicePeriod int64  `json:"vesting_slice_period" binding:"required"`
	VestingDuration    int64  `json:"vesting_duration" binding:"required"`
}

// SetPool 设置池参数。查询 sale 表验证 status=deployed，取 ContractAddress 作为 target_address。
func (s *AdminService) SetPool(ctx context.Context, input SetPoolInput) (*domain.PrepareTx, error) {
	contractAddr, err := s.getSaleContractAddress(ctx, input.SaleID)
	if err != nil {
		return nil, err
	}

	offeringAmount, ok := new(big.Int).SetString(input.OfferingAmount, 10)
	if !ok {
		return nil, fmt.Errorf("无效发售金额格式: %s", input.OfferingAmount)
	}
	raisingAmount, ok := new(big.Int).SetString(input.RaisingAmount, 10)
	if !ok {
		return nil, fmt.Errorf("无效募资金额格式: %s", input.RaisingAmount)
	}
	limitPerUser, ok := new(big.Int).SetString(input.LimitPerUser, 10)
	if !ok {
		return nil, fmt.Errorf("无效用户限额格式: %s", input.LimitPerUser)
	}

	calldata, err := s.encoder.EncodeCall(s.padContractName, "setPool",
		big.NewInt(input.PoolIndex),
		offeringAmount,
		raisingAmount,
		limitPerUser,
		input.IsSpecialSale,
		input.HasTax,
		big.NewInt(input.VestingPercentage),
		big.NewInt(input.VestingCliff),
		big.NewInt(input.VestingSlicePeriod),
		big.NewInt(input.VestingDuration),
	)
	if err != nil {
		return nil, fmt.Errorf("编码 setPool: %w", err)
	}

	saleID := input.SaleID
	poolIdx := input.PoolIndex
	return s.prepareSvc.Create(ctx, CreatePrepareInput{
		OperationType: string(domain.OpSetPool),
		CallerAddress: input.CallerAddress,
		SaleID:        &saleID,
		PoolIndex:     &poolIdx,
		Calldata:      calldata,
		TargetAddress: contractAddr,
		Value:         "0",
	})
}

// SetTierLimitsInput 设置 Tier 额度的输入。
type SetTierLimitsInput struct {
	CallerAddress string `json:"caller_address" binding:"required"`
	SaleID        int64  `json:"sale_id" binding:"required"`
	Tier          int64  `json:"tier" binding:"required"`
	Limit         string `json:"limit" binding:"required"`
}

// SetTierLimits 设置 Tier 额度。通过 SaleID 查询 sale 表获取 target_address。
func (s *AdminService) SetTierLimits(ctx context.Context, input SetTierLimitsInput) (*domain.PrepareTx, error) {
	contractAddr, err := s.getSaleContractAddress(ctx, input.SaleID)
	if err != nil {
		return nil, err
	}

	limit, ok := new(big.Int).SetString(input.Limit, 10)
	if !ok {
		return nil, fmt.Errorf("无效 Tier 额度格式: %s", input.Limit)
	}

	calldata, err := s.encoder.EncodeCall(s.padContractName, "setTierLimits",
		big.NewInt(input.Tier),
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("编码 setTierLimits: %w", err)
	}

	saleID := input.SaleID
	return s.prepareSvc.Create(ctx, CreatePrepareInput{
		OperationType: string(domain.OpSetTierLimits),
		CallerAddress: input.CallerAddress,
		SaleID:        &saleID,
		Calldata:      calldata,
		TargetAddress: contractAddr,
		Value:         "0",
	})
}

// WhitelistInput 白名单操作的输入。
type WhitelistInput struct {
	CallerAddress string   `json:"caller_address" binding:"required"`
	SaleID        int64    `json:"sale_id" binding:"required"`
	Users         []string `json:"users" binding:"required"`
}

// AddWhitelist 添加白名单。查询 sale 表验证 status=deployed。
func (s *AdminService) AddWhitelist(ctx context.Context, input WhitelistInput) (*domain.PrepareTx, error) {
	contractAddr, err := s.getSaleContractAddress(ctx, input.SaleID)
	if err != nil {
		return nil, err
	}

	users := make([]common.Address, len(input.Users))
	for i, u := range input.Users {
		users[i] = common.HexToAddress(u)
	}

	calldata, err := s.encoder.EncodeCall(s.padContractName, "addWhitelist", users)
	if err != nil {
		return nil, fmt.Errorf("编码 addWhitelist: %w", err)
	}

	saleID := input.SaleID
	return s.prepareSvc.Create(ctx, CreatePrepareInput{
		OperationType: string(domain.OpAddWhitelist),
		CallerAddress: input.CallerAddress,
		SaleID:        &saleID,
		Calldata:      calldata,
		TargetAddress: contractAddr,
		Value:         "0",
	})
}

// RemoveWhitelist 移除白名单。查询 sale 表验证 status=deployed。
func (s *AdminService) RemoveWhitelist(ctx context.Context, input WhitelistInput) (*domain.PrepareTx, error) {
	contractAddr, err := s.getSaleContractAddress(ctx, input.SaleID)
	if err != nil {
		return nil, err
	}

	users := make([]common.Address, len(input.Users))
	for i, u := range input.Users {
		users[i] = common.HexToAddress(u)
	}

	calldata, err := s.encoder.EncodeCall(s.padContractName, "removeWhitelist", users)
	if err != nil {
		return nil, fmt.Errorf("编码 removeWhitelist: %w", err)
	}

	saleID := input.SaleID
	return s.prepareSvc.Create(ctx, CreatePrepareInput{
		OperationType: string(domain.OpRemoveWhitelist),
		CallerAddress: input.CallerAddress,
		SaleID:        &saleID,
		Calldata:      calldata,
		TargetAddress: contractAddr,
		Value:         "0",
	})
}

// SetStartEndBlockInput 设置销售时间窗的输入。
type SetStartEndBlockInput struct {
	CallerAddress string `json:"caller_address" binding:"required"`
	SaleID        int64  `json:"sale_id" binding:"required"`
	StartBlock    int64  `json:"start_block" binding:"required"`
	EndBlock      int64  `json:"end_block" binding:"required"`
}

// SetStartEndBlock 设置销售时间窗。查询 sale 表验证 status=deployed。
func (s *AdminService) SetStartEndBlock(ctx context.Context, input SetStartEndBlockInput) (*domain.PrepareTx, error) {
	contractAddr, err := s.getSaleContractAddress(ctx, input.SaleID)
	if err != nil {
		return nil, err
	}

	calldata, err := s.encoder.EncodeCall(s.padContractName, "setStartEndBlock",
		big.NewInt(input.StartBlock),
		big.NewInt(input.EndBlock),
	)
	if err != nil {
		return nil, fmt.Errorf("编码 setStartEndBlock: %w", err)
	}

	saleID := input.SaleID
	return s.prepareSvc.Create(ctx, CreatePrepareInput{
		OperationType: string(domain.OpSetStartEndBlock),
		CallerAddress: input.CallerAddress,
		SaleID:        &saleID,
		Calldata:      calldata,
		TargetAddress: contractAddr,
		Value:         "0",
	})
}

// RevokeInput 撤销 vesting 的输入。
type RevokeInput struct {
	CallerAddress string `json:"caller_address" binding:"required"`
	SaleID        int64  `json:"sale_id" binding:"required"`
	PoolIndex     int64  `json:"pool_index" binding:"required"`
}

// Revoke 撤销 vesting。查询 sale 表验证 status=deployed。
func (s *AdminService) Revoke(ctx context.Context, input RevokeInput) (*domain.PrepareTx, error) {
	contractAddr, err := s.getSaleContractAddress(ctx, input.SaleID)
	if err != nil {
		return nil, err
	}

	calldata, err := s.encoder.EncodeCall(s.padContractName, "revoke")
	if err != nil {
		return nil, fmt.Errorf("编码 revoke: %w", err)
	}

	saleID := input.SaleID
	poolIdx := input.PoolIndex
	return s.prepareSvc.Create(ctx, CreatePrepareInput{
		OperationType: string(domain.OpRevoke),
		CallerAddress: input.CallerAddress,
		SaleID:        &saleID,
		PoolIndex:     &poolIdx,
		Calldata:      calldata,
		TargetAddress: contractAddr,
		Value:         "0",
	})
}

// FinalWithdrawInput 提取资金的输入。
type FinalWithdrawInput struct {
	CallerAddress  string `json:"caller_address" binding:"required"`
	SaleID         int64  `json:"sale_id" binding:"required"`
	RaisingAmount  string `json:"raising_amount" binding:"required"`
	OfferingAmount string `json:"offering_amount" binding:"required"`
}

// FinalWithdraw 提取资金。通过 SaleID 查询 sale 表获取 target_address。
func (s *AdminService) FinalWithdraw(ctx context.Context, input FinalWithdrawInput) (*domain.PrepareTx, error) {
	contractAddr, err := s.getSaleContractAddress(ctx, input.SaleID)
	if err != nil {
		return nil, err
	}

	raiseAmount, ok := new(big.Int).SetString(input.RaisingAmount, 10)
	if !ok {
		return nil, fmt.Errorf("无效募资金额格式: %s", input.RaisingAmount)
	}
	offeringAmount, ok := new(big.Int).SetString(input.OfferingAmount, 10)
	if !ok {
		return nil, fmt.Errorf("无效发售金额格式: %s", input.OfferingAmount)
	}

	calldata, err := s.encoder.EncodeCall(s.padContractName, "finalWithdraw",
		raiseAmount, offeringAmount,
	)
	if err != nil {
		return nil, fmt.Errorf("编码 finalWithdraw: %w", err)
	}

	saleID := input.SaleID
	return s.prepareSvc.Create(ctx, CreatePrepareInput{
		OperationType: string(domain.OpFinalWithdraw),
		CallerAddress: input.CallerAddress,
		SaleID:        &saleID,
		Calldata:      calldata,
		TargetAddress: contractAddr,
		Value:         "0",
	})
}

// RecoverTokenInput 救援误转代币的输入。
type RecoverTokenInput struct {
	CallerAddress string `json:"caller_address" binding:"required"`
	SaleID        int64  `json:"sale_id" binding:"required"`
	TokenAddress  string `json:"token_address" binding:"required"`
	To            string `json:"to" binding:"required"`
	Amount        string `json:"amount" binding:"required"`
}

// RecoverToken 救援误转代币。通过 SaleID 查询 sale 表获取 target_address。
func (s *AdminService) RecoverToken(ctx context.Context, input RecoverTokenInput) (*domain.PrepareTx, error) {
	contractAddr, err := s.getSaleContractAddress(ctx, input.SaleID)
	if err != nil {
		return nil, err
	}

	amount, ok := new(big.Int).SetString(input.Amount, 10)
	if !ok {
		return nil, fmt.Errorf("无效金额格式: %s", input.Amount)
	}

	calldata, err := s.encoder.EncodeCall(s.padContractName, "recoverToken",
		common.HexToAddress(input.TokenAddress),
		common.HexToAddress(input.To),
		amount,
	)
	if err != nil {
		return nil, fmt.Errorf("编码 recoverToken: %w", err)
	}

	saleID := input.SaleID
	return s.prepareSvc.Create(ctx, CreatePrepareInput{
		OperationType: string(domain.OpRecoverToken),
		CallerAddress: input.CallerAddress,
		SaleID:        &saleID,
		Calldata:      calldata,
		TargetAddress: contractAddr,
		Value:         "0",
	})
}
