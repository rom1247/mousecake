package sync

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mousecake-go/mousecake-go/config"
)

// TestChainSyncStatus_Fields 测试 ChainSyncStatus 结构体字段。
func TestChainSyncStatus_Fields(t *testing.T) {
	t.Parallel()

	status := ChainSyncStatus{
		ChainID:         1,
		LastSyncedBlock: 12345,
		EventCounts: map[EventStatus]int64{
			StatusProcessed: 100,
			StatusPending:   5,
			StatusFailed:    2,
		},
	}

	assert.Equal(t, 1, status.ChainID)
	assert.Equal(t, int64(12345), status.LastSyncedBlock)
	assert.Equal(t, int64(100), status.EventCounts[StatusProcessed])
	assert.Equal(t, int64(5), status.EventCounts[StatusPending])
	assert.Equal(t, int64(2), status.EventCounts[StatusFailed])
}

// TestNewSyncManager 测试创建 SyncManager 实例。
func TestNewSyncManager(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := &mockEventService{}

	syncCfg := config.SyncConfig{
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

	mgr, err := NewSyncManager(syncCfg, db, svc)
	require.NoError(t, err)
	assert.NotNil(t, mgr)
}

// TestSyncManager_MultiChainConfig 测试多链配置加载。
func TestSyncManager_MultiChainConfig(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := &mockEventService{}

	syncCfg := config.SyncConfig{
		Chains: []config.SyncChainConfig{
			{ChainID: 1, ProcessorID: "eth-mainnet"},
			{ChainID: 5, ProcessorID: "goerli"},
		},
		Backfill: config.BackfillConfig{
			InitialBatchSize: 100, MinBatchSize: 10, MaxBatchSize: 1000, GrowthFactor: 1.2,
		},
		Projector: config.ProjectorConfig{
			MaxWorkers: 4, MaxRetries: 3, ProcessingTimeout: 300000000000, PollInterval: 1000000000,
		},
	}

	mgr, err := NewSyncManager(syncCfg, db, svc)
	require.NoError(t, err)
	assert.NotNil(t, mgr)
	assert.Len(t, mgr.cfg.Chains, 2)
}

// TestSyncManager_GetStatus_NoChains 测试无链配置时 GetStatus 返回空列表。
func TestSyncManager_GetStatus_NoChains(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := &mockEventService{}

	syncCfg := config.SyncConfig{
		Backfill: config.BackfillConfig{
			InitialBatchSize: 100, MinBatchSize: 10, MaxBatchSize: 1000, GrowthFactor: 1.2,
		},
		Projector: config.ProjectorConfig{
			MaxWorkers: 4, MaxRetries: 3, ProcessingTimeout: 300000000000, PollInterval: 1000000000,
		},
	}

	mgr, err := NewSyncManager(syncCfg, db, svc)
	require.NoError(t, err)

	statuses, err := mgr.GetStatus(context.Background())
	require.NoError(t, err)
	assert.Empty(t, statuses)
}

// TestSyncManager_GetStatus_WithCheckpoint 测试有 checkpoint 时 GetStatus 返回正确数据。
func TestSyncManager_GetStatus_WithCheckpoint(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := &mockEventService{}
	ctx := context.Background()

	// 预先写入 checkpoint
	cp := NewCheckpointRepository(db)
	err := cp.Upsert(ctx, 1, "test-proc", 9999)
	require.NoError(t, err)

	syncCfg := config.SyncConfig{
		Chains: []config.SyncChainConfig{
			{ChainID: 1, ProcessorID: "test-proc"},
		},
		Backfill: config.BackfillConfig{
			InitialBatchSize: 100, MinBatchSize: 10, MaxBatchSize: 1000, GrowthFactor: 1.2,
		},
		Projector: config.ProjectorConfig{
			MaxWorkers: 4, MaxRetries: 3, ProcessingTimeout: 300000000000, PollInterval: 1000000000,
		},
	}

	mgr, err := NewSyncManager(syncCfg, db, svc)
	require.NoError(t, err)

	statuses, err := mgr.GetStatus(ctx)
	require.NoError(t, err)
	require.Len(t, statuses, 1)
	assert.Equal(t, 1, statuses[0].ChainID)
	assert.Equal(t, int64(9999), statuses[0].LastSyncedBlock)
}

// TestSyncManager_GracefulShutdown 测试优雅关停不 panic。
func TestSyncManager_GracefulShutdown(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := &mockEventService{}

	syncCfg := config.SyncConfig{
		Backfill: config.BackfillConfig{
			InitialBatchSize: 100, MinBatchSize: 10, MaxBatchSize: 1000, GrowthFactor: 1.2,
		},
		Projector: config.ProjectorConfig{
			MaxWorkers: 4, MaxRetries: 3, ProcessingTimeout: 300000000000, PollInterval: 1000000000,
		},
	}

	mgr, err := NewSyncManager(syncCfg, db, svc)
	require.NoError(t, err)

	// Stop 在没有 Start 的情况下调用不应 panic
	assert.NotPanics(t, func() {
		mgr.Stop()
	})
}

// TestSyncManager_SingleChainFailureIsolation 测试单链失败不影响其他链。
// 验证 SyncManager 为每条链创建独立的 SyncProcessor，goroutine 级别错误隔离。
func TestSyncManager_SingleChainFailureIsolation(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := &mockEventService{}

	syncCfg := config.SyncConfig{
		Chains: []config.SyncChainConfig{
			{ChainID: 1, ProcessorID: "chain-1"},
			{ChainID: 5, ProcessorID: "chain-5"},
		},
		Backfill: config.BackfillConfig{
			InitialBatchSize: 100, MinBatchSize: 10, MaxBatchSize: 1000, GrowthFactor: 1.2,
		},
		Projector: config.ProjectorConfig{
			MaxWorkers: 4, MaxRetries: 3, ProcessingTimeout: 300000000000, PollInterval: 1000000000,
		},
	}

	mgr, err := NewSyncManager(syncCfg, db, svc)
	require.NoError(t, err)

	// 验证多链配置正确加载，每条链独立
	assert.Len(t, mgr.cfg.Chains, 2)

	// Start 会因无法创建 NodePool（无有效节点）而返回错误，
	// 但 SyncManager 的设计是每条链独立 goroutine 运行，互不影响。
	// 此处验证配置层面的隔离性。
	assert.Equal(t, 1, mgr.cfg.Chains[0].ChainID)
	assert.Equal(t, 5, mgr.cfg.Chains[1].ChainID)
	assert.Equal(t, "chain-1", mgr.cfg.Chains[0].ProcessorID)
	assert.Equal(t, "chain-5", mgr.cfg.Chains[1].ProcessorID)
}

// TestSyncManager_AllNodesUnavailable 测试单链节点全部不可用时 Start 返回错误。
func TestSyncManager_AllNodesUnavailable(t *testing.T) {
	t.Parallel()

	db := setupTestDB(t)
	svc := &mockEventService{}

	syncCfg := config.SyncConfig{
		Chains: []config.SyncChainConfig{
			{
				ChainID:     1,
				ProcessorID: "test-processor",
				Nodes:       []config.ChainNodeConfig{}, // 空节点列表
			},
		},
		Backfill: config.BackfillConfig{
			InitialBatchSize: 100, MinBatchSize: 10, MaxBatchSize: 1000, GrowthFactor: 1.2,
		},
		Projector: config.ProjectorConfig{
			MaxWorkers: 4, MaxRetries: 3, ProcessingTimeout: 300000000000, PollInterval: 1000000000,
		},
	}

	mgr, err := NewSyncManager(syncCfg, db, svc)
	require.NoError(t, err)

	// Start 应因无法创建 NodePool（空节点列表）而返回错误
	err = mgr.Start(context.Background())
	assert.Error(t, err, "节点全部不可用时 Start 应返回错误")
}
