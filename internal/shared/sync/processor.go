package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/mousecake-go/mousecake-go/config"
	"github.com/mousecake-go/mousecake-go/internal/chain"
)

// SyncProcessor 编排单链的完整同步生命周期：历史回填 + 实时订阅 + 补偿循环 + 事件投影。
type SyncProcessor struct {
	pool       chain.NodePool
	store      *EventStore
	checkpoint *CheckpointRepository
	projector  *Projector
	subscriber *chain.Subscriber
	chainCfg   config.SyncChainConfig
	syncCfg    config.SyncConfig
	metrics    *syncMetrics
	addresses  []common.Address
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

	addresses := parseContractAddresses(chainCfg.Contracts)

	projector := NewProjector(store, svc, chainCfg.ChainID, syncCfg.Projector, metrics)

	subscriber := chain.NewSubscriber(pool, chainCfg.ChainID, addresses)

	return &SyncProcessor{
		pool:       pool,
		store:      store,
		checkpoint: checkpoint,
		projector:  projector,
		subscriber: subscriber,
		chainCfg:   chainCfg,
		syncCfg:    syncCfg,
		metrics:    metrics,
		addresses:  addresses,
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
	addresses := addressesToStrings(p.addresses)
	backfiller := NewBackfiller(
		p.pool, p.store, p.checkpoint,
		p.chainCfg.ChainID, p.chainCfg.ProcessorID,
		addresses, p.chainCfg.ConfirmationBlocks,
		p.syncCfg.Backfill,
	)

	if err := backfiller.Run(ctx, startBlock); err != nil {
		return fmt.Errorf("历史回填: %w", err)
	}

	// 启动实时订阅
	logCh, err := p.subscriber.Start(ctx)
	if err != nil {
		return fmt.Errorf("启动订阅: %w", err)
	}

	// 消费实时事件
	go p.consumeRealtimeEvents(ctx, logCh)

	// 启动补偿循环
	go p.compensationLoop(ctx)

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
	var lastFlushedBlock int64

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
			event := convertLog(log, p.chainCfg.ChainID, p.chainCfg.ProcessorID)
			if _, err := p.store.BatchInsert(ctx, []ChainEvent{event}); err != nil {
				slog.Warn("写入实时事件失败",
					"chain_id", p.chainCfg.ChainID,
					"block", log.BlockNumber,
					"error", err)
				continue
			}

			// 按区块更新检查点，避免每条事件都写 DB
			blockNum := int64(log.BlockNumber)
			if blockNum > lastFlushedBlock {
				lastFlushedBlock = blockNum
				if err := p.checkpoint.Upsert(ctx, p.chainCfg.ChainID, p.chainCfg.ProcessorID, blockNum); err != nil {
					slog.Warn("更新检查点失败", "chain_id", p.chainCfg.ChainID, "error", err)
				}
			}
		}
	}
}

// compensationLoop 定期补偿 WS 断连期间可能遗漏的事件。
func (p *SyncProcessor) compensationLoop(ctx context.Context) {
	if p.syncCfg.CompensateInterval <= 0 {
		return
	}

	ticker := time.NewTicker(p.syncCfg.CompensateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.compensateOnce(ctx)
		}
	}
}

