package chain

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
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
	// BatchCallContract 批量执行合约 view 调用，将多个 eth_call 合并为单次 JSON-RPC 2.0 批量请求。
	BatchCallContract(ctx context.Context, calls []ethereum.CallMsg) ([][]byte, error)
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
	// PendingNonceAt 查询指定地址的 pending nonce。
	PendingNonceAt(ctx context.Context, account common.Address) (uint64, error)
	// EstimateGas 估算交易 gas 用量。
	EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error)
	// SendTransaction 广播签名后的交易。
	SendTransaction(ctx context.Context, tx *types.Transaction) error
	// SuggestGasPrice 获取当前链建议的 gas 价格。
	SuggestGasPrice(ctx context.Context) (*big.Int, error)
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

	// 过滤掉被禁用的节点，保留启用的节点
	active := make([]config.ChainNodeConfig, 0, len(nodesCfg))
	for _, cfg := range nodesCfg {
		if cfg.IsEnabled() {
			active = append(active, cfg)
		} else {
			slog.Info("节点已禁用，跳过创建", "node", cfg.Name)
		}
	}
	if len(active) == 0 {
		return nil, ErrAllNodesUnavailable
	}

	// 按 priority 排序，确保优先级高的节点排在前面
	sort.Slice(active, func(i, j int) bool {
		return active[i].Priority < active[j].Priority
	})

	nodes := make([]*node, 0, len(active))
	for _, cfg := range active {
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

// tryNodes 依次尝试可用节点执行 fn，返回第一个成功的结果。
// 如果所有节点都失败，返回聚合错误（包含每个节点的具体错误），
// 并在日志中记录每个节点的失败原因，便于快速定位 RPC 问题。
func (np *nodePool) tryNodes(ctx context.Context, opName string, fn func(*node) (any, error)) (any, error) {
	var errs error
	for _, n := range np.nodes {
		if !n.isAvailable() {
			continue
		}
		result, err := n.execute(func(n *node) (any, error) {
			return fn(n)
		})
		if err != nil {
			slog.WarnContext(ctx, "节点请求失败", "operation", opName, "node", n.name, "err", err)
			if is429Error(err) {
				n.setRateLimited(time.Now().Add(5 * time.Second))
			}
			errs = errors.Join(errs, fmt.Errorf("节点 %s: %#v", n.name, err))
			continue
		}
		return result, nil
	}
	if errs != nil {
		slog.ErrorContext(ctx, "所有节点不可用", "operation", opName, "errors", errs)
		return nil, fmt.Errorf("chain: %s 所有节点不可用: %w: %w", opName, ErrAllNodesUnavailable, errs)
	}
	return nil, ErrAllNodesUnavailable
}

// CallContract 执行合约 view 调用。
func (np *nodePool) CallContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	result, err := np.tryNodes(ctx, "callContract", func(n *node) (any, error) {
		return n.callContract(ctx, msg, blockNumber)
	})
	if err != nil {
		return nil, err
	}
	return result.([]byte), nil
}

// BatchCallContract 批量执行合约 view 调用，复用现有节点选择和熔断机制。
func (np *nodePool) BatchCallContract(ctx context.Context, calls []ethereum.CallMsg) ([][]byte, error) {
	if len(calls) == 0 {
		return nil, nil
	}
	result, err := np.tryNodes(ctx, "batchCallContract", func(n *node) (any, error) {
		return n.batchCallContract(ctx, calls)
	})
	if err != nil {
		return nil, err
	}
	return result.([][]byte), nil
}

// FilterLogs 使用过滤器查询日志。
func (np *nodePool) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	result, err := np.tryNodes(ctx, "filterLogs", func(n *node) (any, error) {
		return n.filterLogs(ctx, query)
	})
	if err != nil {
		return nil, err
	}
	return result.([]types.Log), nil
}

// TransactionReceipt 查询交易 Receipt。
func (np *nodePool) TransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	result, err := np.tryNodes(ctx, "transactionReceipt", func(n *node) (any, error) {
		return n.transactionReceipt(ctx, txHash)
	})
	if err != nil {
		return nil, err
	}
	return result.(*types.Receipt), nil
}

// BlockNumber 查询最新区块号。
func (np *nodePool) BlockNumber(ctx context.Context) (uint64, error) {
	result, err := np.tryNodes(ctx, "blockNumber", func(n *node) (any, error) {
		return n.blockNumber(ctx)
	})
	if err != nil {
		return 0, err
	}
	return result.(uint64), nil
}

// HeaderByNumber 按区块号查询区块头。
func (np *nodePool) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	result, err := np.tryNodes(ctx, "headerByNumber", func(n *node) (any, error) {
		return n.headerByNumber(ctx, number)
	})
	if err != nil {
		return nil, err
	}
	return result.(*types.Header), nil
}

// SubscribeLogs 订阅实时日志。
func (np *nodePool) SubscribeLogs(ctx context.Context, query ethereum.FilterQuery) (chan types.Log, ethereum.Subscription, error) {
	slog.InfoContext(ctx, "订阅链上日志",
		"from", query.FromBlock, "to", query.ToBlock,
		"addresses", query.Addresses, "topics", query.Topics,
	)
	var errs error
	for _, n := range np.nodes {
		if !n.isAvailable() {
			continue
		}
		ch, sub, err := n.subscribeLogs(ctx, query)
		if err != nil {
			slog.WarnContext(ctx, "节点请求失败", "operation", "subscribeLogs", "node", n.name, "err", err)
			if is429Error(err) {
				n.setRateLimited(time.Now().Add(5 * time.Second))
			}
			errs = errors.Join(errs, fmt.Errorf("节点 %s: %#v", n.name, err))
			continue
		}
		slog.InfoContext(ctx, "订阅链上日志成功", "node", n.name)
		return ch, sub, nil
	}
	if errs != nil {
		slog.ErrorContext(ctx, "所有节点不可用", "operation", "subscribeLogs", "errors", fmt.Sprintf("%#v", errs))
		return nil, nil, fmt.Errorf("chain: subscribeLogs 所有节点不可用: %w: %w", ErrAllNodesUnavailable, errs)
	}
	return nil, nil, ErrAllNodesUnavailable
}

// PendingNonceAt 查询指定地址的 pending nonce。
// 写操作策略：只需优先级最高的节点成功即可，无需遍历所有节点。
func (np *nodePool) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	result, err := np.tryNodes(ctx, "pendingNonceAt", func(n *node) (any, error) {
		return n.pendingNonceAt(ctx, account)
	})
	if err != nil {
		return 0, err
	}
	return result.(uint64), nil
}

// EstimateGas 估算交易 gas 用量。
func (np *nodePool) EstimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	result, err := np.tryNodes(ctx, "estimateGas", func(n *node) (any, error) {
		return n.estimateGas(ctx, msg)
	})
	if err != nil {
		return 0, err
	}
	return result.(uint64), nil
}

// SendTransaction 广播签名后的交易。
func (np *nodePool) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	_, err := np.tryNodes(ctx, "sendTransaction", func(n *node) (any, error) {
		return nil, n.sendTransaction(ctx, tx)
	})
	return err
}

// SuggestGasPrice 获取当前链建议的 gas 价格。
func (np *nodePool) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	result, err := np.tryNodes(ctx, "suggestGasPrice", func(n *node) (any, error) {
		return n.suggestGasPrice(ctx)
	})
	if err != nil {
		return nil, err
	}
	return result.(*big.Int), nil
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
