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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

// mockPrepareTxRepo 是 PrepareTxRepository 查询方法的 mock。
// 为了测试 DevExecuteService，我们通过 interface 解耦。
type mockPrepareTxReader struct {
	tx  *domain.PrepareTx
	err error
}

type mockPrepareTxWriter struct {
	err error
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

	svc := &DevExecuteService{
		findByID: func(_ context.Context, id int64) (*domain.PrepareTx, error) {
			assert.Equal(t, int64(1), id)
			return ptx, nil
		},
		signer: &mockTxSigner{hash: expectedHash},
		updateStatus: func(_ context.Context, id int64, status domain.PrepareTxStatus, updates map[string]any) error {
			assert.Equal(t, int64(1), id)
			assert.Equal(t, domain.PrepareTxBroadcast, status)
			assert.Equal(t, expectedHash.Hex(), updates["tx_hash"])
			return nil
		},
	}

	txHash, err := svc.Execute(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, expectedHash.Hex(), txHash)
}

func TestDevExecuteService_Execute_PrepareTx不存在(t *testing.T) {
	svc := &DevExecuteService{
		findByID: func(_ context.Context, _ int64) (*domain.PrepareTx, error) {
			return nil, domain.ErrNotFound
		},
		signer:       &mockTxSigner{},
		updateStatus: func(_ context.Context, _ int64, _ domain.PrepareTxStatus, _ map[string]any) error { return nil },
	}

	_, err := svc.Execute(context.Background(), 999)
	require.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrNotFound))
}

func TestDevExecuteService_Execute_PrepareTx非pending状态(t *testing.T) {
	ptx := buildValidPendingPrepareTx()
	ptx.Status = domain.PrepareTxBroadcast

	svc := &DevExecuteService{
		findByID: func(_ context.Context, _ int64) (*domain.PrepareTx, error) {
			return ptx, nil
		},
		signer:       &mockTxSigner{},
		updateStatus: func(_ context.Context, _ int64, _ domain.PrepareTxStatus, _ map[string]any) error { return nil },
	}

	_, err := svc.Execute(context.Background(), 1)
	require.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "只能执行 pending"))
}

func TestDevExecuteService_Execute_PrepareTx已过期(t *testing.T) {
	ptx := buildValidPendingPrepareTx()
	ptx.ExpiresAt = time.Now().Add(-1 * time.Hour)

	svc := &DevExecuteService{
		findByID: func(_ context.Context, _ int64) (*domain.PrepareTx, error) {
			return ptx, nil
		},
		signer:       &mockTxSigner{},
		updateStatus: func(_ context.Context, _ int64, _ domain.PrepareTxStatus, _ map[string]any) error { return nil },
	}

	_, err := svc.Execute(context.Background(), 1)
	require.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "已过期"))
}

func TestDevExecuteService_Execute_缺少target_address(t *testing.T) {
	ptx := buildValidPendingPrepareTx()
	ptx.TargetAddress = ""

	svc := &DevExecuteService{
		findByID: func(_ context.Context, _ int64) (*domain.PrepareTx, error) {
			return ptx, nil
		},
		signer:       &mockTxSigner{},
		updateStatus: func(_ context.Context, _ int64, _ domain.PrepareTxStatus, _ map[string]any) error { return nil },
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
	}

	_, err := svc.Execute(context.Background(), 1)
	require.Error(t, err)
	assert.True(t, errors.Is(err, signErr))
	// 确认状态保持 pending
	assert.Equal(t, domain.PrepareTxPending, ptx.Status)
}
