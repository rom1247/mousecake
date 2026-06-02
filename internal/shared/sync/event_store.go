// Package sync 提供链上链下同步框架，支持历史回填、实时订阅、异步投影。
package sync

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// EventStatus 事件状态类型。
type EventStatus string

const (
	// StatusPending 待处理。
	StatusPending EventStatus = "pending"
	// StatusProcessing 处理中。
	StatusProcessing EventStatus = "processing"
	// StatusProcessed 已处理。
	StatusProcessed EventStatus = "processed"
	// StatusFailed 处理失败。
	StatusFailed EventStatus = "failed"
	// StatusDeadLetter 死信。
	StatusDeadLetter EventStatus = "dead_letter"
)

// ChainEvent 表示 chain_events 表的一条记录。
type ChainEvent struct {
	ID              int64       `gorm:"primaryKey;autoIncrement"`
	ChainID         int         `gorm:"column:chain_id;not null"`
	BlockNumber     int64       `gorm:"column:block_number;not null"`
	TxHash          string      `gorm:"column:tx_hash;type:char(66);not null"`
	TxIndex         int         `gorm:"column:tx_index;not null"`
	LogIndex        int         `gorm:"column:log_index;not null"`
	ContractAddress string      `gorm:"column:contract_address;type:char(42);not null"`
	EventName       string      `gorm:"column:event_name;type:varchar(64);not null"`
	EventData       string      `gorm:"column:event_data;type:text;not null"`
	Status          EventStatus `gorm:"column:status;type:varchar(16);not null;default:pending"`
	RetryCount      int         `gorm:"column:retry_count;not null;default:0"`
	ErrorMessage    *string     `gorm:"column:error_message"`
	ProcessorID     string      `gorm:"column:processor_id;type:varchar(64);not null;default:''"`
	LastFailedAt    *time.Time  `gorm:"column:last_failed_at"`
	CreatedAt       time.Time   `gorm:"column:created_at;not null;autoCreateTime"`
	UpdatedAt       time.Time   `gorm:"column:updated_at;not null;autoUpdateTime"`
}

// TableName 指定表名。
func (ChainEvent) TableName() string { return "chain_events" }

// EventStore 封装 chain_events 表的读写操作。
type EventStore struct {
	db *gorm.DB
}

// NewEventStore 创建 EventStore 实例。
func NewEventStore(db *gorm.DB) *EventStore {
	return &EventStore{db: db}
}

// BatchInsert 批量插入链上事件，忽略重复（唯一约束冲突时跳过）。
func (s *EventStore) BatchInsert(ctx context.Context, events []ChainEvent) (int64, error) {
	if len(events) == 0 {
		return 0, nil
	}
	result := s.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(events, 500)
	if result.Error != nil {
		return 0, fmt.Errorf("批量插入链上事件: %w", result.Error)
	}
	return result.RowsAffected, nil
}

// ClaimPending 使用行级乐观锁获取待处理事件：将 status 从 pending 改为 processing。
// 返回被成功获取的事件 ID，如果 affected_rows=0 表示被其他 Projector 抢先获取。
func (s *EventStore) ClaimPending(ctx context.Context, id int64, processorID string) (bool, error) {
	result := s.db.WithContext(ctx).
		Model(&ChainEvent{}).
		Where("id = ? AND status = ?", id, StatusPending).
		Updates(map[string]any{
			"status":       StatusProcessing,
			"processor_id": processorID,
			"updated_at":   time.Now(),
		})
	if result.Error != nil {
		return false, fmt.Errorf("获取事件 %d: %w", id, result.Error)
	}
	return result.RowsAffected > 0, nil
}

// MarkProcessed 标记事件为已处理。
func (s *EventStore) MarkProcessed(ctx context.Context, id int64) error {
	result := s.db.WithContext(ctx).
		Model(&ChainEvent{}).
		Where("id = ? AND status = ?", id, StatusProcessing).
		Updates(map[string]any{
			"status":     StatusProcessed,
			"updated_at": time.Now(),
		})
	if result.Error != nil {
		return fmt.Errorf("标记事件 %d 为已处理: %w", id, result.Error)
	}
	return nil
}

// MarkFailed 标记事件为失败，增加重试计数。
func (s *EventStore) MarkFailed(ctx context.Context, id int64, errMsg string) error {
	now := time.Now()
	result := s.db.WithContext(ctx).
		Model(&ChainEvent{}).
		Where("id = ? AND status = ?", id, StatusProcessing).
		Updates(map[string]any{
			"status":         StatusPending,
			"retry_count":    gorm.Expr("retry_count + 1"),
			"error_message":  errMsg,
			"last_failed_at": now,
			"updated_at":     now,
		})
	if result.Error != nil {
		return fmt.Errorf("标记事件 %d 为失败: %w", id, result.Error)
	}
	return nil
}

