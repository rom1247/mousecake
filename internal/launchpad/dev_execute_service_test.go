package launchpad

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mousecake-go/mousecake-go/internal/chain"
	"github.com/mousecake-go/mousecake-go/internal/launchpad/domain"
)

// mockTxSigner 是 TxSigner 接口的 mock 实现。
type mockTxSigner struct {
	hash common.Hash
	err  error
}

func (m *mockTxSigner) SignAndBroadcast(_ context.Context, _ common.Address, _ []byte, _ *big.Int) (common.Hash, error) {
	return m.hash, m.err
}

// buildValidPendingPrepareTx 构建一个合法的 pending 状态 PrepareTx。
func buildValidPendingPrepareTx() *domain.PrepareTx {
	now := time.Now()
	saleID := int64(1)
	return domain.ReconstructPrepareTx(
		1, &saleID, nil,
		domain.OpCreateSale,
		"0xcaller",
		"0x1234567890abcdef1234567890abcdef12345678",
		"0",
		"0xdeadbeef",
		"hash123",
		domain.PrepareTxPending,
		nil, nil, nil,
		now.Add(1*time.Hour), nil,
		now, now,
	)
}

func TestDevExecuteService_Execute_成功路径(t *testing.T) {
	ptx := buildValidPendingPrepareTx()
	expectedHash := common.HexToHash("0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890ab")

	receipt := &types.Receipt{
		Status:      types.ReceiptStatusSuccessful,
		BlockNumber: big.NewInt(12345),
		GasUsed:     150000,
		Logs: []*types.Log{
			{
				Address: common.HexToAddress("0x1234"),
				Topics:  []common.Hash{common.HexToHash("0x01")},
				Data:    []byte{},
				Index:   0,
			},
		},
	}
	parsedEvents := []chain.ParsedEvent{
		{Name: "PoolSet", Fields: map[string]any{"pid": "1"}, Address: common.HexToAddress("0x1234"), LogIndex: 0},
	}

	svc := &DevExecuteService{
		findByID: func(_ context.Context, id int64) (*domain.PrepareTx, error) {
			assert.Equal(t, int64(1), id)
			return ptx, nil
		},
		signer: &mockTxSigner{hash: expectedHash},
		updateStatus: func(_ context.Context, id int64, status domain.PrepareTxStatus, updates map[string]any) error {
			assert.Equal(t, int64(1), id)
			assert.Equal(t, domain.PrepareTxConfirmed, status)
			return nil
		},
		waitForRcpt: func(_ context.Context, txHash common.Hash) (*types.Receipt, error) {
			assert.Equal(t, expectedHash, txHash)
			return receipt, nil
		},
		parseEvents: func(logs []*types.Log) []chain.ParsedEvent {
			return parsedEvents
		},
	}

	result, err := svc.Execute(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, expectedHash.Hex(), result.TxHash)
	assert.Equal(t, uint64(12345), result.BlockNumber)
	assert.Equal(t, uint64(150000), result.GasUsed)
	assert.Equal(t, uint64(1), result.Status)
	require.Len(t, result.Events, 1)
	assert.Equal(t, "PoolSet", result.Events[0].Name)
}

func TestDevExecuteService_Execute_PrepareTx不存在(t *testing.T) {
	svc := &DevExecuteService{
		findByID: func(_ context.Context, _ int64) (*domain.PrepareTx, error) {
			return nil, domain.ErrNotFound
		},
		signer:       &mockTxSigner{},
		updateStatus: func(_ context.Context, _ int64, _ domain.PrepareTxStatus, _ map[string]any) error { return nil },
		waitForRcpt:  func(_ context.Context, _ common.Hash) (*types.Receipt, error) { return nil, nil },
		parseEvents:  func(_ []*types.Log) []chain.ParsedEvent { return nil },
	}

	_, err := svc.Execute(context.Background(), 999)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
}

func TestDevExecuteService_Execute_PrepareTx非pending状态(t *testing.T) {
	ptx := buildValidPendingPrepareTx()
	ptx.Status = domain.PrepareTxBroadcast

	svc := &DevExecuteService{
		findByID:     func(_ context.Context, _ int64) (*domain.PrepareTx, error) { return ptx, nil },
		signer:       &mockTxSigner{},
		updateStatus: func(_ context.Context, _ int64, _ domain.PrepareTxStatus, _ map[string]any) error { return nil },
		waitForRcpt:  func(_ context.Context, _ common.Hash) (*types.Receipt, error) { return nil, nil },
		parseEvents:  func(_ []*types.Log) []chain.ParsedEvent { return nil },
	}

	_, err := svc.Execute(context.Background(), 1)
	require.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "只能执行 pending"))
}

