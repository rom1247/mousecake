package chain

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	// hardhat0PrivateKey Hardhat 网络 #0 账户私钥，仅用于测试。
	hardhat0PrivateKey = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	// hardhat0Address Hardhat 网络 #0 账户地址。
	hardhat0Address = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
)

// mockSignerPool 是 NodePool 接口的 mock 实现，每个方法通过函数字段配置行为。
type mockSignerPool struct {
	// callContractFunc 执行合约 view 调用的 mock 函数。
	callContractFunc func(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	// batchCallContractFunc 批量执行合约 view 调用的 mock 函数。
	batchCallContractFunc func(ctx context.Context, calls []ethereum.CallMsg) ([][]byte, error)
	// filterLogsFunc 使用过滤器查询日志的 mock 函数。
	filterLogsFunc func(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error)
	// transactionReceiptFunc 查询交易 Receipt 的 mock 函数。
	transactionReceiptFunc func(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	// blockNumberFunc 查询最新区块号的 mock 函数。
	blockNumberFunc func(ctx context.Context) (uint64, error)
	// headerByNumberFunc 按区块号查询区块头的 mock 函数。
	headerByNumberFunc func(ctx context.Context, number *big.Int) (*types.Header, error)
	// subscribeLogsFunc 订阅实时日志的 mock 函数。
	subscribeLogsFunc func(ctx context.Context, query ethereum.FilterQuery) (chan types.Log, ethereum.Subscription, error)
	// pendingNonceAtFunc 查询 pending nonce 的 mock 函数。
	pendingNonceAtFunc func(ctx context.Context, account common.Address) (uint64, error)
	// estimateGasFunc 估算交易 gas 的 mock 函数。
	estimateGasFunc func(ctx context.Context, msg ethereum.CallMsg) (uint64, error)
	// sendTransactionFunc 广播交易的 mock 函数。
	sendTransactionFunc func(ctx context.Context, tx *types.Transaction) error
	// suggestGasPriceFunc 获取建议 gas 价格的 mock 函数。
	suggestGasPriceFunc func(ctx context.Context) (*big.Int, error)
	// closeFunc 关闭连接的 mock 函数。
	closeFunc func()
}

func (m *mockSignerPool) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	if m.callContractFunc != nil {
		return m.callContractFunc(ctx, msg, blockNumber)
	}
	return nil, nil
}

func (m *mockSignerPool) BatchCallContract(ctx context.Context, calls []ethereum.CallMsg) ([][]byte, error) {
	if m.batchCallContractFunc != nil {
		return m.batchCallContractFunc(ctx, calls)
	}
	return nil, nil
}

func (m *mockSignerPool) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	if m.filterLogsFunc != nil {
		return m.filterLogsFunc(ctx, query)
	}
	return nil, nil
}

func (m *mockSignerPool) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	if m.transactionReceiptFunc != nil {
		return m.transactionReceiptFunc(ctx, txHash)
	}
	return nil, nil
}

func (m *mockSignerPool) BlockNumber(ctx context.Context) (uint64, error) {
	if m.blockNumberFunc != nil {
		return m.blockNumberFunc(ctx)
	}
	return 0, nil
}

func (m *mockSignerPool) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	if m.headerByNumberFunc != nil {
		return m.headerByNumberFunc(ctx, number)
	}
	return nil, nil
}

func (m *mockSignerPool) SubscribeLogs(ctx context.Context, query ethereum.FilterQuery) (chan types.Log, ethereum.Subscription, error) {
	if m.subscribeLogsFunc != nil {
		return m.subscribeLogsFunc(ctx, query)
	}
	return nil, nil, nil
}

func (m *mockSignerPool) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	if m.pendingNonceAtFunc != nil {
		return m.pendingNonceAtFunc(ctx, account)
	}
	return 0, nil
}

func (m *mockSignerPool) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	if m.estimateGasFunc != nil {
		return m.estimateGasFunc(ctx, msg)
	}
	return 0, nil
}

func (m *mockSignerPool) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	if m.sendTransactionFunc != nil {
		return m.sendTransactionFunc(ctx, tx)
	}
	return nil
}

func (m *mockSignerPool) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	if m.suggestGasPriceFunc != nil {
		return m.suggestGasPriceFunc(ctx)
	}
	return big.NewInt(0), nil
}

func (m *mockSignerPool) Close() {
	if m.closeFunc != nil {
		m.closeFunc()
	}
}

// newTestSigner 创建测试用 Signer 实例，使用 Hardhat #0 私钥和指定 chainID。
func newTestSigner(t *testing.T, pool NodePool, chainID *big.Int) *Signer {
	t.Helper()
	signer, err := NewSigner(hardhat0PrivateKey, pool, chainID)
	require.NoError(t, err)
	return signer
}

