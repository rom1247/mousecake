package sync

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Checkpoint 表示 sync_checkpoints 表的一条记录。
type Checkpoint struct {
	ID              int64     `gorm:"primaryKey;autoIncrement"`
	ChainID         int       `gorm:"column:chain_id;not null;uniqueIndex:idx_chain_processor"`
	ProcessorID     string    `gorm:"column:processor_id;type:varchar(64);not null;uniqueIndex:idx_chain_processor"`
	LastSyncedBlock int64     `gorm:"column:last_synced_block;not null;default:0"`
	UpdatedAt       time.Time `gorm:"column:updated_at;not null;autoUpdateTime"`
}

// TableName 指定表名。
func (Checkpoint) TableName() string { return "sync_checkpoints" }

// CheckpointRepository 封装 sync_checkpoints 表操作。
type CheckpointRepository struct {
	db *gorm.DB
}

// NewCheckpointRepository 创建 CheckpointRepository 实例。
func NewCheckpointRepository(db *gorm.DB) *CheckpointRepository {
	return &CheckpointRepository{db: db}
}

// Upsert 写入或更新检查点。
func (r *CheckpointRepository) Upsert(ctx context.Context, chainID int, processorID string, lastSyncedBlock int64) error {
	cp := Checkpoint{
		ChainID:         chainID,
		ProcessorID:     processorID,
		LastSyncedBlock: lastSyncedBlock,
		UpdatedAt:       time.Now(),
	}
	result := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "chain_id"}, {Name: "processor_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"last_synced_block", "updated_at"}),
		}).
		Create(&cp)
	if result.Error != nil {
		return fmt.Errorf("写入检查点 chain_id=%d processor_id=%s: %w", chainID, processorID, result.Error)
	}
	return nil
}

// Get 查询检查点。
func (r *CheckpointRepository) Get(ctx context.Context, chainID int, processorID string) (*Checkpoint, error) {
	var cp Checkpoint
	err := r.db.WithContext(ctx).
		Where("chain_id = ? AND processor_id = ?", chainID, processorID).
		First(&cp).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询检查点 chain_id=%d processor_id=%s: %w", chainID, processorID, err)
	}
	return &cp, nil
}
