package sync

import (
	"context"
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

// Backfiller 历史区块回填器，支持自适应 batch_size。
type Backfiller struct {
	pool               chain.NodePool
	store              *EventStore
	checkpoint         *CheckpointRepository
	chainID            int
	processorID        string
	addresses          []common.Address
	confirmationBlocks int64
	cfg                config.BackfillConfig
}

// NewBackfiller 创建回填器实例。
func NewBackfiller(
	pool chain.NodePool,
	store *EventStore,
	checkpoint *CheckpointRepository,
	chainID int,
	processorID string,
	addresses []string,
	confirmationBlocks int64,
	cfg config.BackfillConfig,
) *Backfiller {
	return &Backfiller{
		pool:               pool,
		store:              store,
		checkpoint:         checkpoint,
		chainID:            chainID,
		processorID:        processorID,
		addresses:          parseStringAddresses(addresses),
		confirmationBlocks: confirmationBlocks,
		cfg:                cfg,
	}
}

// Run 执行历史回填，从 startBlock 到 currentBlock - confirmationBlocks。
func (b *Backfiller) Run(ctx context.Context, startBlock int64) error {
	currentBlock, err := b.pool.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("获取当前区块号: %w", err)
	}

	targetBlock := int64(currentBlock) - b.confirmationBlocks
	if targetBlock < startBlock {
		slog.Info("无需回填，起始区块已超过目标",
			"chain_id", b.chainID, "start", startBlock, "target", targetBlock)
		return nil
	}

	slog.Info("开始历史回填",
		"chain_id", b.chainID,
		"from", startBlock,
		"to", targetBlock,
		"batch_size", b.cfg.InitialBatchSize)

	batchSize := b.cfg.InitialBatchSize
	fromBlock := startBlock

	for fromBlock <= targetBlock {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		toBlock := fromBlock + int64(batchSize) - 1
		if toBlock > targetBlock {
			toBlock = targetBlock
		}

		batchStart := time.Now()
		logs, err := b.fetchLogs(ctx, fromBlock, toBlock)
		if err != nil {
			// 超时减半
			newBatch := batchSize / 2
			if newBatch < b.cfg.MinBatchSize {
				newBatch = b.cfg.MinBatchSize
			}
			slog.Warn("FilterLogs 失败，减小 batch_size",
				"chain_id", b.chainID,
				"from", fromBlock, "to", toBlock,
				"old_batch", batchSize, "new_batch", newBatch,
				"error", err)
			batchSize = newBatch
			continue
		}

		if len(logs) > 0 {
			events := convertLogs(logs, b.chainID, b.processorID)
			if _, err := b.store.BatchInsert(ctx, events); err != nil {
				return fmt.Errorf("写入区块 %d-%d 事件: %w", fromBlock, toBlock, err)
			}
		}

		if err := b.checkpoint.Upsert(ctx, b.chainID, b.processorID, toBlock); err != nil {
			return fmt.Errorf("写入检查点 %d: %w", toBlock, err)
		}

		// 成功时增长 batch_size
		newBatch := float64(batchSize) * b.cfg.GrowthFactor
		if newBatch > float64(b.cfg.MaxBatchSize) {
			newBatch = float64(b.cfg.MaxBatchSize)
		}
		batchSize = int(newBatch)

		slog.Debug("回填进度",
			"chain_id", b.chainID,
			"from", fromBlock, "to", toBlock,
			"logs", len(logs),
			"batch_size", batchSize,
			"elapsed", time.Since(batchStart).Round(time.Millisecond))

		fromBlock = toBlock + 1
	}

	slog.Info("历史回填完成",
		"chain_id", b.chainID,
		"last_block", targetBlock)

	return nil
}

// fetchLogs 查询指定区块范围的日志。
func (b *Backfiller) fetchLogs(ctx context.Context, fromBlock, toBlock int64) ([]types.Log, error) {
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(fromBlock),
		ToBlock:   big.NewInt(toBlock),
		Addresses: b.addresses,
	}
	return b.pool.FilterLogs(ctx, query)
}