// TestNewSigner_有效私钥初始化 测试有效私钥初始化 Signer。
func TestNewSigner_有效私钥初始化(t *testing.T) {
	// GIVEN: 提供有效的 64 字符 hex 编码私钥字符串
	// AND: 提供 NodePool 实例和 chainID
	pool := &mockSignerPool{}
	chainID := big.NewInt(31337)

	// WHEN: 调用 NewSigner(privateKeyHex, pool, chainID)
	signer, err := NewSigner(hardhat0PrivateKey, pool, chainID)

	// THEN: 返回 Signer 实例，无错误
	require.NoError(t, err)
	require.NotNil(t, signer)

	// AND: signer.Address() 返回正确的以太坊地址
	assert.Equal(t, hardhat0Address, signer.Address().Hex())
}

// TestNewSigner_无效私钥 测试无效私钥时返回错误。
func TestNewSigner_无效私钥(t *testing.T) {
	// GIVEN: 提供无效的私钥字符串
	pool := &mockSignerPool{}
	chainID := big.NewInt(31337)

	// WHEN: 调用 NewSigner("not-a-valid-hex", pool, chainID)
	signer, err := NewSigner("not-a-valid-hex", pool, chainID)

	// THEN: 返回错误
	require.Error(t, err)
	assert.Nil(t, signer)

	// AND: 错误信息包含 "解析私钥"
	assert.Contains(t, err.Error(), "解析私钥")
}

// TestSignAndBroadcast_Legacy路径 测试成功签名广播 Legacy 交易。
func TestSignAndBroadcast_Legacy路径(t *testing.T) {
	// GIVEN: Signer 已初始化
	// AND: 链不支持 EIP-1559（BaseFee 为 nil）
	pool := &mockSignerPool{
		pendingNonceAtFunc: func(_ context.Context, _ common.Address) (uint64, error) {
			return 5, nil
		},
		headerByNumberFunc: func(_ context.Context, _ *big.Int) (*types.Header, error) {
			return &types.Header{BaseFee: nil}, nil
		},
		suggestGasPriceFunc: func(_ context.Context) (*big.Int, error) {
			// 2 Gwei = 2 * 10^9
			return new(big.Int).Mul(big.NewInt(2), new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil)), nil
		},
		estimateGasFunc: func(_ context.Context, _ ethereum.CallMsg) (uint64, error) {
			return 21000, nil
		},
		sendTransactionFunc: func(_ context.Context, _ *types.Transaction) error {
			return nil
		},
	}
	chainID := big.NewInt(31337)
	signer := newTestSigner(t, pool, chainID)

	to := common.HexToAddress("0x1234567890123456789012345678901234567890")
	data := []byte{0xab, 0xcd}
	value := big.NewInt(1000000000000000000) // 1 ETH

	// WHEN: 调用 SignAndBroadcast(ctx, to, data, value)
	txHash, err := signer.SignAndBroadcast(context.Background(), to, data, value)

	// THEN: 返回非零 txHash，无错误
	require.NoError(t, err)
	assert.NotEqual(t, common.Hash{}, txHash)
}

// TestSignAndBroadcast_EIP1559路径 测试成功签名广播 EIP-1559 交易。
func TestSignAndBroadcast_EIP1559路径(t *testing.T) {
	// GIVEN: Signer 已初始化
	// AND: 链支持 EIP-1559（BaseFee 不为 nil）
	gwei := new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil)
	baseFee := new(big.Int).Mul(big.NewInt(30), gwei) // 30 Gwei

	pool := &mockSignerPool{
		pendingNonceAtFunc: func(_ context.Context, _ common.Address) (uint64, error) {
			return 10, nil
		},
		headerByNumberFunc: func(_ context.Context, _ *big.Int) (*types.Header, error) {
			return &types.Header{BaseFee: baseFee}, nil
		},
		suggestGasPriceFunc: func(_ context.Context) (*big.Int, error) {
			// 2 Gwei tip
			return new(big.Int).Mul(big.NewInt(2), gwei), nil
		},
		estimateGasFunc: func(_ context.Context, _ ethereum.CallMsg) (uint64, error) {
			return 150000, nil
		},
		sendTransactionFunc: func(_ context.Context, _ *types.Transaction) error {
			return nil
		},
	}
	chainID := big.NewInt(31337)
	signer := newTestSigner(t, pool, chainID)

	to := common.HexToAddress("0x1234567890123456789012345678901234567890")
	data := []byte{0xab, 0xcd}
	value := big.NewInt(5000000000000000000) // 5 ETH

	// WHEN: 调用 SignAndBroadcast(ctx, to, data, value)
	txHash, err := signer.SignAndBroadcast(context.Background(), to, data, value)

	// THEN: 返回非零 txHash，无错误
	require.NoError(t, err)
	assert.NotEqual(t, common.Hash{}, txHash)
}

// TestSignAndBroadcast_nonce查询失败 测试 nonce 查询失败时返回包装错误。
func TestSignAndBroadcast_nonce查询失败(t *testing.T) {
	// GIVEN: Signer 已初始化
	// AND: PendingNonceAt 返回错误（如连接超时）
	pool := &mockSignerPool{
		pendingNonceAtFunc: func(_ context.Context, _ common.Address) (uint64, error) {
			return 0, errors.New("connection refused")
		},
	}
	chainID := big.NewInt(31337)
	signer := newTestSigner(t, pool, chainID)

	to := common.HexToAddress("0x1234567890123456789012345678901234567890")

	// WHEN: 调用 SignAndBroadcast
	txHash, err := signer.SignAndBroadcast(context.Background(), to, nil, nil)

	// THEN: 返回包装后的错误
	require.Error(t, err)
	assert.Equal(t, common.Hash{}, txHash)

	// AND: 错误信息包含 "查询 nonce"
	assert.Contains(t, err.Error(), "查询 nonce")
}

