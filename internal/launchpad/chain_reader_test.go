package launchpad

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockNodePool 是 chain.NodePool 的测试 mock。
type mockNodePool struct {
	callContractFn   func(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	receiptFn        func(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	filterLogsFn     func(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error)
	blockNumberFn    func(ctx context.Context) (uint64, error)
	headerByNumberFn func(ctx context.Context, number *big.Int) (*types.Header, error)
}

func (m *mockNodePool) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	if m.callContractFn != nil {
		return m.callContractFn(ctx, msg, blockNumber)
	}
	return nil, nil
}

func (m *mockNodePool) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	if m.filterLogsFn != nil {
		return m.filterLogsFn(ctx, query)
	}
	return nil, nil
}

func (m *mockNodePool) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	if m.receiptFn != nil {
		return m.receiptFn(ctx, txHash)
	}
	return nil, nil
}

func (m *mockNodePool) BlockNumber(ctx context.Context) (uint64, error) {
	if m.blockNumberFn != nil {
		return m.blockNumberFn(ctx)
	}
	return 0, nil
}

func (m *mockNodePool) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	if m.headerByNumberFn != nil {
		return m.headerByNumberFn(ctx, number)
	}
	return nil, nil
}

func (m *mockNodePool) SubscribeLogs(_ context.Context, _ ethereum.FilterQuery) (chan types.Log, ethereum.Subscription, error) {
	return nil, nil, nil
}

func (m *mockNodePool) Close() {}

// TestChainReader_GetTransactionReceipt 测试 receipt 查询场景。
func TestChainReader_GetTransactionReceipt(t *testing.T) {
	t.Run("查询成功的 receipt", func(t *testing.T) {
		pool := &mockNodePool{
			receiptFn: func(_ context.Context, _ common.Hash) (*types.Receipt, error) {
				return &types.Receipt{
					Status:      1,
					BlockNumber: big.NewInt(100),
					GasUsed:     0x5208,
				}, nil
			},
		}
		reader := NewChainReader(pool)

		receipt, err := reader.GetTransactionReceipt(context.Background(), "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
		require.NoError(t, err)
		require.NotNil(t, receipt)
		assert.Equal(t, uint64(1), receipt.Status)
		assert.Equal(t, uint64(100), receipt.BlockNumber)
		assert.Equal(t, uint64(0x5208), receipt.GasUsed)
	})

	t.Run("交易尚未被打包", func(t *testing.T) {
		pool := &mockNodePool{
			receiptFn: func(_ context.Context, _ common.Hash) (*types.Receipt, error) {
				return nil, ethereum.NotFound
			},
		}
		reader := NewChainReader(pool)

		receipt, err := reader.GetTransactionReceipt(context.Background(), "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
		require.NoError(t, err)
		assert.Nil(t, receipt)
	})

	t.Run("RPC 错误", func(t *testing.T) {
		pool := &mockNodePool{
			receiptFn: func(_ context.Context, _ common.Hash) (*types.Receipt, error) {
				return nil, fmt.Errorf("节点不可用")
			},
		}
		reader := NewChainReader(pool)

		_, err := reader.GetTransactionReceipt(context.Background(), "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
		assert.Error(t, err)
	})
}

// TestChainReader_GetUserTier 测试链上 Tier 查询。
func TestChainReader_GetUserTier(t *testing.T) {
	t.Run("查询用户 Tier 成功", func(t *testing.T) {
		pool := &mockNodePool{
			callContractFn: func(_ context.Context, _ ethereum.CallMsg, _ *big.Int) ([]byte, error) {
				result := make([]byte, 32)
				new(big.Int).SetInt64(2).FillBytes(result)
				return result, nil
			},
		}
		reader := NewChainReader(pool)

		tier, err := reader.GetUserTier(context.Background(), "0xTierContract", "0xUserAddress")
		require.NoError(t, err)
		assert.Equal(t, 2, tier)
	})

	t.Run("合约调用失败", func(t *testing.T) {
		pool := &mockNodePool{
			callContractFn: func(_ context.Context, _ ethereum.CallMsg, _ *big.Int) ([]byte, error) {
				return nil, fmt.Errorf("execution reverted")
			},
		}
		reader := NewChainReader(pool)

		_, err := reader.GetUserTier(context.Background(), "0xTierContract", "0xUserAddress")
		assert.Error(t, err)
	})
}

// TestChainReader_GetUserCredit 测试链上 Credit 查询。
func TestChainReader_GetUserCredit(t *testing.T) {
	t.Run("查询用户 Credit 成功", func(t *testing.T) {
		pool := &mockNodePool{
			callContractFn: func(_ context.Context, _ ethereum.CallMsg, _ *big.Int) ([]byte, error) {
				result := make([]byte, 32)
				new(big.Int).SetInt64(1e18).FillBytes(result)
				return result, nil
			},
		}
		reader := NewChainReader(pool)

		credit, err := reader.GetUserCredit(context.Background(), "0xTierContract", "0xUserAddress")
		require.NoError(t, err)
		assert.Equal(t, 0, credit.Cmp(big.NewInt(1e18)))
	})
}
