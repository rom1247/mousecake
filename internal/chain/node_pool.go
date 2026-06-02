package chain

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/mousecake-go/mousecake-go/config"
)

// NodePool 定义 RPC 多节点池接口，提供统一的链上操作方法。
type NodePool interface {
	// CallContract 执行合约 view 调用。
	CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error)
	// FilterLogs 使用过滤器查询日志。
	FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error)
	// TransactionReceipt 查询交易 Receipt。
	TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error)
	// BlockNumber 查询最新区块号。
	BlockNumber(ctx context.Context) (uint64, error)
	// HeaderByNumber 按区块号查询区块头。
	HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error)
	// SubscribeLogs 订阅实时日志。
	SubscribeLogs(ctx context.Context, query ethereum.FilterQuery) (chan types.Log, ethereum.Subscription, error)
	// Close 关闭所有节点连接。
	Close()
}

// nodePool 是 NodePool 的实现。
type nodePool struct {
	nodes   []*node
	health  *healthChecker
	metrics *rpcMetrics
}

// NewNodePool 创建多节点池实例。
func NewNodePool(nodesCfg []config.ChainNodeConfig, chainID int) (NodePool, error) {
	if len(nodesCfg) == 0 {
		return nil, ErrAllNodesUnavailable
	}

	metrics := newRPCMetrics(chainID)

	// 按 priority 排序，确保优先级高的节点排在前面
	sorted := make([]config.ChainNodeConfig, len(nodesCfg))
	copy(sorted, nodesCfg)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Priority < sorted[j].Priority
	})

	nodes := make([]*node, 0, len(sorted))
	for _, cfg := range sorted {
		n, err := newNode(cfg, chainID, metrics)
		if err != nil {
			return nil, fmt.Errorf("创建节点 %s: %w", cfg.Name, err)
		}
		nodes = append(nodes, n)
	}

	np := &nodePool{
		nodes:   nodes,
		metrics: metrics,
	}

	np.health = newHealthChecker(nodes, np.metrics)

	return np, nil
}

// Start 启动健康检查。
func (np *nodePool) Start(ctx context.Context) {
	np.health.start(ctx)
}

// CallContract 执行合约 view 调用。
func (np *nodePool) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	for _, n := range np.nodes {
		if !n.isAvailable() {
			continue
		}
		result, err := n.execute(func(n *node) (any, error) {
			return n.callContract(ctx, msg, blockNumber)
		})
		if err != nil {
			if is429Error(err) {
				n.setRateLimited(time.Now().Add(5 * time.Second))
			}
			continue
		}
		return result.([]byte), nil
	}
	return nil, ErrAllNodesUnavailable
}

// FilterLogs 使用过滤器查询日志。
func (np *nodePool) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	for _, n := range np.nodes {
		if !n.isAvailable() {
			continue
		}
		result, err := n.execute(func(n *node) (any, error) {
			return n.filterLogs(ctx, query)
		})
		if err != nil {
			if is429Error(err) {
				n.setRateLimited(time.Now().Add(5 * time.Second))
			}
			continue
		}
		return result.([]types.Log), nil
	}
	return nil, ErrAllNodesUnavailable
}

// TransactionReceipt 查询交易 Receipt。
func (np *nodePool) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	for _, n := range np.nodes {
		if !n.isAvailable() {
			continue
		}
		result, err := n.execute(func(n *node) (any, error) {
			return n.transactionReceipt(ctx, txHash)
		})
		if err != nil {
			if is429Error(err) {
				n.setRateLimited(time.Now().Add(5 * time.Second))
			}
			continue
		}
		return result.(*types.Receipt), nil
	}
	return nil, ErrAllNodesUnavailable
}

// BlockNumber 查询最新区块号。
func (np *nodePool) BlockNumber(ctx context.Context) (uint64, error) {
	for _, n := range np.nodes {
		if !n.isAvailable() {
			continue
		}
		result, err := n.execute(func(n *node) (any, error) {
			return n.blockNumber(ctx)
		})
		if err != nil {
			if is429Error(err) {
				n.setRateLimited(time.Now().Add(5 * time.Second))
			}
			continue
		}
		return result.(uint64), nil
	}
	return 0, ErrAllNodesUnavailable
}

// HeaderByNumber 按区块号查询区块头。
func (np *nodePool) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	for _, n := range np.nodes {
		if !n.isAvailable() {
			continue
		}
		result, err := n.execute(func(n *node) (any, error) {
			return n.headerByNumber(ctx, number)
		})
		if err != nil {
			if is429Error(err) {
				n.setRateLimited(time.Now().Add(5 * time.Second))
			}
			continue
		}
		return result.(*types.Header), nil
	}
	return nil, ErrAllNodesUnavailable
}

// SubscribeLogs 订阅实时日志。
func (np *nodePool) SubscribeLogs(ctx context.Context, query ethereum.FilterQuery) (chan types.Log, ethereum.Subscription, error) {
	for _, n := range np.nodes {
		if !n.isAvailable() {
			continue
		}
		ch, sub, err := n.subscribeLogs(ctx, query)
		if err != nil {
			continue
		}
		return ch, sub, nil
	}
	return nil, nil, ErrAllNodesUnavailable
}

// Close 关闭所有节点连接。
func (np *nodePool) Close() {
	for _, n := range np.nodes {
		n.close()
	}
}

// connectWS 连接指定节点的 WebSocket 客户端。
func (np *nodePool) connectWS(ctx context.Context, nodeName string) (*ethclient.Client, error) {
	for _, n := range np.nodes {
		if n.name == nodeName && n.wsURL != "" {
			return ethclient.DialContext(ctx, n.wsURL)
		}
	}
	return nil, errors.New("chain: 未找到 WS 节点")
}
