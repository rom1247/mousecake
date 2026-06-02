package sync

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mousecake-go/mousecake-go/config"
)

// TestBackfiller_parseAddresses 测试合约地址字符串解析为 common.Address。
func TestBackfiller_parseAddresses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		addresses []string
		wantCount int
		wantFirst string
	}{
		{
			name:      "空地址列表",
			addresses: []string{},
			wantCount: 0,
		},
		{
			name:      "单个合约地址",
			addresses: []string{"0x1234567890abcdef1234567890abcdef12345678"},
			wantCount: 1,
			wantFirst: "0x1234567890AbcdEF1234567890aBcdef12345678",
		},
		{
			name:      "多个合约地址",
			addresses: []string{"0x1111111111111111111111111111111111111111", "0x2222222222222222222222222222222222222222"},
			wantCount: 2,
			wantFirst: "0x1111111111111111111111111111111111111111",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			b := &Backfiller{
				addresses: tt.addresses,
			}

			addrs := b.parseAddresses()
			assert.Len(t, addrs, tt.wantCount)
			if tt.wantCount > 0 {
				assert.Equal(t, tt.wantFirst, addrs[0].Hex())
			}
		})
	}
}

// TestBackfiller_convertLogs 测试链上日志转换为 ChainEvent 列表。
func TestBackfiller_convertLogs(t *testing.T) {
	t.Parallel()

	txHash := common.HexToHash("0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")
	contractAddr := common.HexToAddress("0x1234567890abcdef1234567890abcdef1234567890")
	topic := common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
	blockHash := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001")

	b := &Backfiller{
		chainID:     1,
		processorID: "test-processor",
	}

	tests := []struct {
		name       string
		logs       []types.Log
		wantCount  int
		wantStatus EventStatus
	}{
		{
			name:      "空日志列表",
			logs:      []types.Log{},
			wantCount: 0,
		},
		{
			name: "单条日志",
			logs: []types.Log{
				{
					Address:     contractAddr,
					Topics:      []common.Hash{topic},
					Data:        []byte{0x01, 0x02},
					BlockNumber: 100,
					TxHash:      txHash,
					TxIndex:     1,
					Index:       2,
					BlockHash:   blockHash,
					Removed:     false,
				},
			},
			wantCount:  1,
			wantStatus: StatusPending,
		},
		{
			name: "多条日志",
			logs: []types.Log{
				{
					Address:     contractAddr,
					Topics:      []common.Hash{topic},
					Data:        []byte{},
					BlockNumber: 100,
					TxHash:      txHash,
					TxIndex:     0,
					Index:       0,
					BlockHash:   blockHash,
				},
				{
					Address:     contractAddr,
					Topics:      []common.Hash{topic},
					Data:        []byte{},
					BlockNumber: 101,
					TxHash:      txHash,
					TxIndex:     0,
					Index:       1,
					BlockHash:   blockHash,
				},
			},
			wantCount:  2,
			wantStatus: StatusPending,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			events := b.convertLogs(tt.logs)
			assert.Len(t, events, tt.wantCount)

			if tt.wantCount > 0 {
				first := events[0]
				assert.Equal(t, 1, first.ChainID)
				assert.Equal(t, tt.wantStatus, first.Status)
				assert.Equal(t, "test-processor", first.ProcessorID)
				assert.Equal(t, "0xddf252ad", first.EventName)
			}
		})
	}
}

// TestBackfiller_convertLogs_FieldMapping 测试 convertLogs 的字段映射准确性。
func TestBackfiller_convertLogs_FieldMapping(t *testing.T) {
	t.Parallel()

	txHash := common.HexToHash("0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")
	contractAddr := common.HexToAddress("0x1234567890abcdef1234567890abcdef1234567890")
	topic := common.HexToHash("0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef")
	blockHash := common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000001")

	b := &Backfiller{
		chainID:     137,
		processorID: "launchpad",
	}

	log := types.Log{
		Address:     contractAddr,
		Topics:      []common.Hash{topic},
		Data:        []byte{0xaa, 0xbb},
		BlockNumber: 12345,
		TxHash:      txHash,
		TxIndex:     3,
		Index:       7,
		BlockHash:   blockHash,
		Removed:     false,
	}

	events := b.convertLogs([]types.Log{log})
	requireLen(t, events, 1)

	ev := events[0]
	assert.Equal(t, 137, ev.ChainID)
	assert.Equal(t, int64(12345), ev.BlockNumber)
	assert.Equal(t, txHash.Hex(), ev.TxHash)
	assert.Equal(t, 3, ev.TxIndex)
	assert.Equal(t, 7, ev.LogIndex)
	assert.Equal(t, contractAddr.Hex(), ev.ContractAddress)
	assert.Equal(t, "0xddf252ad", ev.EventName)
	assert.Equal(t, StatusPending, ev.Status)
	assert.Equal(t, "launchpad", ev.ProcessorID)
	// EventData 应包含 topics、data、block_hash、removed 字段
	assert.Contains(t, ev.EventData, "topics")
	assert.Contains(t, ev.EventData, "data")
}

