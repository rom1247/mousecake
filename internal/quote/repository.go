package quote

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"

	"github.com/mousecake-go/mousecake-go/internal/quote/domain"
)

// swapRecordPO SwapRecord 持久化对象。
type swapRecordPO struct {
	ID              int64     `gorm:"column:id;primaryKey"`
	Provider        string    `gorm:"column:provider"`
	ChainID         int       `gorm:"column:chain_id"`
	FromToken       string    `gorm:"column:from_token"`
	ToToken         string    `gorm:"column:to_token"`
	FromAmount      string    `gorm:"column:from_amount"`
	ToAmount        string    `gorm:"column:to_amount"`
	SlippagePercent float64   `gorm:"column:slippage_percent"`
	SwapMode        string    `gorm:"column:swap_mode"`
	Status          string    `gorm:"column:status"`
	TxHash          string    `gorm:"column:tx_hash"`
	CreatedAt       time.Time `gorm:"column:created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at"`
}

func (swapRecordPO) TableName() string { return "swap_records" }

// SwapRecordRepository 实现 domain.SwapRecordRepository 接口。
type SwapRecordRepository struct {
	db  *gorm.DB
	log *slog.Logger
}

// NewSwapRecordRepository 创建 SwapRecord 仓库实现。
func NewSwapRecordRepository(db *gorm.DB) *SwapRecordRepository {
	return &SwapRecordRepository{
		db:  db,
		log: slog.Default().With("module", "quote", "layer", "repository"),
	}
}

// Save 保存 SwapRecord。
func (r *SwapRecordRepository) Save(ctx context.Context, record *domain.SwapRecord) error {
	po := toSwapRecordPO(record)
	if err := r.db.WithContext(ctx).Create(&po).Error; err != nil {
		return fmt.Errorf("save swap record: %w", err)
	}
	return nil
}

// FindByID 根据 ID 查找 SwapRecord。
func (r *SwapRecordRepository) FindByID(ctx context.Context, id int64) (*domain.SwapRecord, error) {
	var po swapRecordPO
	if err := r.db.WithContext(ctx).
		Select("id, provider, chain_id, from_token, to_token, from_amount, to_amount, slippage_percent, swap_mode, status, tx_hash, created_at, updated_at").
		Where("id = ?", id).
		First(&po).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrSwapRecordNotFound
		}
		return nil, fmt.Errorf("find swap record by id %d: %w", id, err)
	}
	return toSwapRecordEntity(&po), nil
}

// UpdateStatus 原子更新状态，使用 WHERE status='pending' 防止并发重复提交。
func (r *SwapRecordRepository) UpdateStatus(ctx context.Context, id int64, txHash string) error {
	result := r.db.WithContext(ctx).
		Model(&swapRecordPO{}).
		Where("id = ? AND status = ?", id, string(domain.SwapStatusPending)).
		Updates(map[string]interface{}{
			"status":     string(domain.SwapStatusSubmitted),
			"tx_hash":    txHash,
			"updated_at": time.Now(),
		})
	if result.Error != nil {
		return fmt.Errorf("update swap record status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return domain.ErrAlreadySubmitted
	}
	return nil
}

// toSwapRecordPO 将领域模型转为 PO。
func toSwapRecordPO(record *domain.SwapRecord) swapRecordPO {
	return swapRecordPO{
		ID:              record.ID,
		Provider:        record.Provider,
		ChainID:         record.ChainID,
		FromToken:       record.FromToken,
		ToToken:         record.ToToken,
		FromAmount:      record.FromAmount,
		ToAmount:        record.ToAmount,
		SlippagePercent: record.SlippagePercent,
		SwapMode:        string(record.SwapMode),
		Status:          string(record.Status),
		TxHash:          record.TxHash,
		CreatedAt:       record.CreatedAt,
		UpdatedAt:       record.UpdatedAt,
	}
}

// toSwapRecordEntity 将 PO 转为领域模型。
func toSwapRecordEntity(po *swapRecordPO) *domain.SwapRecord {
	return domain.ReconstructSwapRecord(domain.SwapRecordSnapshot{
		ID:              po.ID,
		Provider:        po.Provider,
		ChainID:         po.ChainID,
		FromToken:       po.FromToken,
		ToToken:         po.ToToken,
		FromAmount:      po.FromAmount,
		ToAmount:        po.ToAmount,
		SlippagePercent: po.SlippagePercent,
		SwapMode:        domain.SwapMode(po.SwapMode),
		Status:          domain.SwapStatus(po.Status),
		TxHash:          po.TxHash,
	})
}
