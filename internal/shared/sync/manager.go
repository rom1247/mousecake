package sync

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"gorm.io/gorm"

	"github.com/mousecake-go/mousecake-go/config"
	"github.com/mousecake-go/mousecake-go/internal/chain"
)

// SyncManager 编排多条链的同步处理器。
type SyncManager struct {
	cfg        config.SyncConfig
	db         *gorm.DB
	store      *EventStore
	checkpoint *CheckpointRepository
	svc        EventService
	metrics    *syncMetrics

	processors []*SyncProcessor
	pools      []chain.NodePool
	wg         sync.WaitGroup
	errCh      chan error
	mu         sync.Mutex
}

// NewSyncManager 创建 SyncManager 实例。
func NewSyncManager(cfg config.SyncConfig, db *gorm.DB, svc EventService) (*SyncManager, error) {
	store := NewEventStore(db)
	checkpoint := NewCheckpointRepository(db)

	return &SyncManager{
		cfg:        cfg,
		db:         db,
		store:      store,
		checkpoint: checkpoint,
		svc:        svc,
	}, nil
}

// Start 启动所有链的同步处理器。
func (m *SyncManager) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.errCh = make(chan error, len(m.cfg.Chains))

	for _, chainCfg := range m.cfg.Chains {
		pool, err := chain.NewNodePool(chainCfg.Nodes, chainCfg.ChainID)
		if err != nil {
			return fmt.Errorf("创建 NodePool chain_id=%d: %w", chainCfg.ChainID, err)
		}

		processor := NewSyncProcessor(
			pool,
			m.store,
			m.checkpoint,
			m.svc,
			chainCfg,
			m.cfg,
		)

		m.processors = append(m.processors, processor)
		m.pools = append(m.pools, pool)

		// 每条链在独立 goroutine 中启动，错误通过 errCh 传播
		m.wg.Add(1)
		go func(chainID int, p *SyncProcessor) {
			defer m.wg.Done()
			if err := p.Start(ctx); err != nil {
				slog.Error("SyncProcessor 启动失败",
					"chain_id", chainID,
					"error", err)
				m.errCh <- fmt.Errorf("chain_id=%d: %w", chainID, err)
			}
		}(chainCfg.ChainID, processor)
	}

	// 非阻塞检查是否存在立即启动失败的情况
	select {
	case err := <-m.errCh:
		m.stopProcessors()
		return fmt.Errorf("SyncProcessor 启动失败: %w", err)
	default:
	}

	slog.Info("SyncManager 启动完成", "chains", len(m.processors))
	return nil
}

// Stop 停止所有同步处理器，等待 goroutine 退出。
func (m *SyncManager) Stop() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.stopProcessors()

	// 等待所有 SyncProcessor goroutine 退出
	m.wg.Wait()

	slog.Info("SyncManager 已停止", "chains", len(m.processors))
}

// stopProcessors 停止所有 processor 和 pool。
func (m *SyncManager) stopProcessors() {
	for _, p := range m.processors {
		p.Stop()
	}
}

// GetStatus 获取所有链的同步状态。
func (m *SyncManager) GetStatus(ctx context.Context) ([]ChainSyncStatus, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var statuses []ChainSyncStatus
	for _, chainCfg := range m.cfg.Chains {
		cp, err := m.checkpoint.Get(ctx, chainCfg.ChainID, chainCfg.ProcessorID)
		if err != nil {
			return nil, fmt.Errorf("查询检查点 chain_id=%d: %w", chainCfg.ChainID, err)
		}

		lastSyncedBlock := int64(0)
		if cp != nil {
			lastSyncedBlock = cp.LastSyncedBlock
		}

		counts, err := m.store.CountByStatus(ctx, chainCfg.ChainID)
		if err != nil {
			return nil, fmt.Errorf("统计事件 chain_id=%d: %w", chainCfg.ChainID, err)
		}

		statuses = append(statuses, ChainSyncStatus{
			ChainID:         chainCfg.ChainID,
			LastSyncedBlock: lastSyncedBlock,
			EventCounts:     counts,
		})
	}
	return statuses, nil
}

// ChainSyncStatus 单条链的同步状态。
type ChainSyncStatus struct {
	ChainID         int                   `json:"chain_id"`
	LastSyncedBlock int64                 `json:"last_synced_block"`
	EventCounts     map[EventStatus]int64 `json:"event_counts"`
}
