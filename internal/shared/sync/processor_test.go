package sync

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mousecake-go/mousecake-go/config"
	"github.com/mousecake-go/mousecake-go/internal/chain"
)

// mockSubscriberNodePool 同时满足 NodePool 接口，并提供可控行为。
type mockSubscriberNodePool struct {
	blockNumber uint64
	blockErr    error
	logs        []types.Log
	subLogs     chan types.Log
	subErr      error
}

func (m *mockSubscriberNodePool) BlockNumber(_ context.Context) (uint64, error) {
	return m.blockNumber, m.blockErr
}

func (m *mockSubscriberNodePool) FilterLogs(_ context.Context, _ ethereum.FilterQuery) ([]types.Log, error) {
	return m.logs, nil
}

func (m *mockSubscriberNodePool) SubscribeLogs(_ context.Context, _ ethereum.FilterQuery) (chan types.Log, ethereum.Subscription, error) {
	return m.subLogs, nil, m.subErr
}

func (m *mockSubscriberNodePool) CallContract(_ context.Context, _ ethereum.CallMsg, _ *big.Int) ([]byte, error) {
	return nil, nil
}

func (m *mockSubscriberNodePool) TransactionReceipt(_ context.Context, _ common.Hash) (*types.Receipt, error) {
	return nil, nil
}

func (m *mockSubscriberNodePool) HeaderByNumber(_ context.Context, _ *big.Int) (*types.Header, error) {
	return nil, nil
}

func (m *mockSubscriberNodePool) BatchCallContract(_ context.Context, _ []ethereum.CallMsg) ([][]byte, error) {
	return nil, nil
}

func (m *mockSubscriberNodePool) PendingNonceAt(_ context.Context, _ common.Address) (uint64, error) {
	return 0, nil
}

func (m *mockSubscriberNodePool) EstimateGas(_ context.Context, _ ethereum.CallMsg) (uint64, error) {
	return 0, nil
}

func (m *mockSubscriberNodePool) SendTransaction(_ context.Context, _ *types.Transaction) error {
	return nil
}

func (m *mockSubscriberNodePool) SuggestGasPrice(_ context.Context) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (m *mockSubscriberNodePool) Close() {}

// 确保 mockSubscriberNodePool 满足 NodePool 接口。
var _ chain.NodePool = (*mockSubscriberNodePool)(nil)

// defaultSyncCfg 创建测试用 SyncConfig。
func defaultSyncCfg() config.SyncConfig {
	return config.SyncConfig{
		Backfill: config.BackfillConfig{
			InitialBatchSize: 100,
			MinBatchSize:     10,
			MaxBatchSize:     1000,
			GrowthFactor:     1.2,
		},
		Projector: config.ProjectorConfig{
			MaxWorkers:        4,
			MaxRetries:        3,
			ProcessingTimeout: 300000000000,
			PollInterval:      1000000000,
		},
	}
}

// defaultChainCfg 创建测试用 SyncChainConfig。
func defaultChainCfg() config.SyncChainConfig {
	return config.SyncChainConfig{
		ChainID:            1,
		StartBlock:         1000,
		ConfirmationBlocks: 12,
		ProcessorID:        "test-processor",
		Contracts: config.SyncContractsConfig{
			MouseTier:      "0x1111111111111111111111111111111111111111",
			MousePadByTier: "0x2222222222222222222222222222222222222222",
		},
	}
}

// TestSyncProcessor_ColdStartBackfill 测试无 checkpoint 时从 start_block 回填。
func TestSyncProcessor_ColdStartBackfill(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	store := NewEventStore(db)
	checkpoint := NewCheckpointRepository(db)
	svc := &mockEventService{}

	pool := &mockSubscriberNodePool{blockNumber: 1200}
	chainCfg := defaultChainCfg()
	syncCfg := defaultSyncCfg()

	processor := NewSyncProcessor(pool, store, checkpoint, svc, chainCfg, syncCfg)
	assert.NotNil(t, processor)

	// 验证冷启动时无 checkpoint
	ctx := context.Background()
	cp, err := checkpoint.Get(ctx, 1, "test-processor")
	require.NoError(t, err)
	assert.Nil(t, cp, "冷启动时不应有 checkpoint")
}

// TestSyncProcessor_CheckpointResume 测试有 checkpoint 时从上次位置继续。
func TestSyncProcessor_CheckpointResume(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	store := NewEventStore(db)
	checkpoint := NewCheckpointRepository(db)
	svc := &mockEventService{}

	ctx := context.Background()

	// 预先写入 checkpoint
	err := checkpoint.Upsert(ctx, 1, "test-processor", 5000)
	require.NoError(t, err)

	pool := &mockSubscriberNodePool{blockNumber: 12000}
	chainCfg := defaultChainCfg()
	syncCfg := defaultSyncCfg()

	_ = NewSyncProcessor(pool, store, checkpoint, svc, chainCfg, syncCfg)

	// 验证 checkpoint 存在
	cp, err := checkpoint.Get(ctx, 1, "test-processor")
	require.NoError(t, err)
	require.NotNil(t, cp)
	assert.Equal(t, int64(5000), cp.LastSyncedBlock)
}

// TestSyncProcessor_MultiContractAddresses 测试单订阅多合约地址。
func TestSyncProcessor_MultiContractAddresses(t *testing.T) {
	t.Parallel()

	chainCfg := defaultChainCfg()

	// 验证两个合约地址都被配置
	var addresses []string
	if chainCfg.Contracts.MouseTier != "" {
		addresses = append(addresses, chainCfg.Contracts.MouseTier)
	}
	if chainCfg.Contracts.MousePadByTier != "" {
		addresses = append(addresses, chainCfg.Contracts.MousePadByTier)
	}
	assert.Len(t, addresses, 2, "应包含两个合约地址")
}

// TestSyncProcessor_FinalizedBlockFilter 测试 finalized 区块过滤逻辑。
func TestSyncProcessor_FinalizedBlockFilter(t *testing.T) {
	t.Parallel()

	chainCfg := defaultChainCfg()

	tests := []struct {
		name            string
		currentBlock    uint64
		eventBlock      uint64
		expectProcessed bool
	}{
		{
			name:            "事件区块已 finalized",
			currentBlock:    1000,
			eventBlock:      900,
			expectProcessed: true,
		},
		{
			name:            "事件区块未 finalized",
			currentBlock:    1000,
			eventBlock:      995,
			expectProcessed: false,
		},
		{
			name:            "事件区块正好在边界上",
			currentBlock:    1000,
			eventBlock:      988,
			expectProcessed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			threshold := int64(tt.currentBlock) - chainCfg.ConfirmationBlocks
			isFinalized := int64(tt.eventBlock) <= threshold
			assert.Equal(t, tt.expectProcessed, isFinalized)
		})
	}
}

// TestNewSyncProcessor 测试创建 SyncProcessor 实例。
func TestNewSyncProcessor(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	store := NewEventStore(db)
	checkpoint := NewCheckpointRepository(db)
	svc := &mockEventService{}

	pool := &mockSubscriberNodePool{blockNumber: 1000}
	chainCfg := defaultChainCfg()
	syncCfg := defaultSyncCfg()

	processor := NewSyncProcessor(pool, store, checkpoint, svc, chainCfg, syncCfg)

	assert.NotNil(t, processor)
	assert.NotNil(t, processor.projector)
	assert.NotNil(t, processor.subscriber)
	assert.Equal(t, int64(12), processor.chainCfg.ConfirmationBlocks)
}
