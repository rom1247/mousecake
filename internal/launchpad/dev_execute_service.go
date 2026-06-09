package launchpad

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/mousecake-go/mousecake-go/internal/chain"
	"github.com/mousecake-go/mousecake-go/internal/launchpad/domain"
)

// TxSigner 定义交易签名广播接口。
type TxSigner interface {
	// SignAndBroadcast 签名并广播交易到链上，返回交易哈希。
	SignAndBroadcast(ctx context.Context, to common.Address, data []byte, value *big.Int) (common.Hash, error)
}

// DevExecuteResult 开发执行 PrepareTx 的完整结果，包含 receipt 和事件解析数据。
type DevExecuteResult struct {
	// TxHash 交易哈希。
	TxHash string
	// BlockNumber 交易所在区块号。
	BlockNumber uint64
	// GasUsed 交易消耗的 gas。
	GasUsed uint64
	// Status 交易状态（1=成功，0=revert）。
	Status uint64
	// Events 解析后的合约事件列表。
	Events []chain.ParsedEvent
}

// DevExecuteService 用于开发环境一键执行 pending 状态的 PrepareTx。
type DevExecuteService struct {
	findByID     func(ctx context.Context, id int64) (*domain.PrepareTx, error)
	signer       TxSigner
	updateStatus func(ctx context.Context, id int64, status domain.PrepareTxStatus, updates map[string]any) error
	waitForRcpt  func(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	parseEvents  func(logs []*types.Log) []chain.ParsedEvent
}

// NewDevExecuteService 创建 DevExecuteService 实例。
// signer 需同时提供 SignAndBroadcast 和 WaitForReceipt 能力。
func NewDevExecuteService(
	prepareTxRepo *PrepareTxRepository,
	signer *chain.Signer,
) *DevExecuteService {
	parser := chain.NewEventParser()
	return &DevExecuteService{
		findByID:     prepareTxRepo.FindByID,
		signer:       signer,
		updateStatus: prepareTxRepo.UpdateStatus,
		waitForRcpt:  signer.WaitForReceipt,
		parseEvents:  parser.ParseLogs,
	}
}

// Execute 执行 pending 状态的 PrepareTx：验证 → 签名广播 → 等待 receipt → 解析事件 → 更新状态 → 返回结果。
func (s *DevExecuteService) Execute(ctx context.Context, id int64) (*DevExecuteResult, error) {
	ptx, err := s.findByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("查询 prepare_tx id=%d: %w", id, err)
	}

	if ptx.Status != domain.PrepareTxPending {
		return nil, fmt.Errorf("只能执行 pending 状态的 PrepareTx，当前状态: %s", ptx.Status)
	}

	if ptx.IsExpired() {
		return nil, fmt.Errorf("PrepareTx id=%d 已过期，过期时间: %s", id, ptx.ExpiresAt.Format("2006-01-02 15:04:05"))
	}

	if strings.TrimSpace(ptx.TargetAddress) == "" {
		return nil, fmt.Errorf("PrepareTx id=%d 缺少目标合约地址", id)
	}

	to := common.HexToAddress(ptx.TargetAddress)
	data := common.Hex2Bytes(strings.TrimPrefix(ptx.Calldata, "0x"))
	value := new(big.Int)
	if ptx.Value != "" {
		value.SetString(ptx.Value, 10)
	}

	txHash, err := s.signer.SignAndBroadcast(ctx, to, data, value)
	if err != nil {
		//打印错误日志 要具体的错误信息
		slog.Error("签名广播失败", "error", err)
		return nil, fmt.Errorf("签名广播 PrepareTx id=%d: %w", id, err)
	}

	// 等待 receipt
	receipt, err := s.waitForRcpt(ctx, txHash)
	if err != nil {
		// 等待超时，更新为 broadcast 状态（后续可通过 PollTimeout 继续轮询）
		_ = s.updateStatus(ctx, id, domain.PrepareTxBroadcast, map[string]any{
			"tx_hash": txHash.Hex(),
		})
		return nil, fmt.Errorf("等待 PrepareTx id=%d receipt: %w", id, err)
	}

	// 交易 revert
	if receipt.Status == types.ReceiptStatusFailed {
		errMsg := "链上交易执行失败"
		_ = s.updateStatus(ctx, id, domain.PrepareTxReverted, map[string]any{
			"tx_hash":       txHash.Hex(),
			"block_number":  receipt.BlockNumber.Int64(),
			"error_message": errMsg,
		})
		return nil, fmt.Errorf("PrepareTx id=%d %s", id, errMsg)
	}

	// 解析事件
	events := s.parseEvents(receipt.Logs)

	// 更新状态为 confirmed
	now := time.Now()
	if err := s.updateStatus(ctx, id, domain.PrepareTxConfirmed, map[string]any{
		"tx_hash":      txHash.Hex(),
		"block_number": receipt.BlockNumber.Int64(),
		"confirmed_at": now,
	}); err != nil {
		return nil, fmt.Errorf("更新 PrepareTx 状态 id=%d: %w", id, err)
	}

	return &DevExecuteResult{
		TxHash:      txHash.Hex(),
		BlockNumber: receipt.BlockNumber.Uint64(),
		GasUsed:     receipt.GasUsed,
		Status:      receipt.Status,
		Events:      events,
	}, nil
}