// compensateOnce 执行一次补偿查询。
func (p *SyncProcessor) compensateOnce(ctx context.Context) {
	cp, err := p.checkpoint.Get(ctx, p.chainCfg.ChainID, p.chainCfg.ProcessorID)
	if err != nil {
		slog.Warn("补偿查询 checkpoint 失败",
			"chain_id", p.chainCfg.ChainID, "error", err)
		return
	}

	fromBlock := p.chainCfg.StartBlock
	if cp != nil {
		fromBlock = cp.LastSyncedBlock + 1
	}

	currentBlock, err := p.pool.BlockNumber(ctx)
	if err != nil {
		slog.Warn("补偿查询获取区块号失败",
			"chain_id", p.chainCfg.ChainID, "error", err)
		return
	}

	toBlock := int64(currentBlock) - p.chainCfg.ConfirmationBlocks
	if toBlock <= fromBlock {
		return
	}

	// 限制单次补偿范围，避免查询过大
	const maxRange int64 = 5
	if toBlock-fromBlock > maxRange {
		toBlock = fromBlock + maxRange
	}

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(fromBlock),
		ToBlock:   big.NewInt(toBlock),
		Addresses: p.addresses,
	}

	logs, err := p.pool.FilterLogs(ctx, query)
	if err != nil {
		slog.Warn("补偿 FilterLogs 失败",
			"chain_id", p.chainCfg.ChainID,
			"from", fromBlock, "to", toBlock,
			"error", err)
		return
	}

	if len(logs) == 0 {
		return
	}

	events := convertLogs(logs, p.chainCfg.ChainID, p.chainCfg.ProcessorID)
	inserted, err := p.store.BatchInsert(ctx, events)
	if err != nil {
		slog.Warn("补偿写入事件失败",
			"chain_id", p.chainCfg.ChainID, "error", err)
		return
	}

	if inserted > 0 {
		slog.Info("补偿发现遗漏事件",
			"chain_id", p.chainCfg.ChainID,
			"from", fromBlock, "to", toBlock,
			"inserted", inserted)
	}

	// 更新 checkpoint
	if err := p.checkpoint.Upsert(ctx, p.chainCfg.ChainID, p.chainCfg.ProcessorID, toBlock); err != nil {
		slog.Warn("补偿更新检查点失败",
			"chain_id", p.chainCfg.ChainID, "error", err)
	}
}

// --- 包级公共函数 ---

// convertLog 将单条链上日志转换为 ChainEvent。
func convertLog(log types.Log, chainID int, processorID string) ChainEvent {
	eventData, _ := json.Marshal(map[string]any{
		"topics":     topicsToHex(log.Topics),
		"data":       fmt.Sprintf("0x%x", log.Data),
		"block_hash": log.BlockHash.Hex(),
		"removed":    log.Removed,
	})

	return ChainEvent{
		ChainID:         chainID,
		BlockNumber:     int64(log.BlockNumber),
		TxHash:          log.TxHash.Hex(),
		TxIndex:         int(log.TxIndex),
		LogIndex:        int(log.Index),
		ContractAddress: log.Address.Hex(),
		EventName:       extractEventNameFromLog(log),
		EventData:       string(eventData),
		Status:          StatusPending,
		ProcessorID:     processorID,
	}
}

// convertLogs 将多条链上日志转换为 ChainEvent 列表。
func convertLogs(logs []types.Log, chainID int, processorID string) []ChainEvent {
	events := make([]ChainEvent, 0, len(logs))
	for _, log := range logs {
		events = append(events, convertLog(log, chainID, processorID))
	}
	return events
}

// parseContractAddresses 从合约配置解析地址列表。
func parseContractAddresses(contracts config.SyncContractsConfig) []common.Address {
	var addresses []common.Address
	if contracts.MouseTier != "" {
		addresses = append(addresses, common.HexToAddress(contracts.MouseTier))
	}
	if contracts.MousePadByTier != "" {
		addresses = append(addresses, common.HexToAddress(contracts.MousePadByTier))
	}
	return addresses
}

// parseStringAddresses 将字符串地址列表转换为 common.Address 列表。
func parseStringAddresses(addresses []string) []common.Address {
	result := make([]common.Address, 0, len(addresses))
	for _, a := range addresses {
		result = append(result, common.HexToAddress(a))
	}
	return result
}

// addressesToStrings 将 common.Address 列表转换为字符串列表。
func addressesToStrings(addresses []common.Address) []string {
	result := make([]string, 0, len(addresses))
	for _, a := range addresses {
		result = append(result, a.Hex())
	}
	return result
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
