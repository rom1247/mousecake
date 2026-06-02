package domain

import "context"

// SwapRecordRepository 定义 SwapRecord 的持久化接口。
type SwapRecordRepository interface {
	// Save 保存 SwapRecord。
	Save(ctx context.Context, record *SwapRecord) error
	// FindByID 根据 ID 查找 SwapRecord。
	FindByID(ctx context.Context, id int64) (*SwapRecord, error)
	// UpdateStatus 原子更新状态，使用 WHERE status='pending' 防止并发重复提交。
	UpdateStatus(ctx context.Context, id int64, txHash string) error
}