// MarkDeadLetter 标记事件为死信。
func (s *EventStore) MarkDeadLetter(ctx context.Context, id int64, errMsg string) error {
	now := time.Now()
	result := s.db.WithContext(ctx).
		Model(&ChainEvent{}).
		Where("id = ? AND status = ?", id, StatusProcessing).
		Updates(map[string]any{
			"status":         StatusDeadLetter,
			"error_message":  errMsg,
			"last_failed_at": now,
			"updated_at":     now,
		})
	if result.Error != nil {
		return fmt.Errorf("标记事件 %d 为死信: %w", id, result.Error)
	}
	return nil
}

// ListPending 查询待处理事件列表。
func (s *EventStore) ListPending(ctx context.Context, chainID int, processorID string, limit int) ([]ChainEvent, error) {
	var events []ChainEvent
	err := s.db.WithContext(ctx).
		Where("status = ? AND chain_id = ? AND processor_id = ?", StatusPending, chainID, processorID).
		Order("block_number ASC, tx_index ASC, log_index ASC").
		Limit(limit).
		Find(&events).Error
	if err != nil {
		return nil, fmt.Errorf("查询待处理事件: %w", err)
	}
	return events, nil
}

// ListFailed 查询失败/死信事件列表。
func (s *EventStore) ListFailed(ctx context.Context, chainID int, status EventStatus, page, pageSize int) ([]ChainEvent, int64, error) {
	var events []ChainEvent
	var total int64

	query := s.db.WithContext(ctx).Model(&ChainEvent{}).
		Where("status = ? AND chain_id = ?", status, chainID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计事件数量: %w", err)
	}

	offset := (page - 1) * pageSize
	err := query.Order("id DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&events).Error
	if err != nil {
		return nil, 0, fmt.Errorf("查询失败事件: %w", err)
	}
	return events, total, nil
}

// ResetToPending 将指定事件重置为 pending（管理员重试用）。
func (s *EventStore) ResetToPending(ctx context.Context, id int64) error {
	result := s.db.WithContext(ctx).
		Model(&ChainEvent{}).
		Where("id = ? AND status IN ?", id, []EventStatus{StatusFailed, StatusDeadLetter}).
		Updates(map[string]any{
			"status":        StatusPending,
			"retry_count":   0,
			"error_message": nil,
			"updated_at":    time.Now(),
		})
	if result.Error != nil {
		return fmt.Errorf("重置事件 %d: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("事件 %d 不存在或状态不允许重置", id)
	}
	return nil
}

// ResetBlockRange 批量重置区块范围内的事件为 pending。
func (s *EventStore) ResetBlockRange(ctx context.Context, chainID int, fromBlock, toBlock int64) (int64, error) {
	result := s.db.WithContext(ctx).
		Model(&ChainEvent{}).
		Where("chain_id = ? AND block_number BETWEEN ? AND ? AND status IN ?",
			chainID, fromBlock, toBlock, []EventStatus{StatusFailed, StatusDeadLetter}).
		Updates(map[string]any{
			"status":        StatusPending,
			"retry_count":   0,
			"error_message": nil,
			"updated_at":    time.Now(),
		})
	if result.Error != nil {
		return 0, fmt.Errorf("重置区块范围事件: %w", result.Error)
	}
	return result.RowsAffected, nil
}

// ResetProcessingTimeout 重置超时的 processing 事件为 pending。
func (s *EventStore) ResetProcessingTimeout(ctx context.Context, timeout time.Duration) (int64, error) {
	cutoff := time.Now().Add(-timeout)
	result := s.db.WithContext(ctx).
		Model(&ChainEvent{}).
		Where("status = ? AND updated_at < ?", StatusProcessing, cutoff).
		Updates(map[string]any{
			"status":     StatusPending,
			"updated_at": time.Now(),
		})
	if result.Error != nil {
		return 0, fmt.Errorf("重置超时 processing 事件: %w", result.Error)
	}
	return result.RowsAffected, nil
}

// CountByStatus 按状态统计事件数量。
func (s *EventStore) CountByStatus(ctx context.Context, chainID int) (map[EventStatus]int64, error) {
	type statusCount struct {
		Status EventStatus `gorm:"column:status"`
		Count  int64       `gorm:"column:count"`
	}

	var counts []statusCount
	err := s.db.WithContext(ctx).
		Model(&ChainEvent{}).
		Select("status, COUNT(*) as count").
		Where("chain_id = ?", chainID).
		Group("status").
		Find(&counts).Error
	if err != nil {
		return nil, fmt.Errorf("统计事件状态: %w", err)
	}

	result := make(map[EventStatus]int64)
	for _, c := range counts {
		result[c.Status] = c.Count
	}
	return result, nil
}

// GetByID 按 ID 查询事件。
func (s *EventStore) GetByID(ctx context.Context, id int64) (*ChainEvent, error) {
	var event ChainEvent
	err := s.db.WithContext(ctx).First(&event, id).Error
	if err != nil {
		return nil, fmt.Errorf("查询事件 %d: %w", id, err)
	}
	return &event, nil
}