func TestDevExecuteService_Execute_PrepareTx已过期(t *testing.T) {
	ptx := buildValidPendingPrepareTx()
	ptx.ExpiresAt = time.Now().Add(-1 * time.Hour)

	svc := &DevExecuteService{
		findByID:     func(_ context.Context, _ int64) (*domain.PrepareTx, error) { return ptx, nil },
		signer:       &mockTxSigner{},
		updateStatus: func(_ context.Context, _ int64, _ domain.PrepareTxStatus, _ map[string]any) error { return nil },
		waitForRcpt:  func(_ context.Context, _ common.Hash) (*types.Receipt, error) { return nil, nil },
		parseEvents:  func(_ []*types.Log) []chain.ParsedEvent { return nil },
	}

	_, err := svc.Execute(context.Background(), 1)
	require.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "已过期"))
}

func TestDevExecuteService_Execute_缺少target_address(t *testing.T) {
	ptx := buildValidPendingPrepareTx()
	ptx.TargetAddress = ""

	svc := &DevExecuteService{
		findByID:     func(_ context.Context, _ int64) (*domain.PrepareTx, error) { return ptx, nil },
		signer:       &mockTxSigner{},
		updateStatus: func(_ context.Context, _ int64, _ domain.PrepareTxStatus, _ map[string]any) error { return nil },
		waitForRcpt:  func(_ context.Context, _ common.Hash) (*types.Receipt, error) { return nil, nil },
		parseEvents:  func(_ []*types.Log) []chain.ParsedEvent { return nil },
	}

	_, err := svc.Execute(context.Background(), 1)
	require.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "缺少目标合约地址"))
}

func TestDevExecuteService_Execute_签名广播失败(t *testing.T) {
	ptx := buildValidPendingPrepareTx()
	signErr := fmt.Errorf("网络超时: %w", errors.New("connection refused"))

	svc := &DevExecuteService{
		findByID: func(_ context.Context, _ int64) (*domain.PrepareTx, error) {
			return ptx, nil
		},
		signer: &mockTxSigner{err: signErr},
		updateStatus: func(_ context.Context, _ int64, _ domain.PrepareTxStatus, _ map[string]any) error {
			t.Fatal("签名广播失败时不应调用 updateStatus")
			return nil
		},
		waitForRcpt: func(_ context.Context, _ common.Hash) (*types.Receipt, error) { return nil, nil },
		parseEvents: func(_ []*types.Log) []chain.ParsedEvent { return nil },
	}

	_, err := svc.Execute(context.Background(), 1)
	require.Error(t, err)
	assert.True(t, errors.Is(err, signErr))
	assert.Equal(t, domain.PrepareTxPending, ptx.Status)
}

func TestDevExecuteService_Execute_等待receipt超时(t *testing.T) {
	ptx := buildValidPendingPrepareTx()
	expectedHash := common.HexToHash("0xabc")

	svc := &DevExecuteService{
		findByID:     func(_ context.Context, _ int64) (*domain.PrepareTx, error) { return ptx, nil },
		signer:       &mockTxSigner{hash: expectedHash},
		updateStatus: func(_ context.Context, _ int64, _ domain.PrepareTxStatus, _ map[string]any) error { return nil },
		waitForRcpt: func(_ context.Context, _ common.Hash) (*types.Receipt, error) {
			return nil, fmt.Errorf("等待 receipt 超时")
		},
		parseEvents: func(_ []*types.Log) []chain.ParsedEvent { return nil },
	}

	_, err := svc.Execute(context.Background(), 1)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "receipt")
}

func TestDevExecuteService_Execute_交易revert(t *testing.T) {
	ptx := buildValidPendingPrepareTx()
	expectedHash := common.HexToHash("0xabc")

	svc := &DevExecuteService{
		findByID: func(_ context.Context, _ int64) (*domain.PrepareTx, error) { return ptx, nil },
		signer:   &mockTxSigner{hash: expectedHash},
		updateStatus: func(_ context.Context, _ int64, status domain.PrepareTxStatus, _ map[string]any) error {
			assert.Equal(t, domain.PrepareTxReverted, status)
			return nil
		},
		waitForRcpt: func(_ context.Context, _ common.Hash) (*types.Receipt, error) {
			return &types.Receipt{
				Status:      types.ReceiptStatusFailed,
				BlockNumber: big.NewInt(100),
				GasUsed:     21000,
			}, nil
		},
		parseEvents: func(_ []*types.Log) []chain.ParsedEvent { return nil },
	}

	_, err := svc.Execute(context.Background(), 1)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "执行失败")
}
