// Package domain 定义 launchpad 模块的领域模型，包含实体、值对象、状态机和仓库接口。
// 本包不能导入任何外层包（编译期隔离）。
package domain

import (
	"errors"
	"fmt"
	"math/big"
	"slices"
	"time"
)

// SalePhase 表示销售阶段，由链上区块高度和配置参数派生。
type SalePhase string

const (
	// SalePhaseConfiguring 配置阶段：当前区块 < startBlock。
	SalePhaseConfiguring SalePhase = "configuring"
	// SalePhaseLive 募资阶段：startBlock <= 当前区块 <= endBlock。
	SalePhaseLive SalePhase = "live"
	// SalePhaseEnded 结算阶段：当前区块 > endBlock。
	SalePhaseEnded SalePhase = "ended"
)

// Sale 是 IDO 销售合约聚合根实体，表示一个 MousePadByTier 合约实例。
type Sale struct {
	// ID 数据库主键。
	ID int64
	// ContractAddress 链上销售合约地址。
	ContractAddress string
	// ChainID 所属链 ID。
	ChainID int
	// DeployerAddress 合约部署者地址。
	DeployerAddress string
	// OwnerAddress 合约当前 owner 地址。
	OwnerAddress string
	// RaiseTokenAddress 募资代币地址（如 USDT）。
	RaiseTokenAddress string
	// OfferingTokenAddress 发售代币地址（项目方代币）。
	OfferingTokenAddress string
	// MouseTierAddress Mouse Tier NFT 合约地址，用于白名单等级判定。
	MouseTierAddress string
	// StartBlock 募资开始区块高度。
	StartBlock int64
	// EndBlock 募资结束区块高度。
	EndBlock int64
	// VestingStartTime Vesting 锁仓开始时间（Unix 时间戳）。
	VestingStartTime int64
	// VestingRevoked 是否已撤销 vesting。
	VestingRevoked bool
	// MaxBufferBlocks 链上操作允许的最大缓冲区块数，用于判断交易是否过期。
	MaxBufferBlocks int64
	// CreatedAt 记录创建时间。
	CreatedAt time.Time
	// UpdatedAt 记录更新时间。
	UpdatedAt time.Time
	// DeletedAt 软删除时间，nil 表示未删除。
	DeletedAt *time.Time
}

// Phase 根据当前链上区块高度派生销售阶段。
func (s *Sale) Phase(currentBlock int64) SalePhase {
	if currentBlock < s.StartBlock {
		return SalePhaseConfiguring
	}
	if currentBlock <= s.EndBlock {
		return SalePhaseLive
	}
	return SalePhaseEnded
}

// ReconstructSale 从数据库重建 Sale 实体（跳过业务规则校验）。
func ReconstructSale(id int64, contractAddress string, chainID int, deployerAddress, ownerAddress, raiseTokenAddress, offeringTokenAddress, mouseTierAddress string, startBlock, endBlock, vestingStartTime int64, vestingRevoked bool, maxBufferBlocks int64, createdAt, updatedAt time.Time, deletedAt *time.Time) *Sale {
	return &Sale{
		ID:                   id,
		ContractAddress:      contractAddress,
		ChainID:              chainID,
		DeployerAddress:      deployerAddress,
		OwnerAddress:         ownerAddress,
		RaiseTokenAddress:    raiseTokenAddress,
		OfferingTokenAddress: offeringTokenAddress,
		MouseTierAddress:     mouseTierAddress,
		StartBlock:           startBlock,
		EndBlock:             endBlock,
		VestingStartTime:     vestingStartTime,
		VestingRevoked:       vestingRevoked,
		MaxBufferBlocks:      maxBufferBlocks,
		CreatedAt:            createdAt,
		UpdatedAt:            updatedAt,
		DeletedAt:            deletedAt,
	}
}

