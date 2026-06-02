package launchpad

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"

	"github.com/mousecake-go/mousecake-go/internal/launchpad/domain"
)

// --- 持久化对象（PO）定义 ---

type salePO struct {
	ID                   int64      `gorm:"column:id;primaryKey"`
	ContractAddress      string     `gorm:"column:contract_address"`
	ChainID              int        `gorm:"column:chain_id"`
	DeployerAddress      string     `gorm:"column:deployer_address"`
	OwnerAddress         string     `gorm:"column:owner_address"`
	RaiseTokenAddress    string     `gorm:"column:raise_token_address"`
	OfferingTokenAddress string     `gorm:"column:offering_token_address"`
	MouseTierAddress     string     `gorm:"column:mouse_tier_address"`
	StartBlock           int64      `gorm:"column:start_block"`
	EndBlock             int64      `gorm:"column:end_block"`
	VestingStartTime     int64      `gorm:"column:vesting_start_time"`
	VestingRevoked       bool       `gorm:"column:vesting_revoked"`
	MaxBufferBlocks      int64      `gorm:"column:max_buffer_blocks"`
	CreatedAt            time.Time  `gorm:"column:created_at"`
	UpdatedAt            time.Time  `gorm:"column:updated_at"`
	DeletedAt            *time.Time `gorm:"column:deleted_at"`
}

func (salePO) TableName() string { return "launchpad_sales" }