// TestSignAndBroadcast_gas估算失败 测试 gas 估算失败时返回包装错误。
func TestSignAndBroadcast_gas估算失败(t *testing.T) {
	// GIVEN: Signer 已初始化
	// AND: nonce 查询成功
	// AND: EstimateGas 返回错误（如 invalid address）
	pool := &mockSignerPool{
		pendingNonceAtFunc: func(_ context.Context, _ common.Address) (uint64, error) {
			return 5, nil
		},
		headerByNumberFunc: func(_ context.Context, _ *big.Int) (*types.Header, error) {
			return &types.Header{BaseFee: nil}, nil
		},
		estimateGasFunc: func(_ context.Context, _ ethereum.CallMsg) (uint64, error) {
			return 0, errors.New("invalid address")
		},
	}
	chainID := big.NewInt(31337)
	signer := newTestSigner(t, pool, chainID)

	to := common.HexToAddress("0x1234567890123456789012345678901234567890")

	// WHEN: 调用 SignAndBroadcast
	txHash, err := signer.SignAndBroadcast(context.Background(), to, nil, nil)

	// THEN: 返回包装后的错误
	require.Error(t, err)
	assert.Equal(t, common.Hash{}, txHash)

	// AND: 错误信息包含 "估算 gas"
	assert.Contains(t, err.Error(), "估算 gas")
}

// TestSignAndBroadcast_广播失败 测试广播失败时返回包装错误。
func TestSignAndBroadcast_广播失败(t *testing.T) {
	// GIVEN: Signer 已初始化
	// AND: nonce 和 gas 估算成功
	// AND: SendTransaction 返回错误（如 nonce too low）
	pool := &mockSignerPool{
		pendingNonceAtFunc: func(_ context.Context, _ common.Address) (uint64, error) {
			return 5, nil
		},
		headerByNumberFunc: func(_ context.Context, _ *big.Int) (*types.Header, error) {
			return &types.Header{BaseFee: nil}, nil
		},
		suggestGasPriceFunc: func(_ context.Context) (*big.Int, error) {
			return new(big.Int).Mul(big.NewInt(2), new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil)), nil
		},
		estimateGasFunc: func(_ context.Context, _ ethereum.CallMsg) (uint64, error) {
			return 21000, nil
		},
		sendTransactionFunc: func(_ context.Context, _ *types.Transaction) error {
			return errors.New("nonce too low")
		},
	}
	chainID := big.NewInt(31337)
	signer := newTestSigner(t, pool, chainID)

	to := common.HexToAddress("0x1234567890123456789012345678901234567890")

	// WHEN: 调用 SignAndBroadcast
	txHash, err := signer.SignAndBroadcast(context.Background(), to, nil, nil)

	// THEN: 返回包装后的错误
	require.Error(t, err)
	assert.Equal(t, common.Hash{}, txHash)

	// AND: 错误信息包含 "广播交易"
	assert.Contains(t, err.Error(), "广播交易")
}

// TestSignAndBroadcast_签名审计日志 测试签名成功时记录审计日志。
func TestSignAndBroadcast_签名审计日志(t *testing.T) {
	// GIVEN: Signer 已初始化
	// AND: 所有 mock 调用成功
	pool := &mockSignerPool{
		pendingNonceAtFunc: func(_ context.Context, _ common.Address) (uint64, error) {
			return 5, nil
		},
		headerByNumberFunc: func(_ context.Context, _ *big.Int) (*types.Header, error) {
			return &types.Header{BaseFee: nil}, nil
		},
		suggestGasPriceFunc: func(_ context.Context) (*big.Int, error) {
			return new(big.Int).Mul(big.NewInt(2), new(big.Int).Exp(big.NewInt(10), big.NewInt(9), nil)), nil
		},
		estimateGasFunc: func(_ context.Context, _ ethereum.CallMsg) (uint64, error) {
			return 21000, nil
		},
		sendTransactionFunc: func(_ context.Context, _ *types.Transaction) error {
			return nil
		},
	}
	chainID := big.NewInt(31337)
	signer := newTestSigner(t, pool, chainID)

	to := common.HexToAddress("0x1234567890123456789012345678901234567890")
	value := big.NewInt(1000000000000000000) // 1 ETH

	// WHEN: 调用 SignAndBroadcast，方法内部通过 slog.InfoContext 记录审计日志
	txHash, err := signer.SignAndBroadcast(context.Background(), to, nil, value)

	// THEN: 方法正常返回，无错误（日志验证不在测试范围内）
	require.NoError(t, err)
	assert.NotEqual(t, common.Hash{}, txHash)
}
