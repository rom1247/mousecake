package launchpad

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/mousecake-go/mousecake-go/internal/launchpad/domain"
)

// TxSigner 定义交易签名广播接口。
type TxSigner interface {
	// SignAndBroadcast 签名并广播交易到链上，返回交易哈希。
	SignAndBroadcast(ctx context.Context, to common.Address, data []byte, value *big.Int) (common.Hash, error)
}

// DevExecuteService 用于开发环境一键执行 pending 状态的 PrepareTx。
type DevExecuteService struct {
	findByID     func(ctx context.Context, id int64) (*domain.PrepareTx, error)
	signer       TxSigner
	updateStatus func(ctx context.Context, id int64, status domain.PrepareTxStatus, updates map[string]any) error
}

// NewDevExecuteService 创建 DevExecuteService 实例。
func NewDevExecuteService(
	prepareTxRepo *PrepareTxRepository,
	signer TxSigner,
) *DevExecuteService {
	return &DevExecuteService{
		findByID:     prepareTxRepo.FindByID,
		signer:       signer,
		updateStatus: prepareTxRepo.UpdateStatus,
	}
}

// Execute 执行 pending 状态的 PrepareTx：验证状态 → 验证过期 → 验证目标地址 → 签名广播 → 更新状态。
func (s *DevExecuteService) Execute(ctx context.Context, id int64) (string, error) {
	ptx, err := s.findByID(ctx, id)
	if err != nil {
		return "", fmt.Errorf("查询 prepare_tx id=%d: %w", id, err)
	}

	if ptx.Status != domain.PrepareTxPending {
		return "", fmt.Errorf("只能执行 pending 状态的 PrepareTx，当前状态: %s", ptx.Status)
	}

	if ptx.IsExpired() {
		return "", fmt.Errorf("PrepareTx id=%d 已过期，过期时间: %s", id, ptx.ExpiresAt.Format("2006-01-02 15:04:05"))
	}

	if strings.TrimSpace(ptx.TargetAddress) == "" {
		return "", fmt.Errorf("PrepareTx id=%d 缺少目标合约地址", id)
	}

	to := common.HexToAddress(ptx.TargetAddress)
	data := common.Hex2Bytes(strings.TrimPrefix(ptx.Calldata, "0x"))
	value := new(big.Int)
	if ptx.Value != "" {
		value.SetString(ptx.Value, 10)
	}

	txHash, err := s.signer.SignAndBroadcast(ctx, to, data, value)
	if err != nil {
		return "", fmt.Errorf("签名广播 PrepareTx id=%d: %w", id, err)
	}

	if err := s.updateStatus(ctx, id, domain.PrepareTxBroadcast, map[string]any{
		"tx_hash": txHash.Hex(),
	}); err != nil {
		return "", fmt.Errorf("更新 PrepareTx 状态 id=%d: %w", id, err)
	}

	return txHash.Hex(), nil
}