// SaleMeta 是销售展示元信息，由管理后台维护。
type SaleMeta struct {
	// ID 数据库主键。
	ID int64
	// SaleID 关联的销售 ID。
	SaleID int64
	// Title 销售展示标题。
	Title string
	// Description 销售详细描述。
	Description string
	// BannerURL 横幅图片 URL。
	BannerURL string
	// LogoURL Logo 图片 URL。
	LogoURL string
	// WebsiteURL 项目官网 URL。
	WebsiteURL string
	// SocialLinks 社交媒体链接，JSON 格式存储。
	SocialLinks string
	// Visibility 可见性控制：public 或 hidden。
	Visibility string
	// SortOrder 排序权重，值越小越靠前。
	SortOrder int
	// CreatedAt 记录创建时间。
	CreatedAt time.Time
	// UpdatedAt 记录更新时间。
	UpdatedAt time.Time
}

// IsPublic 返回销售是否公开可见。
func (m *SaleMeta) IsPublic() bool {
	return m.Visibility == "public"
}

// ReconstructSaleMeta 从数据库重建 SaleMeta。
func ReconstructSaleMeta(id, saleID int64, title, description, bannerURL, logoURL, websiteURL, socialLinks, visibility string, sortOrder int, createdAt, updatedAt time.Time) *SaleMeta {
	return &SaleMeta{
		ID: id, SaleID: saleID, Title: title, Description: description,
		BannerURL: bannerURL, LogoURL: logoURL, WebsiteURL: websiteURL,
		SocialLinks: socialLinks, Visibility: visibility, SortOrder: sortOrder,
		CreatedAt: createdAt, UpdatedAt: updatedAt,
	}
}

// Pool 是销售池子实体，每个 Sale 固定 2 个池。
type Pool struct {
	// ID 数据库主键。
	ID int64
	// SaleID 关联的销售 ID。
	SaleID int64
	// PoolIndex 池子索引，0 表示普通池，1 表示特殊池。
	PoolIndex int
	// RaisingAmount 募资目标数量（raise token 精度）。
	RaisingAmount *big.Int
	// OfferingAmount 发售代币总数量（offering token 精度）。
	OfferingAmount *big.Int
	// LimitPerUser 单个用户申购上限（raise token 精度）。
	LimitPerUser *big.Int
	// IsSpecialSale 是否为特殊销售池（如白名单专属池）。
	IsSpecialSale bool
	// HasTax 是否收取税费。
	HasTax bool
	// TaxRate 税率，基数为 10000（如 250 表示 2.5%）。
	TaxRate *big.Int
	// VestingPercentage TGE 释放比例，取值 0-100。
	VestingPercentage int
	// VestingCliff Vesting 锁仓悬崖期时长（秒）。
	VestingCliff int64
	// VestingDuration Vesting 总锁仓时长（秒）。
	VestingDuration int64
	// VestingSlicePeriod Vesting 线性释放间隔（秒）。
	VestingSlicePeriod int64
	// TotalAmount 池子实际募资总额（raise token 精度）。
	TotalAmount *big.Int
	// TotalTax 池子累计收取的税费总额（raise token 精度）。
	TotalTax *big.Int
	// CreatedAt 记录创建时间。
	CreatedAt time.Time
	// UpdatedAt 记录更新时间。
	UpdatedAt time.Time
}

// IsConfigured 检查池子参数是否已配置完成。
func (p *Pool) IsConfigured() bool {
	return p.RaisingAmount.Sign() > 0 && p.OfferingAmount.Sign() > 0
}

// ReconstructPool 从数据库重建 Pool 实体。
func ReconstructPool(id, saleID int64, poolIndex int, raisingAmount, offeringAmount, limitPerUser, taxRate, totalAmount, totalTax *big.Int, isSpecialSale, hasTax bool, vestingPercentage int, vestingCliff, vestingDuration, vestingSlicePeriod int64, createdAt, updatedAt time.Time) *Pool {
	return &Pool{
		ID: id, SaleID: saleID, PoolIndex: poolIndex,
		RaisingAmount: raisingAmount, OfferingAmount: offeringAmount,
		LimitPerUser: limitPerUser, IsSpecialSale: isSpecialSale,
		HasTax: hasTax, TaxRate: taxRate,
		VestingPercentage: vestingPercentage, VestingCliff: vestingCliff,
		VestingDuration: vestingDuration, VestingSlicePeriod: vestingSlicePeriod,
		TotalAmount: totalAmount, TotalTax: totalTax,
		CreatedAt: createdAt, UpdatedAt: updatedAt,
	}
}

