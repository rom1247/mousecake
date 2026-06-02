package launchpad

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mousecake-go/mousecake-go/internal/launchpad/domain"
)

// --- Mock 实现 ---

type mockPrepareTxRepo struct {
	txs         map[int64]*domain.PrepareTx
	nextID      int64
	pendingHash map[string]*domain.PrepareTx
	err         error
}

func newMockPrepareTxRepo() *mockPrepareTxRepo {
	return &mockPrepareTxRepo{
		txs:         make(map[int64]*domain.PrepareTx),
		nextID:      1,
		pendingHash: make(map[string]*domain.PrepareTx),
	}
}

func (m *mockPrepareTxRepo) FindByID(_ context.Context, id int64) (*domain.PrepareTx, error) {
	if m.err != nil {
		return nil, m.err
	}
	tx, ok := m.txs[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return tx, nil
}

func (m *mockPrepareTxRepo) FindPendingByCalldataHash(_ context.Context, calldataHash string) (*domain.PrepareTx, error) {
	if m.err != nil {
		return nil, m.err
	}
	tx, ok := m.pendingHash[calldataHash]
	if !ok {
		return nil, nil
	}
	return tx, nil
}

func (m *mockPrepareTxRepo) FindByCaller(_ context.Context, _ string, _ string, _, _ int) ([]*domain.PrepareTx, int64, error) {
	return nil, 0, nil
}

func (m *mockPrepareTxRepo) FindBySaleID(_ context.Context, _ int64, _, _ int) ([]*domain.PrepareTx, int64, error) {
	return nil, 0, nil
}

func (m *mockPrepareTxRepo) FindPendingExpired(_ context.Context) ([]*domain.PrepareTx, error) {
	if m.err != nil {
		return nil, m.err
	}
	var result []*domain.PrepareTx
	for _, tx := range m.txs {
		if tx.Status == domain.PrepareTxPending && tx.IsExpired() {
			result = append(result, tx)
		}
	}
	return result, nil
}

func (m *mockPrepareTxRepo) FindBroadcastTimeout(_ context.Context, _ time.Duration) ([]*domain.PrepareTx, error) {
	if m.err != nil {
		return nil, m.err
	}
	var result []*domain.PrepareTx
	for _, tx := range m.txs {
		if tx.Status == domain.PrepareTxBroadcast && time.Since(tx.UpdatedAt) > 10*time.Minute {
			result = append(result, tx)
		}
	}
	return result, nil
}

func (m *mockPrepareTxRepo) Create(_ context.Context, tx *domain.PrepareTx) error {
	if m.err != nil {
		return m.err
	}
	tx.ID = m.nextID
	m.nextID++
	m.txs[tx.ID] = tx
	if tx.Status == domain.PrepareTxPending {
		m.pendingHash[tx.CalldataHash] = tx
	}
	return nil
}

func (m *mockPrepareTxRepo) UpdateStatus(_ context.Context, id int64, status domain.PrepareTxStatus, updates map[string]any) error {
	if m.err != nil {
		return m.err
	}
	tx, ok := m.txs[id]
	if !ok {
		return domain.ErrNotFound
	}
	// 从 pendingHash 移除旧记录
	if tx.Status == domain.PrepareTxPending {
		delete(m.pendingHash, tx.CalldataHash)
	}
	tx.Status = status
	tx.UpdatedAt = time.Now()
	for k, v := range updates {
		switch k {
		case "tx_hash":
			if s, ok := v.(string); ok {
				tx.TxHash = &s
			}
		case "block_number":
			if n, ok := v.(int64); ok {
				tx.BlockNumber = &n
			}
		case "error_message":
			if s, ok := v.(string); ok {
				tx.ErrorMessage = &s
			}
		case "confirmed_at":
			if t, ok := v.(time.Time); ok {
				tx.ConfirmedAt = &t
			}
		}
	}
	return nil
}

type mockChainReader struct {
	receipt *domain.ReceiptInfo
	tier    int
	credit  *big.Int
	err     error
}

func (m *mockChainReader) GetTransactionReceipt(_ context.Context, _ string) (*domain.ReceiptInfo, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.receipt, nil
}

func (m *mockChainReader) GetUserTier(_ context.Context, _ string, _ string) (int, error) {
	if m.err != nil {
		return 0, m.err
	}
	return m.tier, nil
}

func (m *mockChainReader) GetUserCredit(_ context.Context, _ string, _ string) (*big.Int, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.credit, nil
}

type mockABIEncoder struct {
	calldata []byte
	hash     string
	err      error
}

func (m *mockABIEncoder) encodeCall(_, _ string, _ ...any) ([]byte, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.calldata, nil
}

func (m *mockABIEncoder) calldataHash(_ []byte) string {
	return m.hash
}

// --- Prepare 交易创建测试（Task 6.1）---

func TestPrepareService_Create(t *testing.T) {
	t.Run("管理员创建 set_pool prepare", func(t *testing.T) {
		repo := newMockPrepareTxRepo()
		encoder := &mockABIEncoder{calldata: []byte{1, 2, 3, 4}, hash: "0xabc123"}
		chain := &mockChainReader{receipt: &domain.ReceiptInfo{Status: 1, BlockNumber: 100, GasUsed: 21000}}
		svc := NewPrepareService(repo, chain, newMockEncoder(encoder), 30*time.Minute)

		tx, err := svc.Create(context.Background(), CreatePrepareInput{
			OperationType: string(domain.OpSetPool),
			CallerAddress: "0xAdmin",
			SaleID:        int64Ptr(1),
			PoolIndex:     int64Ptr(0),
			Calldata:      []byte{1, 2, 3, 4},
		})
		require.NoError(t, err)
		assert.Equal(t, int64(1), tx.ID)
		assert.Equal(t, domain.PrepareTxPending, tx.Status)
		assert.Equal(t, "0xabc123", tx.CalldataHash)
	})

	t.Run("用户创建 deposit prepare", func(t *testing.T) {
		repo := newMockPrepareTxRepo()
		encoder := &mockABIEncoder{calldata: []byte{5, 6, 7, 8}, hash: "0xdef456"}
		chain := &mockChainReader{receipt: &domain.ReceiptInfo{Status: 1, BlockNumber: 100, GasUsed: 21000}}
		svc := NewPrepareService(repo, chain, newMockEncoder(encoder), 30*time.Minute)

		tx, err := svc.Create(context.Background(), CreatePrepareInput{
			OperationType: string(domain.OpDeposit),
			CallerAddress: "0xUser",
			SaleID:        int64Ptr(1),
			PoolIndex:     int64Ptr(0),
			Calldata:      []byte{5, 6, 7, 8},
		})
		require.NoError(t, err)
		assert.Equal(t, domain.PrepareTxPending, tx.Status)
		assert.Equal(t, "0xUser", tx.CallerAddress)
	})

	t.Run("同一用户重复创建相同参数的 prepare 返回已有记录", func(t *testing.T) {
		repo := newMockPrepareTxRepo()
		encoder := &mockABIEncoder{calldata: []byte{1, 2, 3}, hash: "0xduplicate"}
		chain := &mockChainReader{receipt: &domain.ReceiptInfo{Status: 1, BlockNumber: 100, GasUsed: 21000}}
		svc := NewPrepareService(repo, chain, newMockEncoder(encoder), 30*time.Minute)

		tx1, err := svc.Create(context.Background(), CreatePrepareInput{
			OperationType: string(domain.OpDeposit),
			CallerAddress: "0xUser",
			SaleID:        int64Ptr(1),
			PoolIndex:     int64Ptr(0),
			Calldata:      []byte{1, 2, 3},
		})
		require.NoError(t, err)

		tx2, err := svc.Create(context.Background(), CreatePrepareInput{
			OperationType: string(domain.OpDeposit),
			CallerAddress: "0xUser",
			SaleID:        int64Ptr(1),
			PoolIndex:     int64Ptr(0),
			Calldata:      []byte{1, 2, 3},
		})
		require.NoError(t, err)
		assert.Equal(t, tx1.ID, tx2.ID, "应返回已有的 pending 记录")
	})

	t.Run("无效操作类型", func(t *testing.T) {
		repo := newMockPrepareTxRepo()
		svc := NewPrepareService(repo, &mockChainReader{}, nil, 30*time.Minute)

		_, err := svc.Create(context.Background(), CreatePrepareInput{
			OperationType: "invalid_op",
			CallerAddress: "0xUser",
		})
		assert.Error(t, err)
	})
}

// --- Prepare txHash 提交和确认测试（Task 6.3）---

func TestPrepareService_Submit(t *testing.T) {
	t.Run("正常确认——receipt 成功", func(t *testing.T) {
		repo := newMockPrepareTxRepo()
		encoder := &mockABIEncoder{calldata: []byte{1}, hash: "0xhash"}
		chain := &mockChainReader{receipt: &domain.ReceiptInfo{Status: 1, BlockNumber: 100, GasUsed: 21000}}
		svc := NewPrepareService(repo, chain, newMockEncoder(encoder), 30*time.Minute)

		tx, _ := svc.Create(context.Background(), CreatePrepareInput{
			OperationType: string(domain.OpDeposit),
			CallerAddress: "0xUser",
			SaleID:        int64Ptr(1),
		})

		// 先标记为 signed 再 broadcast
		require.NoError(t, repo.txs[tx.ID].Transition(domain.PrepareTxSigned))
		require.NoError(t, repo.txs[tx.ID].Transition(domain.PrepareTxBroadcast))

		updated, err := svc.Submit(context.Background(), tx.ID, "0xTxHash123")
		require.NoError(t, err)
		assert.Equal(t, domain.PrepareTxConfirmed, updated.Status)
	})

	t.Run("链上 revert", func(t *testing.T) {
		repo := newMockPrepareTxRepo()
		encoder := &mockABIEncoder{calldata: []byte{1}, hash: "0xhash"}
		chain := &mockChainReader{
			receipt: &domain.ReceiptInfo{Status: 0, BlockNumber: 100, GasUsed: 21000},
		}
		svc := NewPrepareService(repo, chain, newMockEncoder(encoder), 30*time.Minute)

		tx, _ := svc.Create(context.Background(), CreatePrepareInput{
			OperationType: string(domain.OpDeposit),
			CallerAddress: "0xUser",
			SaleID:        int64Ptr(1),
		})
		require.NoError(t, repo.txs[tx.ID].Transition(domain.PrepareTxSigned))
		require.NoError(t, repo.txs[tx.ID].Transition(domain.PrepareTxBroadcast))

		updated, err := svc.Submit(context.Background(), tx.ID, "0xTxHashRevert")
		require.NoError(t, err)
		assert.Equal(t, domain.PrepareTxReverted, updated.Status)
	})

	t.Run("receipt 尚不可用——保持 broadcast", func(t *testing.T) {
		repo := newMockPrepareTxRepo()
		encoder := &mockABIEncoder{calldata: []byte{1}, hash: "0xhash"}
		chain := &mockChainReader{receipt: nil}
		svc := NewPrepareService(repo, chain, newMockEncoder(encoder), 30*time.Minute)

		tx, _ := svc.Create(context.Background(), CreatePrepareInput{
			OperationType: string(domain.OpDeposit),
			CallerAddress: "0xUser",
			SaleID:        int64Ptr(1),
		})
		require.NoError(t, repo.txs[tx.ID].Transition(domain.PrepareTxSigned))
		require.NoError(t, repo.txs[tx.ID].Transition(domain.PrepareTxBroadcast))

		updated, err := svc.Submit(context.Background(), tx.ID, "0xTxHashPending")
		require.NoError(t, err)
		assert.Equal(t, domain.PrepareTxBroadcast, updated.Status, "receipt 未就绪应保持 broadcast")
	})

	t.Run("submit 已过期的 prepare 失败", func(t *testing.T) {
		repo := newMockPrepareTxRepo()
		encoder := &mockABIEncoder{calldata: []byte{1}, hash: "0xhash"}
		svc := NewPrepareService(repo, &mockChainReader{}, newMockEncoder(encoder), 1*time.Nanosecond)

		tx, _ := svc.Create(context.Background(), CreatePrepareInput{
			OperationType: string(domain.OpDeposit),
			CallerAddress: "0xUser",
			SaleID:        int64Ptr(1),
		})
		// 让其过期
		repo.txs[tx.ID].ExpiresAt = time.Now().Add(-1 * time.Hour)

		_, err := svc.Submit(context.Background(), tx.ID, "0xTxHash")
		assert.Error(t, err)
	})

	t.Run("submit 非 broadcast 状态失败", func(t *testing.T) {
		repo := newMockPrepareTxRepo()
		encoder := &mockABIEncoder{calldata: []byte{1}, hash: "0xhash"}
		chain := &mockChainReader{receipt: &domain.ReceiptInfo{Status: 1, BlockNumber: 100, GasUsed: 21000}}
		svc := NewPrepareService(repo, chain, newMockEncoder(encoder), 30*time.Minute)

		tx, _ := svc.Create(context.Background(), CreatePrepareInput{
			OperationType: string(domain.OpDeposit),
			CallerAddress: "0xUser",
			SaleID:        int64Ptr(1),
		})
		// tx 仍处于 pending 状态

		_, err := svc.Submit(context.Background(), tx.ID, "0xTxHash")
		assert.Error(t, err)
	})
}

// --- 兜底轮询测试（Task 6.5）---

func TestPrepareService_PollTimeout(t *testing.T) {
	t.Run("pending 超时标记 expired", func(t *testing.T) {
		repo := newMockPrepareTxRepo()
		encoder := &mockABIEncoder{calldata: []byte{1}, hash: "0xhash"}
		svc := NewPrepareService(repo, &mockChainReader{}, newMockEncoder(encoder), 1*time.Nanosecond)

		tx, _ := svc.Create(context.Background(), CreatePrepareInput{
			OperationType: string(domain.OpDeposit),
			CallerAddress: "0xUser",
			SaleID:        int64Ptr(1),
		})
		repo.txs[tx.ID].ExpiresAt = time.Now().Add(-1 * time.Hour)

		err := svc.PollTimeout(context.Background())
		require.NoError(t, err)
		assert.Equal(t, domain.PrepareTxExpired, repo.txs[tx.ID].Status)
	})

	t.Run("broadcast 超时补查 receipt 成功", func(t *testing.T) {
		repo := newMockPrepareTxRepo()
		encoder := &mockABIEncoder{calldata: []byte{1}, hash: "0xhash"}
		chain := &mockChainReader{
			receipt: &domain.ReceiptInfo{Status: 1, BlockNumber: 200, GasUsed: 50000},
		}
		svc := NewPrepareService(repo, chain, newMockEncoder(encoder), 30*time.Minute)

		tx, _ := svc.Create(context.Background(), CreatePrepareInput{
			OperationType: string(domain.OpDeposit),
			CallerAddress: "0xUser",
			SaleID:        int64Ptr(1),
		})
		require.NoError(t, repo.txs[tx.ID].Transition(domain.PrepareTxSigned))
		require.NoError(t, repo.txs[tx.ID].Transition(domain.PrepareTxBroadcast))
		repo.txs[tx.ID].TxHash = strPtr("0xBroadcastTxHash")
		repo.txs[tx.ID].UpdatedAt = time.Now().Add(-15 * time.Minute)

		err := svc.PollTimeout(context.Background())
		require.NoError(t, err)
		assert.Equal(t, domain.PrepareTxConfirmed, repo.txs[tx.ID].Status)
	})
}

// --- 取消 Prepare 测试（Task 6.7）---

func TestPrepareService_Cancel(t *testing.T) {
	t.Run("取消 pending 的 prepare", func(t *testing.T) {
		repo := newMockPrepareTxRepo()
		encoder := &mockABIEncoder{calldata: []byte{1}, hash: "0xhash"}
		chain := &mockChainReader{receipt: &domain.ReceiptInfo{Status: 1, BlockNumber: 100, GasUsed: 21000}}
		svc := NewPrepareService(repo, chain, newMockEncoder(encoder), 30*time.Minute)

		tx, _ := svc.Create(context.Background(), CreatePrepareInput{
			OperationType: string(domain.OpDeposit),
			CallerAddress: "0xUser",
			SaleID:        int64Ptr(1),
		})

		err := svc.Cancel(context.Background(), tx.ID)
		require.NoError(t, err)
		assert.Equal(t, domain.PrepareTxExpired, repo.txs[tx.ID].Status)
	})

	t.Run("取消非 pending 状态失败", func(t *testing.T) {
		repo := newMockPrepareTxRepo()
		encoder := &mockABIEncoder{calldata: []byte{1}, hash: "0xhash"}
		chain := &mockChainReader{receipt: &domain.ReceiptInfo{Status: 1, BlockNumber: 100, GasUsed: 21000}}
		svc := NewPrepareService(repo, chain, newMockEncoder(encoder), 30*time.Minute)

		tx, _ := svc.Create(context.Background(), CreatePrepareInput{
			OperationType: string(domain.OpDeposit),
			CallerAddress: "0xUser",
			SaleID:        int64Ptr(1),
		})
		require.NoError(t, repo.txs[tx.ID].Transition(domain.PrepareTxSigned))
		require.NoError(t, repo.txs[tx.ID].Transition(domain.PrepareTxBroadcast))

		err := svc.Cancel(context.Background(), tx.ID)
		assert.Error(t, err)
	})
}

// --- 辅助函数（strPtr/int64Ptr 在 repository.go 中定义）---

// mockEncoder 适配 mockABIEncoder 到 EncoderInterface。
type mockEncoder struct {
	mock *mockABIEncoder
}

func newMockEncoder(m *mockABIEncoder) *mockEncoder { return &mockEncoder{mock: m} }

func (e *mockEncoder) EncodeCall(contractName, methodName string, args ...any) ([]byte, error) {
	return e.mock.encodeCall(contractName, methodName, args...)
}

func (e *mockEncoder) CalldataHash(calldata []byte) string {
	return e.mock.calldataHash(calldata)
}

// 确保 mockEncoder 实现 EncoderInterface
var _ EncoderInterface = (*mockEncoder)(nil)
