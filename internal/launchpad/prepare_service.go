package launchpad

import (
	"context"
	"encoding/hex"
	"fmt"
	"log/slog"
	"time"

	"github.com/mousecake-go/mousecake-go/internal/launchpad/domain"
)

// CreatePrepareInput 创建 Prepare 交易的输入参数。
type CreatePrepareInput struct {
	// OperationType 操作类型。
	OperationType string `json:"operation_type" binding:"required"`
	// CallerAddress 调用者地址。
	CallerAddress string `json:"caller_address" binding:"required"`
	// SaleID 关联的 Sale ID。
	SaleID *int64 `json:"sale_id"`
	// PoolIndex 关联的池序号。
	PoolIndex *int64 `json:"pool_index"`
	// Calldata 可选的预编码 calldata（由上层管理员/用户用例填充）。
	Calldata []byte `json:"-"`
}

// PrepareService 编排 Prepare 交易的创建、提交确认、超时轮询和取消。
type PrepareService struct {
	repo      domain.PrepareTxRepository
	chain     domain.ChainReader
	encoder   EncoderInterface
	expiresIn time.Duration
	log       *slog.Logger
}

// NewPrepareService 创建 PrepareService。
func NewPrepareService(
	repo domain.PrepareTxRepository,
	chain domain.ChainReader,
	encoder EncoderInterface,
	expiresIn time.Duration,
) *PrepareService {
	return &PrepareService{
		repo:      repo,
		chain:     chain,
		encoder:   encoder,
		expiresIn: expiresIn,
		log:       slog.Default().With("module", "launchpad", "layer", "service"),
	}
}

// Create 创建 Prepare 交易。
// 校验操作类型，通过 calldata_hash 去重，写入 pending 记录。
func (s *PrepareService) Create(ctx context.Context, input CreatePrepareInput) (*domain.PrepareTx, error) {
	opType := domain.PrepareTxOperationType(input.OperationType)
	if !domain.IsValidOperationType(opType) {
		return nil, fmt.Errorf("无效操作类型: %s: %w", input.OperationType, domain.ErrInvalidTransition)
	}

	var calldataHex string
	var calldataHash string
	if input.Calldata != nil {
		calldataHex = "0x" + hex.EncodeToString(input.Calldata)
		calldataHash = s.encoder.CalldataHash(input.Calldata)
	}

	// 去重：如果已有相同 calldata_hash 的 pending 记录，返回已有记录
	if calldataHash != "" {
		existing, err := s.repo.FindPendingByCalldataHash(ctx, calldataHash)
		if err != nil {
			return nil, fmt.Errorf("检查重复 prepare: %w", err)
		}
		if existing != nil {
			return existing, nil
		}
	}

	now := time.Now()
	tx := &domain.PrepareTx{
		SaleID:        input.SaleID,
		PoolIndex:     input.PoolIndex,
		OperationType: opType,
		CallerAddress: input.CallerAddress,
		Calldata:      calldataHex,
		CalldataHash:  calldataHash,
		Status:        domain.PrepareTxPending,
		ExpiresAt:     now.Add(s.expiresIn),
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.repo.Create(ctx, tx); err != nil {
		return nil, fmt.Errorf("创建 prepare_tx: %w", err)
	}

	return tx, nil
}

// Submit 提交 txHash，查询 receipt 并更新状态。
// 仅 broadcast 状态的 prepare 可提交。
func (s *PrepareService) Submit(ctx context.Context, id int64, txHash string) (*domain.PrepareTx, error) {
	tx, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("查询 prepare_tx %d: %w", id, err)
	}

	if tx.Status != domain.PrepareTxBroadcast {
		return nil, fmt.Errorf("prepare_tx %d 状态为 %s，只有 broadcast 可提交: %w", id, tx.Status, domain.ErrInvalidTransition)
	}

	if tx.IsExpired() {
		return nil, fmt.Errorf("prepare_tx %d 已过期: %w", id, domain.ErrInvalidTransition)
	}

	// 查询链上 receipt
	receipt, err := s.chain.GetTransactionReceipt(ctx, txHash)
	if err != nil {
		return nil, fmt.Errorf("查询 receipt %s: %w", txHash, err)
	}

	updates := map[string]any{
		"tx_hash": txHash,
	}

	// receipt 尚不可用，保持 broadcast 状态
	if receipt == nil {
		if err := s.repo.UpdateStatus(ctx, id, domain.PrepareTxBroadcast, updates); err != nil {
			return nil, fmt.Errorf("更新 prepare_tx %d broadcast: %w", id, err)
		}
		tx.TxHash = &txHash
		return tx, nil
	}

	// receipt 可用，根据 status 判断成功或失败
	updates["block_number"] = int64(receipt.BlockNumber)

	if receipt.Status == 1 {
		now := time.Now()
		updates["confirmed_at"] = now
		if err := s.repo.UpdateStatus(ctx, id, domain.PrepareTxConfirmed, updates); err != nil {
			return nil, fmt.Errorf("更新 prepare_tx %d confirmed: %w", id, err)
		}
		tx.Status = domain.PrepareTxConfirmed
		tx.TxHash = &txHash
		tx.BlockNumber = int64Ptr(int64(receipt.BlockNumber))
		tx.ConfirmedAt = &now
	} else {
		errMsg := "链上交易执行失败"
		updates["error_message"] = errMsg
		if err := s.repo.UpdateStatus(ctx, id, domain.PrepareTxReverted, updates); err != nil {
			return nil, fmt.Errorf("更新 prepare_tx %d reverted: %w", id, err)
		}
		tx.Status = domain.PrepareTxReverted
		tx.TxHash = &txHash
		tx.BlockNumber = int64Ptr(int64(receipt.BlockNumber))
		tx.ErrorMessage = &errMsg
	}

	return tx, nil
}

