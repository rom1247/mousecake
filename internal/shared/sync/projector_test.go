package sync

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mousecake-go/mousecake-go/config"
)

// TestNewProjector 测试创建 Projector 实例。
func TestNewProjector(t *testing.T) {
	t.Parallel()

	cfg := config.ProjectorConfig{
		MaxWorkers:        4,
		MaxRetries:        3,
		ProcessingTimeout: 30000000000, // 30s
		PollInterval:      1000000000,  // 1s
	}
	metrics := newSyncMetrics(1)

	p := NewProjector(nil, nil, 1, cfg, metrics)

	assert.NotNil(t, p)
	assert.Equal(t, 1, p.chainID)
	assert.Equal(t, 4, p.cfg.MaxWorkers)
	assert.Equal(t, 3, p.cfg.MaxRetries)
	assert.NotNil(t, p.partitions)
	assert.Empty(t, p.partitions)
}

// TestNewProjector_DefaultConfig 测试使用零值配置创建 Projector。
func TestNewProjector_DefaultConfig(t *testing.T) {
	t.Parallel()

	p := NewProjector(nil, nil, 11155111, config.ProjectorConfig{}, nil)

	assert.NotNil(t, p)
	assert.Equal(t, 11155111, p.chainID)
	assert.Nil(t, p.metrics)
	assert.Equal(t, 0, p.cfg.MaxWorkers)
	assert.Equal(t, 0, p.cfg.MaxRetries)
}

// mockEventService 用于测试的事件处理 mock。
type mockEventService struct {
	handled []*ChainEvent
	err     error
}

// HandleEvent 实现 EventService 接口，记录处理的事件。
func (m *mockEventService) HandleEvent(_ context.Context, event *ChainEvent) error {
	m.handled = append(m.handled, event)
	return m.err
}

// TestProjector_ProcessEvent_RetryOnFailure 测试投影失败后事件状态重置并增加重试计数。
func TestProjector_ProcessEvent_RetryOnFailure(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	store := NewEventStore(db)
	ctx := context.Background()

	// 插入一个 pending 事件
	event := ChainEvent{
		ChainID:         1,
		BlockNumber:     100,
		TxHash:          "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
		TxIndex:         0,
		LogIndex:        0,
		ContractAddress: "0x1234567890abcdef1234567890abcdef12345678",
		EventName:       "Transfer",
		EventData:       "{}",
		Status:          StatusPending,
		ProcessorID:     "",
	}
	result := db.WithContext(ctx).Create(&event)
	require.NoError(t, result.Error)

	// 手动将事件设为 processing，模拟 ClaimPending 成功后的状态
	err := db.WithContext(ctx).Model(&ChainEvent{}).
		Where("id = ?", event.ID).
		Updates(map[string]any{"status": StatusProcessing, "updated_at": time.Now()}).Error
	require.NoError(t, err)

	cfg := config.ProjectorConfig{
		MaxWorkers:        4,
		MaxRetries:        3,
		ProcessingTimeout: 5 * time.Minute,
		PollInterval:      1 * time.Second,
	}
	metrics := newSyncMetrics(1)

	svc := &mockEventService{err: fmt.Errorf("处理失败")}
	p := NewProjector(store, svc, 1, cfg, metrics)

	// 调用 processEvent
	p.processEvent(ctx, &event)

	// 验证事件状态变为 pending（MarkFailed 重置），retry_count 增加
	updated, err := store.GetByID(ctx, event.ID)
	require.NoError(t, err)
	assert.Equal(t, StatusPending, updated.Status, "失败后事件应重置为 pending")
	assert.Equal(t, 1, updated.RetryCount, "retry_count 应增加 1")
	assert.NotNil(t, updated.ErrorMessage)
	assert.Equal(t, "处理失败", *updated.ErrorMessage)
}

// TestProjector_Dispatch_BlockPartition 测试区块级分区并发分发。
// 验证 dispatch 方法为不同区块号创建不同分区，同区块号复用同一分区。
func TestProjector_Dispatch_BlockPartition(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	store := NewEventStore(db)

	cfg := config.ProjectorConfig{
		MaxWorkers:        4,
		MaxRetries:        3,
		ProcessingTimeout: 5 * time.Minute,
		PollInterval:      1 * time.Second,
	}
	metrics := newSyncMetrics(1)

	svc := &mockEventService{}
	p := NewProjector(store, svc, 1, cfg, metrics)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 分发不同区块号的事件
	event1 := &ChainEvent{ID: 1, BlockNumber: 100}
	event2 := &ChainEvent{ID: 2, BlockNumber: 200}
	event3 := &ChainEvent{ID: 3, BlockNumber: 100} // 与 event1 同区块

	p.dispatch(ctx, event1)
	p.dispatch(ctx, event2)
	p.dispatch(ctx, event3)

	// 验证分区数量：两个不同区块号应有 2 个分区
	p.mu.Lock()
	partitionCount := len(p.partitions)
	p.mu.Unlock()

	assert.Equal(t, 2, partitionCount, "不同区块号应创建不同分区")

	// 清理分区
	p.mu.Lock()
	for _, ch := range p.partitions {
		close(ch)
	}
	p.partitions = make(map[int64]chan *ChainEvent)
	p.mu.Unlock()
}