// Deposit 是用户申购事件实体。
type Deposit struct {
	// ID 数据库主键。
	ID int64
	// SaleID 关联的销售 ID。
	SaleID int64
	// PoolIndex 用户申购的池子索引。
	PoolIndex int
	// UserAddress 申购用户钱包地址。
	UserAddress string
	// Amount 申购金额（raise token 精度）。
	Amount *big.Int
	// TxHash 链上申购交易哈希。
	TxHash string
	// BlockNumber 链上申购交易所在区块高度。
	BlockNumber int64
	// CreatedAt 记录创建时间。
	CreatedAt time.Time
}

// ReconstructDeposit 从数据库重建 Deposit 实体。
func ReconstructDeposit(id, saleID int64, poolIndex int, userAddress, txHash string, amount *big.Int, blockNumber int64, createdAt time.Time) *Deposit {
	return &Deposit{
		ID: id, SaleID: saleID, PoolIndex: poolIndex,
		UserAddress: userAddress, Amount: amount,
		TxHash: txHash, BlockNumber: blockNumber, CreatedAt: createdAt,
	}
}

// UserPoolState 是用户在某个池的累计状态。
type UserPoolState struct {
	// ID 数据库主键。
	ID int64
	// SaleID 关联的销售 ID。
	SaleID int64
	// PoolIndex 用户参与的池子索引。
	PoolIndex int
	// UserAddress 用户钱包地址。
	UserAddress string
	// TotalDeposited 用户在该池的累计申购金额（raise token 精度）。
	TotalDeposited *big.Int
	// Claimed 用户是否已完成结算领取。
	Claimed bool
	// CreatedAt 记录创建时间。
	CreatedAt time.Time
	// UpdatedAt 记录更新时间。
	UpdatedAt time.Time
}

// ReconstructUserPoolState 从数据库重建 UserPoolState。
func ReconstructUserPoolState(id, saleID int64, poolIndex int, userAddress string, totalDeposited *big.Int, claimed bool, createdAt, updatedAt time.Time) *UserPoolState {
	return &UserPoolState{
		ID: id, SaleID: saleID, PoolIndex: poolIndex,
		UserAddress: userAddress, TotalDeposited: totalDeposited,
		Claimed: claimed, CreatedAt: createdAt, UpdatedAt: updatedAt,
	}
}

// UserCredit 是用户在某 sale 的累计信用使用。
type UserCredit struct {
	// ID 数据库主键。
	ID int64
	// SaleID 关联的销售 ID。
	SaleID int64
	// UserAddress 用户钱包地址。
	UserAddress string
	// CreditUsed 用户在该 sale 已使用的信用额度（raise token 精度）。
	CreditUsed *big.Int
	// CreatedAt 记录创建时间。
	CreatedAt time.Time
	// UpdatedAt 记录更新时间。
	UpdatedAt time.Time
}

// ReconstructUserCredit 从数据库重建 UserCredit。
func ReconstructUserCredit(id, saleID int64, userAddress string, creditUsed *big.Int, createdAt, updatedAt time.Time) *UserCredit {
	return &UserCredit{
		ID: id, SaleID: saleID, UserAddress: userAddress,
		CreditUsed: creditUsed, CreatedAt: createdAt, UpdatedAt: updatedAt,
	}
}