type saleMetaPO struct {
	ID          int64     `gorm:"column:id;primaryKey"`
	SaleID      int64     `gorm:"column:sale_id"`
	Title       *string   `gorm:"column:title"`
	Description *string   `gorm:"column:description"`
	BannerURL   *string   `gorm:"column:banner_url"`
	LogoURL     *string   `gorm:"column:logo_url"`
	WebsiteURL  *string   `gorm:"column:website_url"`
	SocialLinks *string   `gorm:"column:social_links"`
	Visibility  string    `gorm:"column:visibility"`
	SortOrder   int       `gorm:"column:sort_order"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (saleMetaPO) TableName() string { return "launchpad_sale_meta" }

type poolPO struct {
	ID                 int64     `gorm:"column:id;primaryKey"`
	SaleID             int64     `gorm:"column:sale_id"`
	PoolIndex          int       `gorm:"column:pool_index"`
	RaisingAmount      string    `gorm:"column:raising_amount"`
	OfferingAmount     string    `gorm:"column:offering_amount"`
	LimitPerUser       string    `gorm:"column:limit_per_user"`
	IsSpecialSale      bool      `gorm:"column:is_special_sale"`
	HasTax             bool      `gorm:"column:has_tax"`
	TaxRate            string    `gorm:"column:tax_rate"`
	VestingPercentage  int       `gorm:"column:vesting_percentage"`
	VestingCliff       int64     `gorm:"column:vesting_cliff"`
	VestingDuration    int64     `gorm:"column:vesting_duration"`
	VestingSlicePeriod int64     `gorm:"column:vesting_slice_period"`
	TotalAmount        string    `gorm:"column:total_amount"`
	TotalTax           string    `gorm:"column:total_tax"`
	CreatedAt          time.Time `gorm:"column:created_at"`
	UpdatedAt          time.Time `gorm:"column:updated_at"`
}

func (poolPO) TableName() string { return "launchpad_pools" }

type tierLimitPO struct {
	ID          int64     `gorm:"column:id;primaryKey"`
	SaleID      int64     `gorm:"column:sale_id"`
	Tier        int       `gorm:"column:tier"`
	CreditLimit string    `gorm:"column:credit_limit"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (tierLimitPO) TableName() string { return "launchpad_tier_limits" }

type whitelistPO struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	SaleID    int64     `gorm:"column:sale_id"`
	Address   string    `gorm:"column:address"`
	IsActive  bool      `gorm:"column:is_active"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (whitelistPO) TableName() string { return "launchpad_whitelists" }

type depositPO struct {
	ID          int64     `gorm:"column:id;primaryKey"`
	SaleID      int64     `gorm:"column:sale_id"`
	PoolIndex   int       `gorm:"column:pool_index"`
	UserAddress string    `gorm:"column:user_address"`
	Amount      string    `gorm:"column:amount"`
	TxHash      string    `gorm:"column:tx_hash"`
	BlockNumber int64     `gorm:"column:block_number"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

func (depositPO) TableName() string { return "launchpad_deposits" }

type userPoolStatePO struct {
	ID             int64     `gorm:"column:id;primaryKey"`
	SaleID         int64     `gorm:"column:sale_id"`
	PoolIndex      int       `gorm:"column:pool_index"`
	UserAddress    string    `gorm:"column:user_address"`
	TotalDeposited string    `gorm:"column:total_deposited"`
	Claimed        bool      `gorm:"column:claimed"`
	CreatedAt      time.Time `gorm:"column:created_at"`
	UpdatedAt      time.Time `gorm:"column:updated_at"`
}

func (userPoolStatePO) TableName() string { return "launchpad_user_pool_state" }

type userCreditPO struct {
	ID          int64     `gorm:"column:id;primaryKey"`
	SaleID      int64     `gorm:"column:sale_id"`
	UserAddress string    `gorm:"column:user_address"`
	CreditUsed  string    `gorm:"column:credit_used"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (userCreditPO) TableName() string { return "launchpad_user_credit" }

type harvestPO struct {
	ID             int64     `gorm:"column:id;primaryKey"`
	SaleID         int64     `gorm:"column:sale_id"`
	PoolIndex      int       `gorm:"column:pool_index"`
	UserAddress    string    `gorm:"column:user_address"`
	IsOverflow     bool      `gorm:"column:is_overflow"`
	OfferingAmount string    `gorm:"column:offering_amount"`
	PayAmount      string    `gorm:"column:pay_amount"`
	RaiseRefund    string    `gorm:"column:raise_refund"`
	TaxAmount      string    `gorm:"column:tax_amount"`
	TGEAmount      string    `gorm:"column:tge_amount"`
	VestingAmount  string    `gorm:"column:vesting_amount"`
	TxHash         string    `gorm:"column:tx_hash"`
	BlockNumber    int64     `gorm:"column:block_number"`
	CreatedAt      time.Time `gorm:"column:created_at"`
}

func (harvestPO) TableName() string { return "launchpad_harvests" }

type vestingSchedulePO struct {
	ID          int64     `gorm:"column:id;primaryKey"`
	SaleID      int64     `gorm:"column:sale_id"`
	PoolIndex   int       `gorm:"column:pool_index"`
	ScheduleID  int64     `gorm:"column:schedule_id"`
	Beneficiary string    `gorm:"column:beneficiary"`
	AmountTotal string    `gorm:"column:amount_total"`
	Released    string    `gorm:"column:released"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

func (vestingSchedulePO) TableName() string { return "launchpad_vesting_schedules" }

type vestingReleasePO struct {
	ID          int64     `gorm:"column:id;primaryKey"`
	ScheduleID  int64     `gorm:"column:schedule_id"`
	Amount      string    `gorm:"column:amount"`
	TxHash      string    `gorm:"column:tx_hash"`
	BlockNumber int64     `gorm:"column:block_number"`
	CreatedAt   time.Time `gorm:"column:created_at"`
}

func (vestingReleasePO) TableName() string { return "launchpad_vesting_releases" }

type tokenPO struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	Address   string    `gorm:"column:address"`
	ChainID   int       `gorm:"column:chain_id"`
	Name      string    `gorm:"column:name"`
	Symbol    string    `gorm:"column:symbol"`
	Decimals  int       `gorm:"column:decimals"`
	LogoURL   *string   `gorm:"column:logo_url"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (tokenPO) TableName() string { return "launchpad_tokens" }

type prepareTxPO struct {
	ID            int64      `gorm:"column:id;primaryKey"`
	SaleID        *int64     `gorm:"column:sale_id"`
	PoolIndex     *int64     `gorm:"column:pool_index"`
	OperationType string     `gorm:"column:operation_type"`
	CallerAddress string     `gorm:"column:caller_address"`
	Calldata      string     `gorm:"column:calldata"`
	CalldataHash  string     `gorm:"column:calldata_hash"`
	Status        string     `gorm:"column:status"`
	TxHash        *string    `gorm:"column:tx_hash"`
	BlockNumber   *int64     `gorm:"column:block_number"`
	ErrorMessage  *string    `gorm:"column:error_message"`
	ExpiresAt     time.Time  `gorm:"column:expires_at"`
	ConfirmedAt   *time.Time `gorm:"column:confirmed_at"`
	CreatedAt     time.Time  `gorm:"column:created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at"`
}

func (prepareTxPO) TableName() string { return "launchpad_prepare_txs" }

// --- NUMERIC ↔ *big.Int 转换辅助 ---

func bigToInt(s string) *big.Int {
	v := new(big.Int)
	v.SetString(s, 10)
	return v
}

func intToBig(v *big.Int) string {
	if v == nil {
		return "0"
	}
	return v.String()
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func int64Ptr(v int64) *int64 {
	return &v
}

func strVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// --- SaleRepository ---

// SaleRepository 是 domain.SaleRepository 的 Gorm 实现。
type SaleRepository struct {
	db *gorm.DB
}

// NewSaleRepository 创建仓库实例。
func NewSaleRepository(db *gorm.DB) *SaleRepository {
	return &SaleRepository{db: db}
}

// FindByID 根据 ID 查询 Sale。
func (r *SaleRepository) FindByID(ctx context.Context, saleID int64) (*domain.Sale, error) {
	var po salePO
	err := r.db.WithContext(ctx).
		Select("id, contract_address, chain_id, deployer_address, owner_address, raise_token_address, offering_token_address, mouse_tier_address, start_block, end_block, vesting_start_time, vesting_revoked, max_buffer_blocks, created_at, updated_at, deleted_at").
		Where("id = ? AND deleted_at IS NULL", saleID).
		First(&po).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("find sale by id %d: %w", saleID, err)
	}
	return salePOToEntity(&po), nil
}

// FindByContractAddress 根据合约地址查询 Sale。
func (r *SaleRepository) FindByContractAddress(ctx context.Context, address string) (*domain.Sale, error) {
	var po salePO
	err := r.db.WithContext(ctx).
		Select("id, contract_address, chain_id, deployer_address, owner_address, raise_token_address, offering_token_address, mouse_tier_address, start_block, end_block, vesting_start_time, vesting_revoked, max_buffer_blocks, created_at, updated_at, deleted_at").
		Where("contract_address = ? AND deleted_at IS NULL", strings.ToLower(address)).
		First(&po).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("find sale by contract %s: %w", address, err)
	}
	return salePOToEntity(&po), nil
}

// FindPublicList 分页查询公开可见的销售列表。
func (r *SaleRepository) FindPublicList(ctx context.Context, page, pageSize int) ([]*domain.Sale, int64, error) {
	var total int64
	r.db.WithContext(ctx).Model(&saleMetaPO{}).
		Where("visibility = ?", "public").
		Count(&total)

	var metaPOs []saleMetaPO
	err := r.db.WithContext(ctx).
		Select("sale_id").
		Where("visibility = ?", "public").
		Order("sort_order DESC, sale_id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&metaPOs).Error
	if err != nil {
		return nil, 0, fmt.Errorf("find public sales: %w", err)
	}

	if len(metaPOs) == 0 {
		return nil, 0, nil
	}

	saleIDs := make([]int64, len(metaPOs))
	for i, m := range metaPOs {
		saleIDs[i] = m.SaleID
	}

	var salePOs []salePO
	err = r.db.WithContext(ctx).
		Where("id IN ? AND deleted_at IS NULL", saleIDs).
		Find(&salePOs).Error
	if err != nil {
		return nil, 0, fmt.Errorf("find sales by ids: %w", err)
	}

	result := make([]*domain.Sale, len(salePOs))
	for i, po := range salePOs {
		result[i] = salePOToEntity(&po)
	}
	return result, total, nil
}

// Create 创建 Sale。
func (r *SaleRepository) Create(ctx context.Context, sale *domain.Sale) error {
	po := saleEntityToPO(sale)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return fmt.Errorf("create sale: %w", err)
	}
	sale.ID = po.ID
	return nil
}

// Update 更新 Sale。
func (r *SaleRepository) Update(ctx context.Context, sale *domain.Sale) error {
	po := saleEntityToPO(sale)
	err := r.db.WithContext(ctx).
		Select("*").
		Where("id = ?", sale.ID).
		Updates(po).Error
	if err != nil {
		return fmt.Errorf("update sale %d: %w", sale.ID, err)
	}
	return nil
}

// --- SaleMetaRepository ---

// SaleMetaRepository 是 domain.SaleMetaRepository 的 Gorm 实现。
type SaleMetaRepository struct {
	db *gorm.DB
}

// NewSaleMetaRepository 创建仓库实例。
func NewSaleMetaRepository(db *gorm.DB) *SaleMetaRepository {
	return &SaleMetaRepository{db: db}
}

// FindBySaleID 根据 sale_id 查询元信息。
func (r *SaleMetaRepository) FindBySaleID(ctx context.Context, saleID int64) (*domain.SaleMeta, error) {
	var po saleMetaPO
	err := r.db.WithContext(ctx).
		Where("sale_id = ?", saleID).
		First(&po).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("find sale meta by sale_id %d: %w", saleID, err)
	}
	return saleMetaPOToEntity(&po), nil
}

// Create 创建元信息。
func (r *SaleMetaRepository) Create(ctx context.Context, meta *domain.SaleMeta) error {
	po := saleMetaEntityToPO(meta)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return fmt.Errorf("create sale meta: %w", err)
	}
	meta.ID = po.ID
	return nil
}

// Update 更新元信息。
func (r *SaleMetaRepository) Update(ctx context.Context, meta *domain.SaleMeta) error {
	po := saleMetaEntityToPO(meta)
	err := r.db.WithContext(ctx).
		Select("*").
		Where("id = ?", meta.ID).
		Updates(po).Error
	if err != nil {
		return fmt.Errorf("update sale meta %d: %w", meta.ID, err)
	}
	return nil
}

// --- PoolRepository ---

// PoolRepository 是 domain.PoolRepository 的 Gorm 实现。
type PoolRepository struct {
	db *gorm.DB
}

// NewPoolRepository 创建仓库实例。
func NewPoolRepository(db *gorm.DB) *PoolRepository {
	return &PoolRepository{db: db}
}

// FindBySaleID 按 sale_id 查询所有池。
func (r *PoolRepository) FindBySaleID(ctx context.Context, saleID int64) ([]*domain.Pool, error) {
	var pos []poolPO
	err := r.db.WithContext(ctx).
		Where("sale_id = ?", saleID).
		Order("pool_index").
		Find(&pos).Error
	if err != nil {
		return nil, fmt.Errorf("find pools by sale_id %d: %w", saleID, err)
	}
	result := make([]*domain.Pool, len(pos))
	for i, po := range pos {
		result[i] = poolPOToEntity(&po)
	}
	return result, nil
}

// FindBySaleAndPool 按 sale_id + pool_index 查询单池。
func (r *PoolRepository) FindBySaleAndPool(ctx context.Context, saleID int64, poolIndex int) (*domain.Pool, error) {
	var po poolPO
	err := r.db.WithContext(ctx).
		Where("sale_id = ? AND pool_index = ?", saleID, poolIndex).
		First(&po).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("find pool by sale %d pool %d: %w", saleID, poolIndex, err)
	}
	return poolPOToEntity(&po), nil
}

// Create 创建池。
func (r *PoolRepository) Create(ctx context.Context, pool *domain.Pool) error {
	po := poolEntityToPO(pool)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return fmt.Errorf("create pool: %w", err)
	}
	pool.ID = po.ID
	return nil
}

// Update 更新池。
func (r *PoolRepository) Update(ctx context.Context, pool *domain.Pool) error {
	po := poolEntityToPO(pool)
	err := r.db.WithContext(ctx).
		Select("*").
		Where("id = ?", pool.ID).
		Updates(po).Error
	if err != nil {
		return fmt.Errorf("update pool %d: %w", pool.ID, err)
	}
	return nil
}

// --- TierLimitRepository ---

// TierLimitRepository 是 domain.TierLimitRepository 的 Gorm 实现。
type TierLimitRepository struct {
	db *gorm.DB
}

// NewTierLimitRepository 创建仓库实例。
func NewTierLimitRepository(db *gorm.DB) *TierLimitRepository {
	return &TierLimitRepository{db: db}
}

// FindBySaleID 查询 sale 的所有 Tier 额度。
func (r *TierLimitRepository) FindBySaleID(ctx context.Context, saleID int64) ([]*domain.TierLimit, error) {
	var pos []tierLimitPO
	err := r.db.WithContext(ctx).
		Where("sale_id = ?", saleID).
		Order("tier").
		Find(&pos).Error
	if err != nil {
		return nil, fmt.Errorf("find tier limits by sale_id %d: %w", saleID, err)
	}
	result := make([]*domain.TierLimit, len(pos))
	for i, po := range pos {
		result[i] = tierLimitPOToEntity(&po)
	}
	return result, nil
}

// FindBySaleAndTier 查询指定 sale 和 tier 的额度。
func (r *TierLimitRepository) FindBySaleAndTier(ctx context.Context, saleID int64, tier int) (*domain.TierLimit, error) {
	var po tierLimitPO
	err := r.db.WithContext(ctx).
		Where("sale_id = ? AND tier = ?", saleID, tier).
		First(&po).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("find tier limit: %w", err)
	}
	return tierLimitPOToEntity(&po), nil
}

// Save 保存 Tier 额度（upsert）。
func (r *TierLimitRepository) Save(ctx context.Context, limit *domain.TierLimit) error {
	po := tierLimitEntityToPO(limit)
	err := r.db.WithContext(ctx).
		Where("sale_id = ? AND tier = ?", limit.SaleID, limit.Tier).
		Assign(po).
		FirstOrCreate(po).Error
	if err != nil {
		return fmt.Errorf("save tier limit: %w", err)
	}
	limit.ID = po.ID
	return nil
}

// --- WhitelistRepository ---

// WhitelistRepository 是 domain.WhitelistRepository 的 Gorm 实现。
type WhitelistRepository struct {
	db *gorm.DB
}

// NewWhitelistRepository 创建仓库实例。
func NewWhitelistRepository(db *gorm.DB) *WhitelistRepository {
	return &WhitelistRepository{db: db}
}

// IsWhitelisted 检查用户是否在白名单中。
func (r *WhitelistRepository) IsWhitelisted(ctx context.Context, saleID int64, address string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&whitelistPO{}).
		Where("sale_id = ? AND address = ? AND is_active = ?", saleID, strings.ToLower(address), true).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("check whitelist: %w", err)
	}
	return count > 0, nil
}

// --- DepositRepository ---

// DepositRepository 是 domain.DepositRepository 的 Gorm 实现。
type DepositRepository struct {
	db *gorm.DB
}

// NewDepositRepository 创建仓库实例。
func NewDepositRepository(db *gorm.DB) *DepositRepository {
	return &DepositRepository{db: db}
}

// FindByUserAndSale 查询用户在某 sale 的所有申购。
func (r *DepositRepository) FindByUserAndSale(ctx context.Context, userAddress string, saleID int64) ([]*domain.Deposit, error) {
	var pos []depositPO
	err := r.db.WithContext(ctx).
		Where("sale_id = ? AND user_address = ?", saleID, strings.ToLower(userAddress)).
		Order("block_number").
		Find(&pos).Error
	if err != nil {
		return nil, fmt.Errorf("find deposits: %w", err)
	}
	result := make([]*domain.Deposit, len(pos))
	for i, po := range pos {
		result[i] = depositPOToEntity(&po)
	}
	return result, nil
}

// FindByUserAndPool 查询用户在某 sale 某池的所有申购。
func (r *DepositRepository) FindByUserAndPool(ctx context.Context, userAddress string, saleID int64, poolIndex int) ([]*domain.Deposit, error) {
	var pos []depositPO
	err := r.db.WithContext(ctx).
		Where("sale_id = ? AND pool_index = ? AND user_address = ?", saleID, poolIndex, strings.ToLower(userAddress)).
		Order("block_number").
		Find(&pos).Error
	if err != nil {
		return nil, fmt.Errorf("find deposits by pool: %w", err)
	}
	result := make([]*domain.Deposit, len(pos))
	for i, po := range pos {
		result[i] = depositPOToEntity(&po)
	}
	return result, nil
}

// --- UserPoolStateRepository ---

// UserPoolStateRepository 是 domain.UserPoolStateRepository 的 Gorm 实现。
type UserPoolStateRepository struct {
	db *gorm.DB
}

// NewUserPoolStateRepository 创建仓库实例。
func NewUserPoolStateRepository(db *gorm.DB) *UserPoolStateRepository {
	return &UserPoolStateRepository{db: db}
}

// FindByUserAndPool 查询用户在某池的累计状态。
func (r *UserPoolStateRepository) FindByUserAndPool(ctx context.Context, userAddress string, saleID int64, poolIndex int) (*domain.UserPoolState, error) {
	var po userPoolStatePO
	err := r.db.WithContext(ctx).
		Where("sale_id = ? AND pool_index = ? AND user_address = ?", saleID, poolIndex, strings.ToLower(userAddress)).
		First(&po).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("find user pool state: %w", err)
	}
	return userPoolStatePOToEntity(&po), nil
}

// --- UserCreditRepository ---

// UserCreditRepository 是 domain.UserCreditRepository 的 Gorm 实现。
type UserCreditRepository struct {
	db *gorm.DB
}

// NewUserCreditRepository 创建仓库实例。
func NewUserCreditRepository(db *gorm.DB) *UserCreditRepository {
	return &UserCreditRepository{db: db}
}

// FindByUserAndSale 查询用户在某 sale 的累计信用使用。
func (r *UserCreditRepository) FindByUserAndSale(ctx context.Context, userAddress string, saleID int64) (*domain.UserCredit, error) {
	var po userCreditPO
	err := r.db.WithContext(ctx).
		Where("sale_id = ? AND user_address = ?", saleID, strings.ToLower(userAddress)).
		First(&po).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("find user credit: %w", err)
	}
	return userCreditPOToEntity(&po), nil
}

// --- HarvestRepository ---

// HarvestRepository 是 domain.HarvestRepository 的 Gorm 实现。
type HarvestRepository struct {
	db *gorm.DB
}

// NewHarvestRepository 创建仓库实例。
func NewHarvestRepository(db *gorm.DB) *HarvestRepository {
	return &HarvestRepository{db: db}
}

// FindByUserAndSale 查询用户在某 sale 的所有结算记录。
func (r *HarvestRepository) FindByUserAndSale(ctx context.Context, userAddress string, saleID int64) ([]*domain.Harvest, error) {
	var pos []harvestPO
	err := r.db.WithContext(ctx).
		Where("sale_id = ? AND user_address = ?", saleID, strings.ToLower(userAddress)).
		Order("pool_index").
		Find(&pos).Error
	if err != nil {
		return nil, fmt.Errorf("find harvests: %w", err)
	}
	result := make([]*domain.Harvest, len(pos))
	for i, po := range pos {
		result[i] = harvestPOToEntity(&po)
	}
	return result, nil
}

// --- VestingScheduleRepository ---

// VestingScheduleRepository 是 domain.VestingScheduleRepository 的 Gorm 实现。
type VestingScheduleRepository struct {
	db *gorm.DB
}

// NewVestingScheduleRepository 创建仓库实例。
func NewVestingScheduleRepository(db *gorm.DB) *VestingScheduleRepository {
	return &VestingScheduleRepository{db: db}
}

// FindByBeneficiary 查询用户的所有 vesting 计划。
func (r *VestingScheduleRepository) FindByBeneficiary(ctx context.Context, beneficiary string) ([]*domain.VestingSchedule, error) {
	var pos []vestingSchedulePO
	err := r.db.WithContext(ctx).
		Where("beneficiary = ?", strings.ToLower(beneficiary)).
		Order("id").
		Find(&pos).Error
	if err != nil {
		return nil, fmt.Errorf("find vesting schedules: %w", err)
	}
	result := make([]*domain.VestingSchedule, len(pos))
	for i, po := range pos {
		result[i] = vestingSchedulePOToEntity(&po)
	}
	return result, nil
}

// FindByID 根据 ID 查询 vesting 计划。
func (r *VestingScheduleRepository) FindByID(ctx context.Context, id int64) (*domain.VestingSchedule, error) {
	var po vestingSchedulePO
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&po).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("find vesting schedule by id %d: %w", id, err)
	}
	return vestingSchedulePOToEntity(&po), nil
}

// --- VestingReleaseRepository ---

// VestingReleaseRepository 是 domain.VestingReleaseRepository 的 Gorm 实现。
type VestingReleaseRepository struct {
	db *gorm.DB
}

// NewVestingReleaseRepository 创建仓库实例。
func NewVestingReleaseRepository(db *gorm.DB) *VestingReleaseRepository {
	return &VestingReleaseRepository{db: db}
}

// FindByScheduleID 查询指定 schedule 的所有释放记录。
func (r *VestingReleaseRepository) FindByScheduleID(ctx context.Context, scheduleID int64) ([]*domain.VestingRelease, error) {
	var pos []vestingReleasePO
	err := r.db.WithContext(ctx).
		Where("schedule_id = ?", scheduleID).
		Order("block_number").
		Find(&pos).Error
	if err != nil {
		return nil, fmt.Errorf("find vesting releases: %w", err)
	}
	result := make([]*domain.VestingRelease, len(pos))
	for i, po := range pos {
		result[i] = vestingReleasePOToEntity(&po)
	}
	return result, nil
}

// --- TokenRepository ---

// TokenRepository 是 domain.TokenRepository 的 Gorm 实现。
type TokenRepository struct {
	db *gorm.DB
}

// NewTokenRepository 创建仓库实例。
func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

// FindByID 根据 ID 查询代币。
func (r *TokenRepository) FindByID(ctx context.Context, id int64) (*domain.Token, error) {
	var po tokenPO
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&po).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("find token by id %d: %w", id, err)
	}
	return tokenPOToEntity(&po), nil
}

// FindByAddress 根据地址查询代币。
func (r *TokenRepository) FindByAddress(ctx context.Context, address string, chainID int) (*domain.Token, error) {
	var po tokenPO
	err := r.db.WithContext(ctx).
		Where("address = ? AND chain_id = ?", strings.ToLower(address), chainID).
		First(&po).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("find token by address %s: %w", address, err)
	}
	return tokenPOToEntity(&po), nil
}

// Create 创建代币元信息。
func (r *TokenRepository) Create(ctx context.Context, token *domain.Token) error {
	po := tokenEntityToPO(token)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		return fmt.Errorf("create token: %w", err)
	}
	token.ID = po.ID
	return nil
}

// Update 更新代币元信息。
func (r *TokenRepository) Update(ctx context.Context, token *domain.Token) error {
	po := tokenEntityToPO(token)
	err := r.db.WithContext(ctx).
		Select("*").
		Where("id = ?", token.ID).
		Updates(po).Error
	if err != nil {
		return fmt.Errorf("update token %d: %w", token.ID, err)
	}
	return nil
}

// --- PrepareTxRepository ---

// PrepareTxRepository 是 domain.PrepareTxRepository 的 Gorm 实现。
type PrepareTxRepository struct {
	db *gorm.DB
}

// NewPrepareTxRepository 创建仓库实例。
func NewPrepareTxRepository(db *gorm.DB) *PrepareTxRepository {
	return &PrepareTxRepository{db: db}
}

// FindByID 根据 ID 查询 Prepare 交易。
func (r *PrepareTxRepository) FindByID(ctx context.Context, id int64) (*domain.PrepareTx, error) {
	var po prepareTxPO
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&po).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("find prepare_tx by id %d: %w", id, err)
	}
	return prepareTxPOToEntity(&po), nil
}

// FindPendingByCalldataHash 按 calldata_hash 查询 pending 状态的记录。
func (r *PrepareTxRepository) FindPendingByCalldataHash(ctx context.Context, calldataHash string) (*domain.PrepareTx, error) {
	var po prepareTxPO
	err := r.db.WithContext(ctx).
		Where("calldata_hash = ? AND status IN ?", calldataHash, []string{"pending", "signed", "broadcast"}).
		First(&po).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("find pending by calldata_hash: %w", err)
	}
	return prepareTxPOToEntity(&po), nil
}

// FindByCaller 查询指定调用者的 Prepare 交易列表。
func (r *PrepareTxRepository) FindByCaller(ctx context.Context, callerAddress string, status string, page, pageSize int) ([]*domain.PrepareTx, int64, error) {
	query := r.db.WithContext(ctx).Model(&prepareTxPO{}).
		Where("caller_address = ?", strings.ToLower(callerAddress))
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var pos []prepareTxPO
	err := query.Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&pos).Error
	if err != nil {
		return nil, 0, fmt.Errorf("find prepare_tx by caller: %w", err)
	}
	result := make([]*domain.PrepareTx, len(pos))
	for i, po := range pos {
		result[i] = prepareTxPOToEntity(&po)
	}
	return result, total, nil
}

// FindBySaleID 查询指定 sale 的 Prepare 交易列表。
func (r *PrepareTxRepository) FindBySaleID(ctx context.Context, saleID int64, page, pageSize int) ([]*domain.PrepareTx, int64, error) {
	query := r.db.WithContext(ctx).Model(&prepareTxPO{}).
		Where("sale_id = ?", saleID)

	var total int64
	query.Count(&total)

	var pos []prepareTxPO
	err := query.Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&pos).Error
	if err != nil {
		return nil, 0, fmt.Errorf("find prepare_tx by sale_id: %w", err)
	}
	result := make([]*domain.PrepareTx, len(pos))
	for i, po := range pos {
		result[i] = prepareTxPOToEntity(&po)
	}
	return result, total, nil
}

// FindPendingExpired 查询已过期的 pending 记录。
func (r *PrepareTxRepository) FindPendingExpired(ctx context.Context) ([]*domain.PrepareTx, error) {
	var pos []prepareTxPO
	err := r.db.WithContext(ctx).
		Where("status IN ? AND expires_at < ?", []string{"pending", "signed"}, time.Now()).
		Find(&pos).Error
	if err != nil {
		return nil, fmt.Errorf("find expired pending: %w", err)
	}
	result := make([]*domain.PrepareTx, len(pos))
	for i, po := range pos {
		result[i] = prepareTxPOToEntity(&po)
	}
	return result, nil
}

// FindBroadcastTimeout 查询需要轮询的 broadcast 记录。
func (r *PrepareTxRepository) FindBroadcastTimeout(ctx context.Context, timeout time.Duration) ([]*domain.PrepareTx, error) {
	cutoff := time.Now().Add(-timeout)
	var pos []prepareTxPO
	err := r.db.WithContext(ctx).
		Where("status = ? AND updated_at < ?", "broadcast", cutoff).
		Find(&pos).Error
	if err != nil {
		return nil, fmt.Errorf("find broadcast timeout: %w", err)
	}
	result := make([]*domain.PrepareTx, len(pos))
	for i, po := range pos {
		result[i] = prepareTxPOToEntity(&po)
	}
	return result, nil
}

// Create 创建 Prepare 交易。
// 活跃状态的 calldata_hash 存在唯一约束，并发创建时若冲突则查询并返回已有记录。
func (r *PrepareTxRepository) Create(ctx context.Context, tx *domain.PrepareTx) error {
	po := prepareTxEntityToPO(tx)
	if err := r.db.WithContext(ctx).Create(po).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && tx.CalldataHash != "" {
			existing, findErr := r.FindPendingByCalldataHash(ctx, tx.CalldataHash)
			if findErr != nil {
				return fmt.Errorf("create prepare_tx 唯一约束冲突后查询失败: %w", findErr)
			}
			if existing != nil {
				*tx = *existing
				return nil
			}
		}
		return fmt.Errorf("create prepare_tx: %w", err)
	}
	tx.ID = po.ID
	return nil
}

// UpdateStatus 更新 Prepare 交易状态。
func (r *PrepareTxRepository) UpdateStatus(ctx context.Context, id int64, status domain.PrepareTxStatus, updates map[string]any) error {
	updates["status"] = string(status)
	updates["updated_at"] = time.Now()

	err := r.db.WithContext(ctx).
		Model(&prepareTxPO{}).
		Where("id = ?", id).
		Updates(updates).Error
	if err != nil {
		return fmt.Errorf("update prepare_tx status %d: %w", id, err)
	}
	return nil
}

// --- PO ↔ Entity 转换函数 ---

func salePOToEntity(po *salePO) *domain.Sale {
	return domain.ReconstructSale(
		po.ID, po.ContractAddress, po.ChainID,
		po.DeployerAddress, po.OwnerAddress,
		po.RaiseTokenAddress, po.OfferingTokenAddress,
		po.MouseTierAddress, po.StartBlock, po.EndBlock,
		po.VestingStartTime, po.VestingRevoked,
		po.MaxBufferBlocks, po.CreatedAt, po.UpdatedAt, po.DeletedAt,
	)
}

func saleEntityToPO(s *domain.Sale) *salePO {
	return &salePO{
		ID: s.ID, ContractAddress: s.ContractAddress, ChainID: s.ChainID,
		DeployerAddress: s.DeployerAddress, OwnerAddress: s.OwnerAddress,
		RaiseTokenAddress: s.RaiseTokenAddress, OfferingTokenAddress: s.OfferingTokenAddress,
		MouseTierAddress: s.MouseTierAddress, StartBlock: s.StartBlock, EndBlock: s.EndBlock,
		VestingStartTime: s.VestingStartTime, VestingRevoked: s.VestingRevoked,
		MaxBufferBlocks: s.MaxBufferBlocks, CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt,
		DeletedAt: s.DeletedAt,
	}
}

func saleMetaPOToEntity(po *saleMetaPO) *domain.SaleMeta {
	return domain.ReconstructSaleMeta(
		po.ID, po.SaleID, strVal(po.Title), strVal(po.Description),
		strVal(po.BannerURL), strVal(po.LogoURL), strVal(po.WebsiteURL),
		strVal(po.SocialLinks), po.Visibility, po.SortOrder,
		po.CreatedAt, po.UpdatedAt,
	)
}

func saleMetaEntityToPO(m *domain.SaleMeta) *saleMetaPO {
	return &saleMetaPO{
		ID: m.ID, SaleID: m.SaleID, Title: strPtr(m.Title),
		Description: strPtr(m.Description), BannerURL: strPtr(m.BannerURL),
		LogoURL: strPtr(m.LogoURL), WebsiteURL: strPtr(m.WebsiteURL),
		SocialLinks: strPtr(m.SocialLinks), Visibility: m.Visibility,
		SortOrder: m.SortOrder, CreatedAt: m.CreatedAt, UpdatedAt: m.UpdatedAt,
	}
}

func poolPOToEntity(po *poolPO) *domain.Pool {
	return domain.ReconstructPool(
		po.ID, po.SaleID, po.PoolIndex,
		bigToInt(po.RaisingAmount), bigToInt(po.OfferingAmount),
		bigToInt(po.LimitPerUser), bigToInt(po.TaxRate),
		bigToInt(po.TotalAmount), bigToInt(po.TotalTax),
		po.IsSpecialSale, po.HasTax, po.VestingPercentage,
		po.VestingCliff, po.VestingDuration, po.VestingSlicePeriod,
		po.CreatedAt, po.UpdatedAt,
	)
}

func poolEntityToPO(p *domain.Pool) *poolPO {
	return &poolPO{
		ID: p.ID, SaleID: p.SaleID, PoolIndex: p.PoolIndex,
		RaisingAmount: intToBig(p.RaisingAmount), OfferingAmount: intToBig(p.OfferingAmount),
		LimitPerUser: intToBig(p.LimitPerUser), IsSpecialSale: p.IsSpecialSale,
		HasTax: p.HasTax, TaxRate: intToBig(p.TaxRate),
		VestingPercentage: p.VestingPercentage, VestingCliff: p.VestingCliff,
		VestingDuration: p.VestingDuration, VestingSlicePeriod: p.VestingSlicePeriod,
		TotalAmount: intToBig(p.TotalAmount), TotalTax: intToBig(p.TotalTax),
		CreatedAt: p.CreatedAt, UpdatedAt: p.UpdatedAt,
	}
}

func tierLimitPOToEntity(po *tierLimitPO) *domain.TierLimit {
	return domain.ReconstructTierLimit(po.ID, po.SaleID, po.Tier, bigToInt(po.CreditLimit), po.CreatedAt, po.UpdatedAt)
}

func tierLimitEntityToPO(t *domain.TierLimit) *tierLimitPO {
	return &tierLimitPO{ID: t.ID, SaleID: t.SaleID, Tier: t.Tier, CreditLimit: intToBig(t.CreditLimit), CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt}
}

func depositPOToEntity(po *depositPO) *domain.Deposit {
	return domain.ReconstructDeposit(po.ID, po.SaleID, po.PoolIndex, po.UserAddress, po.TxHash, bigToInt(po.Amount), po.BlockNumber, po.CreatedAt)
}

func userPoolStatePOToEntity(po *userPoolStatePO) *domain.UserPoolState {
	return domain.ReconstructUserPoolState(po.ID, po.SaleID, po.PoolIndex, po.UserAddress, bigToInt(po.TotalDeposited), po.Claimed, po.CreatedAt, po.UpdatedAt)
}

func userCreditPOToEntity(po *userCreditPO) *domain.UserCredit {
	return domain.ReconstructUserCredit(po.ID, po.SaleID, po.UserAddress, bigToInt(po.CreditUsed), po.CreatedAt, po.UpdatedAt)
}

func harvestPOToEntity(po *harvestPO) *domain.Harvest {
	return domain.ReconstructHarvest(
		po.ID, po.SaleID, po.PoolIndex, po.UserAddress, po.IsOverflow,
		bigToInt(po.OfferingAmount), bigToInt(po.PayAmount),
		bigToInt(po.RaiseRefund), bigToInt(po.TaxAmount),
		bigToInt(po.TGEAmount), bigToInt(po.VestingAmount),
		po.TxHash, po.BlockNumber, po.CreatedAt,
	)
}

func vestingSchedulePOToEntity(po *vestingSchedulePO) *domain.VestingSchedule {
	return domain.ReconstructVestingSchedule(
		po.ID, po.SaleID, po.PoolIndex, po.ScheduleID,
		po.Beneficiary, bigToInt(po.AmountTotal), bigToInt(po.Released),
		po.CreatedAt, po.UpdatedAt,
	)
}

func vestingReleasePOToEntity(po *vestingReleasePO) *domain.VestingRelease {
	return domain.ReconstructVestingRelease(po.ID, po.ScheduleID, bigToInt(po.Amount), po.TxHash, po.BlockNumber, po.CreatedAt)
}

func tokenPOToEntity(po *tokenPO) *domain.Token {
	return domain.ReconstructToken(po.ID, po.Address, po.ChainID, po.Name, po.Symbol, po.Decimals, strVal(po.LogoURL), po.CreatedAt, po.UpdatedAt)
}

func tokenEntityToPO(t *domain.Token) *tokenPO {
	return &tokenPO{
		ID: t.ID, Address: t.Address, ChainID: t.ChainID,
		Name: t.Name, Symbol: t.Symbol, Decimals: t.Decimals,
		LogoURL: strPtr(t.LogoURL), CreatedAt: t.CreatedAt, UpdatedAt: t.UpdatedAt,
	}
}

func prepareTxPOToEntity(po *prepareTxPO) *domain.PrepareTx {
	return domain.ReconstructPrepareTx(
		po.ID, po.SaleID, po.PoolIndex,
		domain.PrepareTxOperationType(po.OperationType),
		po.CallerAddress, po.Calldata, po.CalldataHash,
		domain.PrepareTxStatus(po.Status), po.TxHash, po.BlockNumber,
		po.ErrorMessage, po.ExpiresAt, po.ConfirmedAt,
		po.CreatedAt, po.UpdatedAt,
	)
}

func prepareTxEntityToPO(tx *domain.PrepareTx) *prepareTxPO {
	return &prepareTxPO{
		ID: tx.ID, SaleID: tx.SaleID, PoolIndex: tx.PoolIndex,
		OperationType: string(tx.OperationType), CallerAddress: tx.CallerAddress,
		Calldata: tx.Calldata, CalldataHash: tx.CalldataHash,
		Status: string(tx.Status), TxHash: tx.TxHash, BlockNumber: tx.BlockNumber,
		ErrorMessage: tx.ErrorMessage, ExpiresAt: tx.ExpiresAt,
		ConfirmedAt: tx.ConfirmedAt, CreatedAt: tx.CreatedAt, UpdatedAt: tx.UpdatedAt,
	}
}