// TestProjector_ProcessEvent_DeadLetter 测试超过最大重试次数后事件进入死信队列。
func TestProjector_ProcessEvent_DeadLetter(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	store := NewEventStore(db)
	ctx := context.Background()

	// 插入一个 pending 事件
	event := ChainEvent{
		ChainID:         1,
		BlockNumber:     100,
		TxHash:          "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
		TxIndex:         0,
		LogIndex:        0,
		ContractAddress: "0x1234567890abcdef1234567890abcdef12345678",
		EventName:       "Transfer",
		EventData:       "{}",
		Status:          StatusPending,
		ProcessorID:     "",
		RetryCount:      2, // 已重试 2 次，maxRetries=3 时下一次将达到上限
	}
	result := db.WithContext(ctx).Create(&event)
	require.NoError(t, result.Error)

	// 手动将事件设为 processing
	err := db.WithContext(ctx).Model(&ChainEvent{}).
		Where("id = ?", event.ID).
		Updates(map[string]any{"status": StatusProcessing, "updated_at": time.Now()}).Error
	require.NoError(t, err)

	cfg := config.ProjectorConfig{
		MaxWorkers:        4,
		MaxRetries:        3,
		ProcessingTimeout: 5 * time.Minute,
		PollInterval:      1 * time.Second,
	}
	metrics := newSyncMetrics(1)

	svc := &mockEventService{err: fmt.Errorf("处理失败")}
	p := NewProjector(store, svc, 1, cfg, metrics)

	// 调用 processEvent，RetryCount(2)+1 >= MaxRetries(3)，应进入死信
	p.processEvent(ctx, &event)

	// 验证事件状态变为 dead_letter
	updated, err := store.GetByID(ctx, event.ID)
	require.NoError(t, err)
	assert.Equal(t, StatusDeadLetter, updated.Status, "超过最大重试次数应变为 dead_letter")
}

// TestEventStore_ResetProcessingTimeout 测试 processing 超时事件重置为 pending。
func TestEventStore_ResetProcessingTimeout(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	store := NewEventStore(db)
	ctx := context.Background()

	// 插入一个 processing 状态事件，updated_at 设为很久以前
	event := ChainEvent{
		ChainID:         1,
		BlockNumber:     100,
		TxHash:          "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
		TxIndex:         0,
		LogIndex:        0,
		ContractAddress: "0x1234567890abcdef1234567890abcdef12345678",
		EventName:       "Transfer",
		EventData:       "{}",
		Status:          StatusProcessing,
		ProcessorID:     "test-processor",
	}
	result := db.WithContext(ctx).Create(&event)
	require.NoError(t, result.Error)

	// 手动将 updated_at 设为 1 小时前（超过任何合理的超时时间）
	pastTime := time.Now().Add(-1 * time.Hour)
	err := db.WithContext(ctx).Model(&ChainEvent{}).
		Where("id = ?", event.ID).
		Updates(map[string]any{"updated_at": pastTime}).Error
	require.NoError(t, err)

	// 使用 5 分钟的超时阈值，1 小时前的事件应被重置
	reset, err := store.ResetProcessingTimeout(ctx, 5*time.Minute)
	require.NoError(t, err)
	assert.Equal(t, int64(1), reset, "应重置 1 个超时事件")

	// 验证事件状态重置为 pending
	updated, err := store.GetByID(ctx, event.ID)
	require.NoError(t, err)
	assert.Equal(t, StatusPending, updated.Status, "超时事件应重置为 pending")
}

// TestProjector_ProcessEvent_Idempotent 测试投影成功后事件状态变为 processed。
func TestProjector_ProcessEvent_Idempotent(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	store := NewEventStore(db)
	ctx := context.Background()

	// 插入一个 pending 事件
	event := ChainEvent{
		ChainID:         1,
		BlockNumber:     100,
		TxHash:          "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
		TxIndex:         0,
		LogIndex:        0,
		ContractAddress: "0x1234567890abcdef1234567890abcdef12345678",
		EventName:       "Transfer",
		EventData:       "{}",
		Status:          StatusPending,
		ProcessorID:     "",
	}
	result := db.WithContext(ctx).Create(&event)
	require.NoError(t, result.Error)

	// 手动将事件设为 processing
	err := db.WithContext(ctx).Model(&ChainEvent{}).
		Where("id = ?", event.ID).
		Updates(map[string]any{"status": StatusProcessing, "updated_at": time.Now()}).Error
	require.NoError(t, err)

	cfg := config.ProjectorConfig{
		MaxWorkers:        4,
		MaxRetries:        3,
		ProcessingTimeout: 5 * time.Minute,
		PollInterval:      1 * time.Second,
	}
	metrics := newSyncMetrics(1)

	svc := &mockEventService{} // err == nil 表示处理成功
	p := NewProjector(store, svc, 1, cfg, metrics)

	// 调用 processEvent
	p.processEvent(ctx, &event)

	// 验证事件状态变为 processed
	updated, err := store.GetByID(ctx, event.ID)
	require.NoError(t, err)
	assert.Equal(t, StatusProcessed, updated.Status, "成功处理后事件应变为 processed")

	// 验证 EventService 确实被调用了
	assert.Len(t, svc.handled, 1)
	assert.Equal(t, event.ID, svc.handled[0].ID)
}