// Harvest 是用户结算结果实体。
type Harvest struct {
	// ID 数据库主键。
	ID int64
	// SaleID 关联的销售 ID。
	SaleID int64
	// PoolIndex 结算的池子索引。
	PoolIndex int
	// UserAddress 结算用户钱包地址。
	UserAddress string
	// IsOverflow 是否超额认购（池子募资总额超过目标）。
	IsOverflow bool
	// OfferingAmount 用户可获得的发售代币数量（offering token 精度）。
	OfferingAmount *big.Int
	// PayAmount 用户实际需支付的金额（raise token 精度）。
	PayAmount *big.Int
	// RaiseRefund 超额部分的退款金额（raise token 精度）。
	RaiseRefund *big.Int
	// TaxAmount 扣除的税费金额（raise token 精度）。
	TaxAmount *big.Int
	// TGEAmount TGE（代币生成事件）时一次性释放的数量（offering token 精度）。
	TGEAmount *big.Int
	// VestingAmount 进入 vesting 锁仓的代币数量（offering token 精度）。
	VestingAmount *big.Int
	// TxHash 链上结算交易哈希。
	TxHash string
	// BlockNumber 链上结算交易所在区块高度。
	BlockNumber int64
	// CreatedAt 记录创建时间。
	CreatedAt time.Time
}

// ReconstructHarvest 从数据库重建 Harvest 实体。
func ReconstructHarvest(id, saleID int64, poolIndex int, userAddress string, isOverflow bool, offeringAmount, payAmount, raiseRefund, taxAmount, tgeAmount, vestingAmount *big.Int, txHash string, blockNumber int64, createdAt time.Time) *Harvest {
	return &Harvest{
		ID: id, SaleID: saleID, PoolIndex: poolIndex,
		UserAddress: userAddress, IsOverflow: isOverflow,
		OfferingAmount: offeringAmount, PayAmount: payAmount,
		RaiseRefund: raiseRefund, TaxAmount: taxAmount,
		TGEAmount: tgeAmount, VestingAmount: vestingAmount,
		TxHash: txHash, BlockNumber: blockNumber, CreatedAt: createdAt,
	}
}

// VestingSchedule 是 vesting 锁仓计划实体。
type VestingSchedule struct {
	// ID 数据库主键。
	ID int64
	// SaleID 关联的销售 ID。
	SaleID int64
	// PoolIndex 关联的池子索引。
	PoolIndex int
	// ScheduleID 链上 vesting 合约分配的计划 ID。
	ScheduleID int64
	// Beneficiary 受益人钱包地址。
	Beneficiary string
	// AmountTotal 锁仓总量（offering token 精度）。
	AmountTotal *big.Int
	// Released 已释放总量（offering token 精度）。
	Released *big.Int
	// CreatedAt 记录创建时间。
	CreatedAt time.Time
	// UpdatedAt 记录更新时间。
	UpdatedAt time.Time
}

// Remaining 返回剩余未释放量。
func (v *VestingSchedule) Remaining() *big.Int {
	r := new(big.Int).Sub(v.AmountTotal, v.Released)
	if r.Sign() < 0 {
		return new(big.Int)
	}
	return r
}

// ReconstructVestingSchedule 从数据库重建 VestingSchedule。
func ReconstructVestingSchedule(id, saleID int64, poolIndex int, scheduleID int64, beneficiary string, amountTotal, released *big.Int, createdAt, updatedAt time.Time) *VestingSchedule {
	return &VestingSchedule{
		ID: id, SaleID: saleID, PoolIndex: poolIndex,
		ScheduleID: scheduleID, Beneficiary: beneficiary,
		AmountTotal: amountTotal, Released: released,
		CreatedAt: createdAt, UpdatedAt: updatedAt,
	}
}

// VestingRelease 是 vesting 释放记录实体。
type VestingRelease struct {
	// ID 数据库主键。
	ID int64
	// ScheduleID 关联的 vesting 计划 ID。
	ScheduleID int64
	// Amount 本次释放数量（offering token 精度）。
	Amount *big.Int
	// TxHash 链上释放交易哈希。
	TxHash string
	// BlockNumber 链上释放交易所在区块高度。
	BlockNumber int64
	// CreatedAt 记录创建时间。
	CreatedAt time.Time
}

// ReconstructVestingRelease 从数据库重建 VestingRelease。
func ReconstructVestingRelease(id, scheduleID int64, amount *big.Int, txHash string, blockNumber int64, createdAt time.Time) *VestingRelease {
	return &VestingRelease{
		ID: id, ScheduleID: scheduleID, Amount: amount,
		TxHash: txHash, BlockNumber: blockNumber, CreatedAt: createdAt,
	}
}

