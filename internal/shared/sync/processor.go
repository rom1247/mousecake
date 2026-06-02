package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/mousecake-go/mousecake-go/config"
	"github.com/mousecake-go/mousecake-go/internal/chain"
)

// SyncProcessor 编排单链的完整同步生命周期：历史回填 + 实时订阅 + 事件投影。
type SyncProcessor struct {
	pool       chain.NodePool
	store      *EventStore
	checkpoint *CheckpointRepository
	projector  *Projector
	subscriber *chain.Subscriber
	chainCfg   config.SyncChainConfig
	syncCfg    config.SyncConfig
	metrics    *syncMetrics
}

// NewSyncProcessor 创建 SyncProcessor 实例。
func NewSyncProcessor(
	pool chain.NodePool,
	store *EventStore,
	checkpoint *CheckpointRepository,
	svc EventService,
	chainCfg config.SyncChainConfig,
	syncCfg config.SyncConfig,
) *SyncProcessor {
	metrics := newSyncMetrics(chainCfg.ChainID)

	var addresses []string
	if chainCfg.Contracts.MouseTier != "" {
		addresses = append(addresses, chainCfg.Contracts.MouseTier)
	}
	if chainCfg.Contracts.MousePadByTier != "" {
		addresses = append(addresses, chainCfg.Contracts.MousePadByTier)
	}

	projector := NewProjector(store, svc, chainCfg.ChainID, syncCfg.Projector, metrics)

	subscriber := chain.NewSubscriber(
		pool,
		chainCfg.ChainID,
		chainCfg.Contracts,
		chainCfg.ConfirmationBlocks,
		chainCfg.BlockInterval,
	)

	return &SyncProcessor{
		pool:       pool,
		store:      store,
		checkpoint: checkpoint,
		projector:  projector,
		subscriber: subscriber,
		chainCfg:   chainCfg,
		syncCfg:    syncCfg,
		metrics:    metrics,
	}
}

// Start 启动同步处理器。
func (p *SyncProcessor) Start(ctx context.Context) error {
	slog.Info("启动 SyncProcessor", "chain_id", p.chainCfg.ChainID)

	// 启动 Projector
	p.projector.Start(ctx)

	// 检查 checkpoint 决定起始位置
	cp, err := p.checkpoint.Get(ctx, p.chainCfg.ChainID, p.chainCfg.ProcessorID)
	if err != nil {
		return fmt.Errorf("查询检查点: %w", err)
	}

	startBlock := p.chainCfg.StartBlock
	if cp != nil {
		startBlock = cp.LastSyncedBlock + 1
		slog.Info("从 checkpoint 继续同步",
			"chain_id", p.chainCfg.ChainID,
			"last_synced", cp.LastSyncedBlock,
			"start_from", startBlock)
	}

	// 历史回填
	var addresses []string
	if p.chainCfg.Contracts.MouseTier != "" {
		addresses = append(addresses, p.chainCfg.Contracts.MouseTier)
	}
	if p.chainCfg.Contracts.MousePadByTier != "" {
		addresses = append(addresses, p.chainCfg.Contracts.MousePadByTier)
	}

	backfiller := NewBackfiller(
		p.pool, p.store, p.checkpoint,
		p.chainCfg.ChainID, p.chainCfg.ProcessorID,
		addresses, p.chainCfg.ConfirmationBlocks,
		p.syncCfg.Backfill,
	)

	if err := backfiller.Run(ctx, startBlock); err != nil {
		return fmt.Errorf("历史回填: %w", err)
	}

	// 设置 Subscriber 起始区块
	if cp != nil {
		p.subscriber.SetLastSeenBlock(cp.LastSyncedBlock)
	} else {
		currentBlock, err := p.pool.BlockNumber(ctx)
		if err == nil {
			p.subscriber.SetLastSeenBlock(int64(currentBlock) - p.chainCfg.ConfirmationBlocks)
		}
	}

	// 启动实时订阅
	logCh, err := p.subscriber.Start(ctx)
	if err != nil {
		return fmt.Errorf("启动订阅: %w", err)
	}

	// 消费实时事件
	go p.consumeRealtimeEvents(ctx, logCh)

	slog.Info("SyncProcessor 启动完成", "chain_id", p.chainCfg.ChainID)
	return nil
}

// Stop 停止同步处理器。
func (p *SyncProcessor) Stop() {
	p.subscriber.Stop()
	p.projector.Stop()
	p.pool.Close()
}

// consumeRealtimeEvents 消费实时订阅事件。
func (p *SyncProcessor) consumeRealtimeEvents(ctx context.Context, logCh <-chan types.Log) {
	for {
		select {
		case <-ctx.Done():
			return
		case log, ok := <-logCh:
			if !ok {
				return
			}

			// finalized 区块过滤
			currentBlock, err := p.pool.BlockNumber(ctx)
			if err != nil {
				slog.Warn("获取当前区块号失败", "chain_id", p.chainCfg.ChainID, "error", err)
				continue
			}

			if int64(log.BlockNumber) > int64(currentBlock)-p.chainCfg.ConfirmationBlocks {
				continue
			}

			// 转换并写入
			eventData, _ := json.Marshal(map[string]any{
				"topics":  topicsToHex(log.Topics),
				"data":    fmt.Sprintf("0x%x", log.Data),
				"removed": log.Removed,
			})

			event := ChainEvent{
				ChainID:         p.chainCfg.ChainID,
				BlockNumber:     int64(log.BlockNumber),
				TxHash:          log.TxHash.Hex(),
				TxIndex:         int(log.TxIndex),
				LogIndex:        int(log.Index),
				ContractAddress: log.Address.Hex(),
				EventName:       extractEventNameFromLog(log),
				EventData:       string(eventData),
				Status:          StatusPending,
				ProcessorID:     p.chainCfg.ProcessorID,
			}

			if _, err := p.store.BatchInsert(ctx, []ChainEvent{event}); err != nil {
				slog.Warn("写入实时事件失败",
					"chain_id", p.chainCfg.ChainID,
					"block", log.BlockNumber,
					"error", err)
			}

			// 更新检查点
			if err := p.checkpoint.Upsert(ctx, p.chainCfg.ChainID, p.chainCfg.ProcessorID, int64(log.BlockNumber)); err != nil {
				slog.Warn("更新检查点失败", "chain_id", p.chainCfg.ChainID, "error", err)
			}
		}
	}
}

// extractEventNameFromLog 从日志中提取事件名称。
func extractEventNameFromLog(log types.Log) string {
	if len(log.Topics) == 0 {
		return "Unknown"
	}
	topic := log.Topics[0].Hex()
	if len(topic) >= 10 {
		return topic[:10]
	}
	return "Unknown"
}

// topicsToHex 将 topics 转换为十六进制字符串列表。
func topicsToHex(topics []common.Hash) []string {
	result := make([]string, len(topics))
	for i, t := range topics {
		result[i] = t.Hex()
	}
	return result
}
