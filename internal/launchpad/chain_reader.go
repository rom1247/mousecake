// Package launchpad 实现 Launchpad IDO 募资模块的业务逻辑和数据访问。
// Package launchpad 实现 Launchpad IDO 募资模块的业务逻辑和数据访问。
package launchpad

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"

	"github.com/mousecake-go/mousecake-go/internal/chain"
	"github.com/mousecake-go/mousecake-go/internal/chain/contract/launchpad"
	"github.com/mousecake-go/mousecake-go/internal/launchpad/domain"
)

// ChainReader 实现 domain.ChainReader 接口，通过 NodePool 执行只读链上调用。
type ChainReader struct {
	pool    chain.NodePool
	tierAbi *launchpad.MouseTier
}

// NewChainReader 创建 ChainReader 实例。
// pool 参数提供多节点容错和熔断能力，由调用方管理生命周期。
func NewChainReader(pool chain.NodePool) *ChainReader {
	return &ChainReader{
		pool:    pool,
		tierAbi: launchpad.NewMouseTier(),
	}
}

// GetTransactionReceipt 查询交易 Receipt。
func (r *ChainReader) GetTransactionReceipt(ctx context.Context, txHash string) (*domain.ReceiptInfo, error) {
	hash := common.HexToHash(txHash)
	receipt, err := r.pool.TransactionReceipt(ctx, hash)
	if err != nil {
		if err == ethereum.NotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("查询 receipt %s: %w", txHash, err)
	}

	return &domain.ReceiptInfo{
		Status:      receipt.Status,
		BlockNumber: receipt.BlockNumber.Uint64(),
		GasUsed:     receipt.GasUsed,
	}, nil
}

// GetUserTier 查询用户在 MouseTier 合约的实时 Tier。
func (r *ChainReader) GetUserTier(ctx context.Context, mouseTierAddress string, userAddress string) (int, error) {
	contract := common.HexToAddress(mouseTierAddress)
	user := common.HexToAddress(userAddress)

	// 使用生成的合约绑定编码调用数据
	data := r.tierAbi.PackGetUserTier(user)

	result, err := r.callContract(ctx, contract, data)
	if err != nil {
		return 0, fmt.Errorf("调用 getUserTier(%s): %w", userAddress, err)
	}

	tier, err := r.tierAbi.UnpackGetUserTier(result)
	if err != nil {
		return 0, fmt.Errorf("解析 getUserTier 返回值: %w", err)
	}

	return int(tier.Int64()), nil
}

// GetUserCredit 查询用户在 MouseTier 合约的实时 Credit。
func (r *ChainReader) GetUserCredit(ctx context.Context, mouseTierAddress string, userAddress string) (*big.Int, error) {
	contract := common.HexToAddress(mouseTierAddress)
	user := common.HexToAddress(userAddress)

	// 使用生成的合约绑定编码调用数据
	data := r.tierAbi.PackGetUserCredit(user)

	result, err := r.callContract(ctx, contract, data)
	if err != nil {
		return nil, fmt.Errorf("调用 getUserCredit(%s): %w", userAddress, err)
	}

	credit, err := r.tierAbi.UnpackGetUserCredit(result)
	if err != nil {
		return nil, fmt.Errorf("解析 getUserCredit 返回值: %w", err)
	}

	return credit, nil
}

// callContract 执行合约 view 调用。
func (r *ChainReader) callContract(ctx context.Context, to common.Address, data []byte) ([]byte, error) {
	msg := ethereum.CallMsg{
		To:   &to,
		Data: data,
	}
	return r.pool.CallContract(ctx, msg, nil)
}