// PollTimeout 扫描超时的 pending 和 broadcast 记录，更新状态。
func (s *PrepareService) PollTimeout(ctx context.Context) error {
	// 处理已过期的 pending 记录
	expired, err := s.repo.FindPendingExpired(ctx)
	if err != nil {
		return fmt.Errorf("查询过期 pending: %w", err)
	}
	for _, tx := range expired {
		if err := s.repo.UpdateStatus(ctx, tx.ID, domain.PrepareTxExpired, nil); err != nil {
			return fmt.Errorf("标记 prepare_tx %d 过期: %w", tx.ID, err)
		}
	}

	// 处理超时的 broadcast 记录
	timeoutTxs, err := s.repo.FindBroadcastTimeout(ctx, 10*time.Minute)
	if err != nil {
		return fmt.Errorf("查询超时 broadcast: %w", err)
	}
	for _, tx := range timeoutTxs {
		if tx.TxHash == nil {
			continue
		}

		receipt, err := s.chain.GetTransactionReceipt(ctx, *tx.TxHash)
		if err != nil {
			continue
		}
		if receipt == nil {
			continue
		}

		updates := map[string]any{
			"block_number": int64(receipt.BlockNumber),
		}
		if receipt.Status == 1 {
			now := time.Now()
			updates["confirmed_at"] = now
			if err := s.repo.UpdateStatus(ctx, tx.ID, domain.PrepareTxConfirmed, updates); err != nil {
				s.log.ErrorContext(ctx, "更新 prepare_tx 状态为 confirmed 失败", "id", tx.ID, "error", err)
			}
		} else {
			updates["error_message"] = "链上交易执行失败"
			if err := s.repo.UpdateStatus(ctx, tx.ID, domain.PrepareTxReverted, updates); err != nil {
				s.log.ErrorContext(ctx, "更新 prepare_tx 状态为 reverted 失败", "id", tx.ID, "error", err)
			}
		}
	}

	return nil
}

// GetByID 根据 ID 查询 Prepare 交易。
func (s *PrepareService) GetByID(ctx context.Context, id int64) (*domain.PrepareTx, error) {
	tx, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("查询 prepare_tx %d: %w", id, err)
	}
	return tx, nil
}

// Cancel 取消指定 pending 的 Prepare 交易。
func (s *PrepareService) Cancel(ctx context.Context, id int64) error {
	tx, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("查询 prepare_tx %d: %w", id, err)
	}

	if tx.Status != domain.PrepareTxPending {
		return fmt.Errorf("prepare_tx %d 状态为 %s，只有 pending 可取消: %w", id, tx.Status, domain.ErrInvalidTransition)
	}

	if err := s.repo.UpdateStatus(ctx, id, domain.PrepareTxExpired, nil); err != nil {
		return fmt.Errorf("取消 prepare_tx %d: %w", id, err)
	}

	return nil
}
