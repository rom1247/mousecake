package chain

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// Signer 交易签名服务，提供交易签名和广播能力。
type Signer struct {
	privateKey *ecdsa.PrivateKey
	address    common.Address
	chainID    *big.Int
	pool       NodePool
}

// NewSigner 创建 Signer 实例，解析 hex 私钥为 *ecdsa.PrivateKey，派生 Address。
func NewSigner(privateKeyHex string, pool NodePool, chainID *big.Int) (*Signer, error) {
	// 去除 0x 前缀（如果有）
	if strings.HasPrefix(privateKeyHex, "0x") || strings.HasPrefix(privateKeyHex, "0X") {
		privateKeyHex = privateKeyHex[2:]
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("解析私钥: %w", err)
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	return &Signer{
		privateKey: privateKey,
		address:    address,
		chainID:    chainID,
		pool:       pool,
	}, nil
}

// Address 返回 Signer 的钱包地址。
func (s *Signer) Address() common.Address {
	return s.address
}

// SignAndBroadcast 方法 — 查 nonce → 查 header（检测 1559）→ 查 gasPrice/tip → 估算 gasLimit → 构造交易 → 签名 → 广播 → 返回 txHash。
func (s *Signer) SignAndBroadcast(ctx context.Context, to common.Address, data []byte, value *big.Int) (common.Hash, error) {
	var txHash common.Hash

	slog.Info("正在签名广播交易...", "address", s.address.Hex())
	// 1. 查询 pending nonce
	nonce, err := s.pool.PendingNonceAt(ctx, s.address)
	if err != nil {
		return txHash, fmt.Errorf("查询 nonce: %w", err)
	}

	// 2. 查询最新区块头（用于检测 EIP-1559 支持）
	header, err := s.pool.HeaderByNumber(ctx, nil)
	if err != nil {
		return txHash, fmt.Errorf("查询区块头: %w", err)
	}

	// 3. 构造 CallMsg
	msg := ethereum.CallMsg{
		From:  s.address,
		To:    &to,
		Data:  data,
		Value: value,
	}

	// 4. 估算 gasLimit
	gasLimit, err := s.pool.EstimateGas(ctx, msg)
	if err != nil {
		return txHash, fmt.Errorf("估算 gas: %w", err)
	}

	// 5. 根据链是否支持 EIP-1559 构造交易
	var tx *types.Transaction
	if header.BaseFee != nil {
		// EIP-1559 交易
		tip, err := s.pool.SuggestGasPrice(ctx)
		if err != nil {
			return txHash, fmt.Errorf("查询 gas 价格: %w", err)
		}

		// 计算 maxFeePerGas = baseFee * 2 + tip
		baseFee := header.BaseFee
		maxFeePerGas := new(big.Int).Mul(baseFee, big.NewInt(2))
		maxFeePerGas = new(big.Int).Add(maxFeePerGas, tip)

		// 构造 EIP-1559 DynamicFee 交易
		tx = types.NewTx(&types.DynamicFeeTx{
			ChainID:   s.chainID,
			Nonce:     nonce,
			To:        &to,
			Value:     value,
			Gas:       gasLimit,
			GasTipCap: tip,
			GasFeeCap: maxFeePerGas,
			Data:      data,
		})
	} else {
		// Legacy 交易
		gasPrice, err := s.pool.SuggestGasPrice(ctx)
		if err != nil {
			return txHash, fmt.Errorf("查询 gas 价格: %w", err)
		}

		// 构造 Legacy 交易
		tx = types.NewTx(&types.LegacyTx{
			Nonce:    nonce,
			To:       &to,
			Value:    value,
			Gas:      gasLimit,
			GasPrice: gasPrice,
			Data:     data,
		})
	}

	// 6. 签名交易（LatestSignerForChainID 同时支持 Legacy 和 EIP-1559 交易类型）
	signer := types.LatestSignerForChainID(s.chainID)
	signedTx, err := types.SignTx(tx, signer, s.privateKey)
	if err != nil {
		return txHash, fmt.Errorf("签名交易: %w", err)
	}

	// 7. 广播交易
	if err := s.pool.SendTransaction(ctx, signedTx); err != nil {
		return txHash, fmt.Errorf("广播交易: %w", err)
	}

	txHash = signedTx.Hash()

	// 8. 记录审计日志
	slog.InfoContext(ctx, "签名广播交易成功",
		"to", to.Hex(),
		"value", value.String(),
		"nonce", nonce,
		"gas_limit", gasLimit,
		"tx_hash", txHash.Hex(),
	)

	return txHash, nil
}

// WaitForReceipt 轮询等待交易 receipt，直到交易上链或超时。
// 仅供开发环境使用。如果 ctx 没有 deadline，默认超时 60s，轮询间隔 1s。
func (s *Signer) WaitForReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	// 如果 ctx 没有 deadline，创建 60s 超时的子 context
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
		defer cancel()
	}

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		receipt, err := s.pool.TransactionReceipt(ctx, txHash)
		if err != nil && !errors.Is(err, ethereum.NotFound) {
			return nil, fmt.Errorf("查询 receipt %s: %w", txHash.Hex(), err)
		}
		if receipt != nil {
			return receipt, nil
		}

		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("等待 receipt 超时 (tx=%s): %w", txHash.Hex(), ctx.Err())
		case <-ticker.C:
			// 继续轮询
		}
	}
}
