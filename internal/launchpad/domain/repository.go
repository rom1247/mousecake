package domain

import (
	"context"
	"math/big"
	"time"
)

// SaleRepository 定义 Sale 聚合的持久化接口。
type SaleRepository interface {
	// FindByID 根据 ID 查询 Sale（排除软删除）。
	FindByID(ctx context.Context, saleID int64) (*Sale, error)
	// FindByContractAddress 根据合约地址查询 Sale。
	FindByContractAddress(ctx context.Context, address string) (*Sale, error)
	// FindPublicList 分页查询公开可见的销售列表（关联 meta）。
	FindPublicList(ctx context.Context, page, pageSize int) ([]*Sale, int64, error)
	// Create 创建 Sale。
	Create(ctx context.Context, sale *Sale) error
	// Update 更新 Sale。
	Update(ctx context.Context, sale *Sale) error
}

// SaleMetaRepository 定义销售元信息的持久化接口。
type SaleMetaRepository interface {
	// FindBySaleID 根据 sale_id 查询元信息。
	FindBySaleID(ctx context.Context, saleID int64) (*SaleMeta, error)
	// Create 创建元信息。
	Create(ctx context.Context, meta *SaleMeta) error
	// Update 更新元信息。
	Update(ctx context.Context, meta *SaleMeta) error
}

// PoolRepository 定义池子的持久化接口。
type PoolRepository interface {
	// FindBySaleID 按 sale_id 查询所有池。
	FindBySaleID(ctx context.Context, saleID int64) ([]*Pool, error)
	// FindBySaleAndPool 按 sale_id + pool_index 查询单池。
	FindBySaleAndPool(ctx context.Context, saleID int64, poolIndex int) (*Pool, error)
	// Create 创建池。
	Create(ctx context.Context, pool *Pool) error
	// Update 更新池。
	Update(ctx context.Context, pool *Pool) error
}

// TierLimitRepository 定义 Tier 额度的持久化接口。
type TierLimitRepository interface {
	// FindBySaleID 查询 sale 的所有 Tier 额度。
	FindBySaleID(ctx context.Context, saleID int64) ([]*TierLimit, error)
	// FindBySaleAndTier 查询指定 sale 和 tier 的额度。
	FindBySaleAndTier(ctx context.Context, saleID int64, tier int) (*TierLimit, error)
	// Save 保存 Tier 额度（upsert）。
	Save(ctx context.Context, limit *TierLimit) error
}

// WhitelistRepository 定义白名单的持久化接口。
type WhitelistRepository interface {
	// IsWhitelisted 检查用户是否在白名单中。
	IsWhitelisted(ctx context.Context, saleID int64, address string) (bool, error)
}

// DepositRepository 定义申购记录的持久化接口。
type DepositRepository interface {
	// FindByUserAndSale 查询用户在某 sale 的所有申购。
	FindByUserAndSale(ctx context.Context, userAddress string, saleID int64) ([]*Deposit, error)
	// FindByUserAndPool 查询用户在某 sale 某池的所有申购。
	FindByUserAndPool(ctx context.Context, userAddress string, saleID int64, poolIndex int) ([]*Deposit, error)
}

// UserPoolStateRepository 定义用户池内累计状态的持久化接口。
type UserPoolStateRepository interface {
	// FindByUserAndPool 查询用户在某池的累计状态。
	FindByUserAndPool(ctx context.Context, userAddress string, saleID int64, poolIndex int) (*UserPoolState, error)
}

// UserCreditRepository 定义用户信用使用的持久化接口。
type UserCreditRepository interface {
	// FindByUserAndSale 查询用户在某 sale 的累计信用使用。
	FindByUserAndSale(ctx context.Context, userAddress string, saleID int64) (*UserCredit, error)
}

// HarvestRepository 定义结算记录的持久化接口。
type HarvestRepository interface {
	// FindByUserAndSale 查询用户在某 sale 的所有结算记录。
	FindByUserAndSale(ctx context.Context, userAddress string, saleID int64) ([]*Harvest, error)
}

// VestingScheduleRepository 定义 vesting 计划的持久化接口。
type VestingScheduleRepository interface {
	// FindByBeneficiary 查询用户的所有 vesting 计划。
	FindByBeneficiary(ctx context.Context, beneficiary string) ([]*VestingSchedule, error)
	// FindByID 根据 ID 查询 vesting 计划。
	FindByID(ctx context.Context, id int64) (*VestingSchedule, error)
}

// VestingReleaseRepository 定义 vesting 释放记录的持久化接口。
type VestingReleaseRepository interface {
	// FindByScheduleID 查询指定 schedule 的所有释放记录。
	FindByScheduleID(ctx context.Context, scheduleID int64) ([]*VestingRelease, error)
}

// TokenRepository 定义代币元信息的持久化接口。
type TokenRepository interface {
	// FindByID 根据 ID 查询代币。
	FindByID(ctx context.Context, id int64) (*Token, error)
	// FindByAddress 根据地址查询代币。
	FindByAddress(ctx context.Context, address string, chainID int) (*Token, error)
	// Create 创建代币元信息。
	Create(ctx context.Context, token *Token) error
	// Update 更新代币元信息。
	Update(ctx context.Context, token *Token) error
}

// PrepareTxRepository 定义 Prepare 交易的持久化接口。
type PrepareTxRepository interface {
	// FindByID 根据 ID 查询 Prepare 交易。
	FindByID(ctx context.Context, id int64) (*PrepareTx, error)
	// FindPendingByCalldataHash 按 calldata_hash 查询 pending 状态的记录（去重用）。
	FindPendingByCalldataHash(ctx context.Context, calldataHash string) (*PrepareTx, error)
	// FindByCaller 查询指定调用者的 Prepare 交易列表。
	FindByCaller(ctx context.Context, callerAddress string, status string, page, pageSize int) ([]*PrepareTx, int64, error)
	// FindBySaleID 查询指定 sale 的 Prepare 交易列表。
	FindBySaleID(ctx context.Context, saleID int64, page, pageSize int) ([]*PrepareTx, int64, error)
	// FindPendingExpired 查询已过期的 pending 记录。
	FindPendingExpired(ctx context.Context) ([]*PrepareTx, error)
	// FindBroadcastTimeout 查询需要轮询的 broadcast 记录。
	FindBroadcastTimeout(ctx context.Context, timeout time.Duration) ([]*PrepareTx, error)
	// Create 创建 Prepare 交易。
	Create(ctx context.Context, tx *PrepareTx) error
	// UpdateStatus 更新 Prepare 交易状态。
	UpdateStatus(ctx context.Context, id int64, status PrepareTxStatus, updates map[string]any) error
}

// ReceiptInfo 表示链上交易 Receipt 信息。
type ReceiptInfo struct {
	Status      uint64
	BlockNumber uint64
	GasUsed     uint64
}

// ChainReader 定义链上只读交互接口。
type ChainReader interface {
	// GetTransactionReceipt 查询交易 Receipt。
	// 返回 nil, nil 表示交易尚未被打包。
	GetTransactionReceipt(ctx context.Context, txHash string) (*ReceiptInfo, error)
	// GetUserTier 查询用户在 MouseTier 合约的实时 Tier。
	GetUserTier(ctx context.Context, mouseTierAddress string, userAddress string) (int, error)
	// GetUserCredit 查询用户在 MouseTier 合约的实时 Credit。
	GetUserCredit(ctx context.Context, mouseTierAddress string, userAddress string) (*big.Int, error)
}