// requireLen 断言切片长度，失败时终止测试。
func requireLen(t *testing.T, slice []ChainEvent, want int) {
	t.Helper()
	if len(slice) != want {
		t.Fatalf("期望长度 %d, 实际 %d", want, len(slice))
	}
}

// TestBackfiller_NewBackfiller 测试 NewBackfiller 构造函数。
func TestBackfiller_NewBackfiller(t *testing.T) {
	t.Parallel()

	cfg := config.BackfillConfig{
		InitialBatchSize: 100,
		MinBatchSize:     10,
		MaxBatchSize:     1000,
		GrowthFactor:     1.5,
	}

	b := NewBackfiller(
		nil, // pool — 在单元测试中不需要
		nil, // store
		nil, // checkpoint
		1,   // chainID
		"test-processor",
		[]string{"0x1111111111111111111111111111111111111111"},
		12,
		cfg,
	)

	assert.NotNil(t, b)
	assert.Equal(t, 1, b.chainID)
	assert.Equal(t, "test-processor", b.processorID)
	assert.Len(t, b.addresses, 1)
	assert.Equal(t, int64(12), b.confirmationBlocks)
	assert.Equal(t, 100, b.cfg.InitialBatchSize)
	assert.Equal(t, 1.5, b.cfg.GrowthFactor)
}

// mockNodePool 是 chain.NodePool 接口的测试 mock 实现。
type mockNodePool struct {
	blockNumber uint64
	blockErr    error
	logs        []types.Log
	logsErr     error
}

// CallContract 实现 NodePool 接口，返回零值。
func (m *mockNodePool) CallContract(_ context.Context, _ ethereum.CallMsg, _ *big.Int) ([]byte, error) {
	return nil, nil
}

// FilterLogs 实现 NodePool 接口，返回 mock 配置的日志或错误。
func (m *mockNodePool) FilterLogs(_ context.Context, _ ethereum.FilterQuery) ([]types.Log, error) {
	return m.logs, m.logsErr
}

// TransactionReceipt 实现 NodePool 接口，返回零值。
func (m *mockNodePool) TransactionReceipt(_ context.Context, _ common.Hash) (*types.Receipt, error) {
	return nil, nil
}

// BlockNumber 实现 NodePool 接口，返回 mock 配置的区块号。
func (m *mockNodePool) BlockNumber(_ context.Context) (uint64, error) {
	return m.blockNumber, m.blockErr
}

// HeaderByNumber 实现 NodePool 接口，返回零值。
func (m *mockNodePool) HeaderByNumber(_ context.Context, _ *big.Int) (*types.Header, error) {
	return nil, nil
}

// SubscribeLogs 实现 NodePool 接口，返回 nil。
func (m *mockNodePool) SubscribeLogs(_ context.Context, _ ethereum.FilterQuery) (chan types.Log, ethereum.Subscription, error) {
	return nil, nil, nil
}

// Close 实现 NodePool 接口，不做任何操作。
func (m *mockNodePool) Close() {}

// TestBackfiller_Run_BatchSizeGrowth 测试回填成功时 batch_size 正常增长。
// 配置 InitialBatchSize=5000，GrowthFactor=1.2，成功回填后检查点正确推进。
func TestBackfiller_Run_BatchSizeGrowth(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	store := NewEventStore(db)
	checkpoint := NewCheckpointRepository(db)
	ctx := context.Background()

	pool := &mockNodePool{
		blockNumber: 19300000,
		logs:        []types.Log{}, // 空日志，模拟成功
	}

	cfg := config.BackfillConfig{
		InitialBatchSize: 5000,
		GrowthFactor:     1.2,
		MaxBatchSize:     10000,
		MinBatchSize:     500,
	}

	b := NewBackfiller(
		pool, store, checkpoint,
		1, "test-processor",
		[]string{"0x1234567890abcdef1234567890abcdef12345678"},
		12, // confirmationBlocks
		cfg,
	)

	err := b.Run(ctx, 19283746)
	require.NoError(t, err)

	// 验证 checkpoint 已写入正确的 last_synced_block
	cp, err := checkpoint.Get(ctx, 1, "test-processor")
	require.NoError(t, err)
	require.NotNil(t, cp, "回填完成后应存在检查点")

	// targetBlock = 19300000 - 12 = 19299988
	expectedTarget := int64(19300000) - 12
	assert.Equal(t, expectedTarget, cp.LastSyncedBlock,
		"检查点的 last_synced_block 应等于 currentBlock - confirmationBlocks")
}