// Token 是代币元信息实体。
type Token struct {
	// ID 数据库主键。
	ID int64
	// Address 链上代币合约地址。
	Address string
	// ChainID 所属链 ID。
	ChainID int
	// Name 代币名称（如 Tether USD）。
	Name string
	// Symbol 代币符号（如 USDT）。
	Symbol string
	// Decimals 代币精度位数。
	Decimals int
	// LogoURL 代币 Logo 图片 URL。
	LogoURL string
	// CreatedAt 记录创建时间。
	CreatedAt time.Time
	// UpdatedAt 记录更新时间。
	UpdatedAt time.Time
}

// ReconstructToken 从数据库重建 Token。
func ReconstructToken(id int64, address string, chainID int, name, symbol string, decimals int, logoURL string, createdAt, updatedAt time.Time) *Token {
	return &Token{
		ID: id, Address: address, ChainID: chainID,
		Name: name, Symbol: symbol, Decimals: decimals,
		LogoURL: logoURL, CreatedAt: createdAt, UpdatedAt: updatedAt,
	}
}

// TierLimit 是 Tier 额度配置。
type TierLimit struct {
	// ID 数据库主键。
	ID int64
	// SaleID 关联的销售 ID。
	SaleID int64
	// Tier Mouse Tier NFT 等级。
	Tier int
	// CreditLimit 该等级对应的信用额度上限（raise token 精度）。
	CreditLimit *big.Int
	// CreatedAt 记录创建时间。
	CreatedAt time.Time
	// UpdatedAt 记录更新时间。
	UpdatedAt time.Time
}

// ReconstructTierLimit 从数据库重建 TierLimit。
func ReconstructTierLimit(id, saleID int64, tier int, creditLimit *big.Int, createdAt, updatedAt time.Time) *TierLimit {
	return &TierLimit{
		ID: id, SaleID: saleID, Tier: tier,
		CreditLimit: creditLimit, CreatedAt: createdAt, UpdatedAt: updatedAt,
	}
}

// WhitelistEntry 是白名单记录。
type WhitelistEntry struct {
	// ID 数据库主键。
	ID int64
	// SaleID 关联的销售 ID。
	SaleID int64
	// Address 白名单用户钱包地址。
	Address string
	// IsActive 白名单是否生效。
	IsActive bool
	// CreatedAt 记录创建时间。
	CreatedAt time.Time
	// UpdatedAt 记录更新时间。
	UpdatedAt time.Time
}

// ReconstructWhitelistEntry 从数据库重建 WhitelistEntry。
func ReconstructWhitelistEntry(id, saleID int64, address string, isActive bool, createdAt, updatedAt time.Time) *WhitelistEntry {
	return &WhitelistEntry{
		ID: id, SaleID: saleID, Address: address,
		IsActive: isActive, CreatedAt: createdAt, UpdatedAt: updatedAt,
	}
}

// PrepareTxStatus 表示 Prepare 交易状态。
type PrepareTxStatus string

const (
	// PrepareTxPending 待签名广播。
	PrepareTxPending PrepareTxStatus = "pending"
	// PrepareTxSigned 已签名待广播。
	PrepareTxSigned PrepareTxStatus = "signed"
	// PrepareTxBroadcast 已广播待确认。
	PrepareTxBroadcast PrepareTxStatus = "broadcast"
	// PrepareTxConfirmed 链上确认成功。
	PrepareTxConfirmed PrepareTxStatus = "confirmed"
	// PrepareTxReverted 链上执行失败。
	PrepareTxReverted PrepareTxStatus = "reverted"
	// PrepareTxExpired 超时过期。
	PrepareTxExpired PrepareTxStatus = "expired"
	// PrepareTxFailed 处理失败。
	PrepareTxFailed PrepareTxStatus = "failed"
)

