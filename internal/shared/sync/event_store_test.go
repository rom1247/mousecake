package sync

import (
	"context"
	"crypto/rand"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// TestChainEvent_TableName 测试 ChainEvent 表名映射。
func TestChainEvent_TableName(t *testing.T) {
	t.Parallel()

	var event ChainEvent
	assert.Equal(t, "chain_events", event.TableName())
}

// TestEventStatus_Constants 测试事件状态常量值正确。
func TestEventStatus_Constants(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status EventStatus
		want   string
	}{
		{"待处理", StatusPending, "pending"},
		{"处理中", StatusProcessing, "processing"},
		{"已处理", StatusProcessed, "processed"},
		{"处理失败", StatusFailed, "failed"},
		{"死信", StatusDeadLetter, "dead_letter"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.want, string(tt.status))
		})
	}
}

// TestChainEvent_Fields 测试 ChainEvent 结构体字段默认值。
func TestChainEvent_Fields(t *testing.T) {
	t.Parallel()

	event := ChainEvent{
		ChainID:     1,
		BlockNumber: 100,
		TxHash:      "0xabc",
	}

	assert.Equal(t, int64(0), event.ID)
	assert.Equal(t, 1, event.ChainID)
	assert.Equal(t, int64(100), event.BlockNumber)
	assert.Equal(t, "0xabc", event.TxHash)
	assert.Equal(t, EventStatus(""), event.Status) // 默认零值
	assert.Equal(t, 0, event.RetryCount)
	assert.Nil(t, event.ErrorMessage)
	assert.Equal(t, "", event.ProcessorID)
}

// TestNewEventStore 测试创建 EventStore 实例。
func TestNewEventStore(t *testing.T) {
	// 不使用 t.Parallel()，因为 NewEventStore 接收 nil 仅验证构造逻辑
	store := NewEventStore(nil)
	assert.NotNil(t, store)
}

// setupTestDB 创建 SQLite 内存数据库并自动迁移所需表，测试结束时自动关闭。
// 使用 cache=shared 确保连接池中的所有连接共享同一个数据库实例，
// 避免并发 goroutine 中出现 "no such table" 错误。
func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	var buf [8]byte
	_, _ = rand.Read(buf[:])
	dsn := fmt.Sprintf("file:%x?mode=memory&cache=shared", buf)

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	err = db.AutoMigrate(&ChainEvent{}, &Checkpoint{})
	require.NoError(t, err)
	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		sqlDB.Close()
	})
	return db
}

// TestEventStore_ClaimPending_ConcurrentRace 测试 ClaimPending 的行级乐观锁并发竞争。
// 两个 goroutine 同时竞争同一个 pending 事件，只有一个应成功获取。
func TestEventStore_ClaimPending_ConcurrentRace(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	store := NewEventStore(db)
	ctx := context.Background()

	// 插入一个 pending 状态的事件
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
	require.NotZero(t, event.ID)

	// 使用原子计数器记录成功获取的数量
	var successCount int32

	// 启动两个 goroutine 并发调用 ClaimPending
	var wg sync.WaitGroup
	wg.Add(2)

	for i := 0; i < 2; i++ {
		go func(idx int) {
			defer wg.Done()
			claimed, err := store.ClaimPending(ctx, event.ID, fmt.Sprintf("processor-%d", idx))
			if err != nil {
				t.Errorf("ClaimPending 返回错误: %v", err)
				return
			}
			if claimed {
				atomic.AddInt32(&successCount, 1)
			}
		}(i)
	}

	wg.Wait()

	// 验证只有一个 goroutine 成功获取事件
	assert.Equal(t, int32(1), atomic.LoadInt32(&successCount),
		"并发竞争中应只有一个 ClaimPending 返回 true")

	// 验证事件状态变为 processing
	updated, err := store.GetByID(ctx, event.ID)
	require.NoError(t, err)
	assert.Equal(t, StatusProcessing, updated.Status)
}