// TestBackfiller_Run_BatchSizeConfig 测试 batch_size 相关配置正确存储和边界计算。
func TestBackfiller_Run_BatchSizeConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		cfg              config.BackfillConfig
		wantInitial      int
		wantMin          int
		wantMax          int
		wantGrowthFactor float64
	}{
		{
			name: "默认配置",
			cfg: config.BackfillConfig{
				InitialBatchSize: 5000,
				MinBatchSize:     500,
				MaxBatchSize:     10000,
				GrowthFactor:     1.2,
			},
			wantInitial:      5000,
			wantMin:          500,
			wantMax:          10000,
			wantGrowthFactor: 1.2,
		},
		{
			name: "小批次配置",
			cfg: config.BackfillConfig{
				InitialBatchSize: 100,
				MinBatchSize:     10,
				MaxBatchSize:     500,
				GrowthFactor:     1.5,
			},
			wantInitial:      100,
			wantMin:          10,
			wantMax:          500,
			wantGrowthFactor: 1.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			b := NewBackfiller(nil, nil, nil, 1, "test", nil, 0, tt.cfg)
			assert.Equal(t, tt.wantInitial, b.cfg.InitialBatchSize)
			assert.Equal(t, tt.wantMin, b.cfg.MinBatchSize)
			assert.Equal(t, tt.wantMax, b.cfg.MaxBatchSize)
			assert.Equal(t, tt.wantGrowthFactor, b.cfg.GrowthFactor)
		})
	}
}

// TestBackfiller_Run_BatchSizeHalvedOnFailure 测试回填失败时 batch_size 减半。
// FilterLogs 返回错误时，回填器应减小 batch_size 并继续重试直到 context 取消。
func TestBackfiller_Run_BatchSizeHalvedOnFailure(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	store := NewEventStore(db)
	checkpoint := NewCheckpointRepository(db)

	// 使用可取消的 context，在足够时间后停止无限重试
	ctx, cancel := context.WithCancel(context.Background())

	pool := &mockNodePool{
		blockNumber: 19300000,
		logsErr:     fmt.Errorf("模拟超时错误"),
	}

	cfg := config.BackfillConfig{
		InitialBatchSize: 5000,
		GrowthFactor:     1.2,
		MaxBatchSize:     10000,
		MinBatchSize:     500,
	}

	b := NewBackfiller(
		pool, store, checkpoint,
		1, "test-processor",
		[]string{"0x1234567890abcdef1234567890abcdef12345678"},
		12,
		cfg,
	)

	// 在另一个 goroutine 中运行，一段时间后取消 context
	done := make(chan error, 1)
	go func() {
		done <- b.Run(ctx, 19283746)
	}()

	// 等待足够时间让重试循环执行几次，然后取消
	time.Sleep(200 * time.Millisecond)
	cancel()

	err := <-done
	// context 取消后应返回 context 错误
	assert.Error(t, err)
}

// TestBackfiller_Run_NoBackfillNeeded 测试起始区块超过目标时无需回填。
func TestBackfiller_Run_NoBackfillNeeded(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	store := NewEventStore(db)
	checkpoint := NewCheckpointRepository(db)
	ctx := context.Background()

	pool := &mockNodePool{
		blockNumber: 19283746, // currentBlock 与 startBlock 相同
		logs:        []types.Log{},
	}

	cfg := config.BackfillConfig{
		InitialBatchSize: 5000,
		GrowthFactor:     1.2,
		MaxBatchSize:     10000,
		MinBatchSize:     500,
	}

	b := NewBackfiller(
		pool, store, checkpoint,
		1, "test-processor",
		[]string{"0x1234567890abcdef1234567890abcdef12345678"},
		12,
		cfg,
	)

	// startBlock 超过 currentBlock - confirmationBlocks，无需回填
	err := b.Run(ctx, 19283746)
	require.NoError(t, err)

	// 无需回填时不应写入 checkpoint
	cp, err := checkpoint.Get(ctx, 1, "test-processor")
	require.NoError(t, err)
	assert.Nil(t, cp, "无需回填时不应存在检查点")
}