// prepareTxTransitions 合法的状态转换表。
var prepareTxTransitions = map[PrepareTxStatus][]PrepareTxStatus{
	PrepareTxPending:   {PrepareTxSigned, PrepareTxBroadcast, PrepareTxExpired},
	PrepareTxSigned:    {PrepareTxBroadcast, PrepareTxExpired},
	PrepareTxBroadcast: {PrepareTxConfirmed, PrepareTxReverted, PrepareTxFailed, PrepareTxExpired},
	PrepareTxConfirmed: {},
	PrepareTxReverted:  {},
	PrepareTxExpired:   {},
	PrepareTxFailed:    {},
}

// CanTransitionTo 检查是否可以转换到目标状态。
func (s PrepareTxStatus) CanTransitionTo(target PrepareTxStatus) bool {
	return slices.Contains(prepareTxTransitions[s], target)
}

// PrepareTxOperationType 表示 Prepare 交易操作类型。
type PrepareTxOperationType string

const (
	// OpCreateSale 创建 IDO 销售合约。
	OpCreateSale PrepareTxOperationType = "create_sale"
	// OpSetPool 设置池参数。
	OpSetPool PrepareTxOperationType = "set_pool"
	// OpSetTierLimits 设置 Tier 额度。
	OpSetTierLimits PrepareTxOperationType = "set_tier_limits"
	// OpAddWhitelist 添加白名单。
	OpAddWhitelist PrepareTxOperationType = "add_whitelist"
	// OpRemoveWhitelist 移除白名单。
	OpRemoveWhitelist PrepareTxOperationType = "remove_whitelist"
	// OpSetStartEndBlock 设置销售时间窗。
	OpSetStartEndBlock PrepareTxOperationType = "set_start_end_block"
	// OpRevoke 撤销 vesting。
	OpRevoke PrepareTxOperationType = "revoke"
	// OpFinalWithdraw 提取资金。
	OpFinalWithdraw PrepareTxOperationType = "final_withdraw"
	// OpRecoverToken 救援误转代币。
	OpRecoverToken PrepareTxOperationType = "recover_token"
	// OpDeposit 用户申购。
	OpDeposit PrepareTxOperationType = "deposit"
	// OpHarvest 用户结算。
	OpHarvest PrepareTxOperationType = "harvest"
	// OpRelease 用户释放 vesting。
	OpRelease PrepareTxOperationType = "release"
)

// adminOperationTypes 管理员操作类型集合。
var adminOperationTypes = map[PrepareTxOperationType]bool{
	OpCreateSale:       true,
	OpSetPool:          true,
	OpSetTierLimits:    true,
	OpAddWhitelist:     true,
	OpRemoveWhitelist:  true,
	OpSetStartEndBlock: true,
	OpRevoke:           true,
	OpFinalWithdraw:    true,
	OpRecoverToken:     true,
}

// userOperationTypes 用户操作类型集合。
var userOperationTypes = map[PrepareTxOperationType]bool{
	OpDeposit: true,
	OpHarvest: true,
	OpRelease: true,
}

// IsAdminOperation 检查是否为管理员操作。
func (t PrepareTxOperationType) IsAdminOperation() bool {
	return adminOperationTypes[t]
}

// IsUserOperation 检查是否为用户操作。
func (t PrepareTxOperationType) IsUserOperation() bool {
	return userOperationTypes[t]
}

// IsValidOperationType 检查操作类型是否合法。
func IsValidOperationType(t PrepareTxOperationType) bool {
	return adminOperationTypes[t] || userOperationTypes[t]
}

