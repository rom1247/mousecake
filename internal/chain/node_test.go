package chain

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mousecake-go/mousecake-go/config"
)

// Test_node_pendingNonceAt_成功获取nonce 测试 node.pendingNonceAt 方法成功获取 nonce。
func Test_node_pendingNonceAt_成功获取nonce(t *testing.T) {
	// GIVEN: node 已连接到 fake RPC 服务器
	// AND: 服务器返回 nonce=5
	s := newMethodRoutingRPCServer(t, map[string]string{
		"eth_getTransactionCount": `{"jsonrpc":"2.0","id":1,"result":"0x5"}`,
	})
	defer s.Close()

	n, err := newNode(config.ChainNodeConfig{
		Name:    "test-node",
		HTTPURL: s.URL,
		Timeout: 5 * time.Second,
	}, 1, nil)
	require.NoError(t, err)

	// WHEN: 调用 pendingNonceAt 查询地址的 nonce
	nonce, err := n.pendingNonceAt(context.Background(), common.HexToAddress("0xAdmin"))

	// THEN: 返回正确的 nonce 值
	require.NoError(t, err)
	assert.Equal(t, uint64(5), nonce)
}

// Test_node_estimateGas_成功估算gas 测试 node.estimateGas 方法成功估算 gas。
func Test_node_estimateGas_成功估算gas(t *testing.T) {
	// GIVEN: node 已连接到 fake RPC 服务器
	// AND: 交易参数合法
	s := newMethodRoutingRPCServer(t, map[string]string{
		"eth_estimateGas": `{"jsonrpc":"2.0","id":1,"result":"0x249f0"}`,
	})
	defer s.Close()

	n, err := newNode(config.ChainNodeConfig{
		Name:    "test-node",
		HTTPURL: s.URL,
		Timeout: 5 * time.Second,
	}, 1, nil)
	require.NoError(t, err)

	// WHEN: 调用 estimateGas 估算交易 gas
	to := common.HexToAddress("0xContract")
	gas, err := n.estimateGas(context.Background(), ethereum.CallMsg{
		From: common.HexToAddress("0xAdmin"),
		To:   &to,
		Data: common.Hex2Bytes("abc123"),
	})

	// THEN: 返回正确的 gas 用量（150000 = 0x249f0）
	require.NoError(t, err)
	assert.Equal(t, uint64(150000), gas)
}

// Test_node_sendTransaction_成功广播 测试 node.sendTransaction 方法成功广播交易。
func Test_node_sendTransaction_成功广播(t *testing.T) {
	// GIVEN: node 已连接到 fake RPC 服务器
	// AND: 交易已签名
	s := newMethodRoutingRPCServer(t, map[string]string{
		"eth_sendRawTransaction": `{"jsonrpc":"2.0","id":1,"result":"0xabc"}`,
	})
	defer s.Close()

	n, err := newNode(config.ChainNodeConfig{
		Name:    "test-node",
		HTTPURL: s.URL,
		Timeout: 5 * time.Second,
	}, 1, nil)
	require.NoError(t, err)

	signedTx := newTestSignedTx(t)

	// WHEN: 调用 sendTransaction 广播交易
	err = n.sendTransaction(context.Background(), signedTx)

	// THEN: 广播成功（无错误）
	require.NoError(t, err)
}

// Test_node_suggestGasPrice_成功获取gasPrice 测试 node.suggestGasPrice 方法成功获取 gas 价格。
func Test_node_suggestGasPrice_成功获取gasPrice(t *testing.T) {
	// GIVEN: node 已连接到 fake RPC 服务器
	// AND: 服务器返回 20 Gwei = 20000000000 = 0x4a817c800
	s := newMethodRoutingRPCServer(t, map[string]string{
		"eth_gasPrice": `{"jsonrpc":"2.0","id":1,"result":"0x4a817c800"}`,
	})
	defer s.Close()

	n, err := newNode(config.ChainNodeConfig{
		Name:    "test-node",
		HTTPURL: s.URL,
		Timeout: 5 * time.Second,
	}, 1, nil)
	require.NoError(t, err)

	// WHEN: 调用 suggestGasPrice 查询当前 gas 价格
	gasPrice, err := n.suggestGasPrice(context.Background())

	// THEN: 返回正确的 gas 价格（20 Gwei）
	require.NoError(t, err)
	assert.Equal(t, int64(20000000000), gasPrice.Int64())
}

// newTestSignedTx 创建一个用于测试的签名交易。
func newTestSignedTx(t *testing.T) *types.Transaction {
	t.Helper()
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	to := common.HexToAddress("0xABC")
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    1,
		To:       &to,
		Value:    nil,
		Gas:      21000,
		GasPrice: big.NewInt(1),
		Data:     nil,
	})
	chainID := big.NewInt(1)
	signer := types.NewEIP155Signer(chainID)
	signedTx, err := types.SignTx(tx, signer, privateKey)
	require.NoError(t, err)
	return signedTx
}