// PrepareTx 是 Prepare 交易记录实体。
type PrepareTx struct {
	// ID 数据库主键。
	ID int64
	// SaleID 关联的销售 ID，部分操作（如创建合约）可能为 nil。
	SaleID *int64
	// PoolIndex 关联的池子索引，部分操作可能为 nil。
	PoolIndex *int64
	// OperationType 操作类型（创建合约、设置池参数、申购、结算等）。
	OperationType PrepareTxOperationType
	// CallerAddress 发起调用的用户或管理员地址。
	CallerAddress string
	// Calldata 编码后的合约调用数据（hex）。
	Calldata string
	// CalldataHash 调用数据的哈希，用于前端验证数据完整性。
	CalldataHash string
	// Status 交易当前状态。
	Status PrepareTxStatus
	// TxHash 链上交易哈希，广播后填充。
	TxHash *string
	// BlockNumber 链上交易所在区块高度，确认后填充。
	BlockNumber *int64
	// ErrorMessage 失败时的错误信息。
	ErrorMessage *string
	// ExpiresAt 交易过期时间，过期后不可再提交。
	ExpiresAt time.Time
	// ConfirmedAt 链上确认时间。
	ConfirmedAt *time.Time
	// CreatedAt 记录创建时间。
	CreatedAt time.Time
	// UpdatedAt 记录更新时间。
	UpdatedAt time.Time
}

// IsExpired 检查 Prepare 交易是否已过期。
func (tx *PrepareTx) IsExpired() bool {
	return time.Now().After(tx.ExpiresAt)
}

// Transition 转换 PrepareTx 状态，校验合法性。
func (tx *PrepareTx) Transition(target PrepareTxStatus) error {
	if !tx.Status.CanTransitionTo(target) {
		return fmt.Errorf("非法状态转换: %s → %s: %w", tx.Status, target, ErrInvalidTransition)
	}
	tx.Status = target
	tx.UpdatedAt = time.Now()
	return nil
}

// ReconstructPrepareTx 从数据库重建 PrepareTx。
func ReconstructPrepareTx(id int64, saleID, poolIndex *int64, operationType PrepareTxOperationType, callerAddress, calldata, calldataHash string, status PrepareTxStatus, txHash *string, blockNumber *int64, errorMessage *string, expiresAt time.Time, confirmedAt *time.Time, createdAt, updatedAt time.Time) *PrepareTx {
	return &PrepareTx{
		ID: id, SaleID: saleID, PoolIndex: poolIndex,
		OperationType: operationType, CallerAddress: callerAddress,
		Calldata: calldata, CalldataHash: calldataHash,
		Status: status, TxHash: txHash, BlockNumber: blockNumber,
		ErrorMessage: errorMessage, ExpiresAt: expiresAt,
		ConfirmedAt: confirmedAt, CreatedAt: createdAt, UpdatedAt: updatedAt,
	}
}

// 领域哨兵错误
var (
	// ErrInvalidTransition 非法的 PrepareTx 状态转换。
	ErrInvalidTransition = errors.New("launchpad: 非法状态转换")
	// ErrNotFound 记录未找到。
	ErrNotFound = errors.New("launchpad: 记录未找到")
	// ErrAlreadyClaimed 该池已结算。
	ErrAlreadyClaimed = errors.New("launchpad: 该池已结算")
	// ErrSaleNotEnded 销售尚未结束。
	ErrSaleNotEnded = errors.New("launchpad: 销售尚未结束")
	// ErrSaleNotLive 销售不在募资阶段。
	ErrSaleNotLive = errors.New("launchpad: 销售不在募资阶段")
	// ErrNotWhitelisted 不在白名单中。
	ErrNotWhitelisted = errors.New("launchpad: 不在白名单中")
	// ErrTierLimitExceeded 超出 Tier 额度。
	ErrTierLimitExceeded = errors.New("launchpad: 超出 Tier 额度限制")
	// ErrPoolLimitExceeded 超出池个人上限。
	ErrPoolLimitExceeded = errors.New("launchpad: 超出池个人上限")
	// ErrNoReleasable 当前无可释放量。
	ErrNoReleasable = errors.New("launchpad: 当前无可释放量")
	// ErrInvalidOperationType 无效操作类型。
	ErrInvalidOperationType = errors.New("launchpad: 无效操作类型")
	// ErrPrepareTxNotPending Prepare 交易不在 pending 状态。
	ErrPrepareTxNotPending = errors.New("launchpad: Prepare 交易不在 pending 状态")
	// ErrPrepareTxStatusImmutable Prepare 交易状态不可变更。
	ErrPrepareTxStatusImmutable = errors.New("launchpad: Prepare 交易状态不可变更")
)
