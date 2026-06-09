package chain

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"

	"github.com/mousecake-go/mousecake-go/config"
)

// node 表示单个 RPC 节点，持有 ethclient、rpc.Client、熔断器、限流器。
type node struct {
	name      string
	wsURL     string
	httpURL   string
	timeout   time.Duration
	chainID   int
	client    *ethclient.Client
	wsClient  *ethclient.Client
	rpcClient *rpc.Client
	breakers  map[string]*gobreaker.CircuitBreaker
	limiter   *rate.Limiter
	metrics   *rpcMetrics
	mu        sync.Mutex

	// 健康状态
	healthy   bool
	failCount int

	// 限流暂停
	rateLimitedUntil time.Time
}

// newNode 创建单个 RPC 节点。
func newNode(cfg config.ChainNodeConfig, chainID int, metrics *rpcMetrics) (*node, error) {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	rpcURL := cfg.HTTPURL
	if rpcURL == "" {
		rpcURL = cfg.WSURL
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("连接 RPC 节点 %s 失败: %w", cfg.Name, err)
	}

	// 创建底层 rpc.Client 用于批量调用，优先使用 HTTP URL
	var rpcClient *rpc.Client
	httpURL := cfg.HTTPURL
	if httpURL != "" {
		rc, err := rpc.Dial(httpURL)
		if err != nil {
			client.Close()
			return nil, fmt.Errorf("创建 rpc.Client %s 失败: %w", cfg.Name, err)
		}
		rpcClient = rc
	} else if cfg.WSURL != "" {
		rc, err := rpc.Dial(cfg.WSURL)
		if err != nil {
			client.Close()
			return nil, fmt.Errorf("创建 rpc.Client (ws) %s 失败: %w", cfg.Name, err)
		}
		rpcClient = rc
	}

	failureThreshold := cfg.CircuitBreaker.FailureThreshold
	if failureThreshold == 0 {
		failureThreshold = 5
	}
	maxRequests := cfg.CircuitBreaker.MaxRequests
	if maxRequests == 0 {
		maxRequests = 3
	}
	cbTimeout := cfg.CircuitBreaker.Timeout
	if cbTimeout == 0 {
		cbTimeout = 30 * time.Second
	}

	breakers := make(map[string]*gobreaker.CircuitBreaker)
	methods := []string{"eth_call", "eth_getLogs", "eth_getTransactionReceipt", "eth_blockNumber", "eth_subscribe", "eth_getBlockByNumber", "eth_pendingNonce", "eth_estimateGas", "eth_sendRawTransaction", "eth_gasPrice"}
	for _, method := range methods {
		threshold := failureThreshold
		cbCfg := gobreaker.Settings{
			Name:        fmt.Sprintf("%s/%s", cfg.Name, method),
			MaxRequests: maxRequests,
			Interval:    0,
			Timeout:     cbTimeout,
			ReadyToTrip: func(counts gobreaker.Counts) bool {
				return counts.ConsecutiveFailures >= threshold
			},
		}
		breakers[method] = gobreaker.NewCircuitBreaker(cbCfg)
	}

	// 创建 WebSocket 客户端用于订阅（eth_subscribe 需要 WebSocket 连接）
	var wsClient *ethclient.Client
	if cfg.WSURL != "" {
		wsc, err := ethclient.Dial(cfg.WSURL)
		if err != nil {
			client.Close()
			if rpcClient != nil {
				rpcClient.Close()
			}
			return nil, fmt.Errorf("连接 WS 节点 %s 失败: %w", cfg.Name, err)
		}
		wsClient = wsc
	}

	var limiter *rate.Limiter
	if cfg.RateLimit > 0 {
		limiter = rate.NewLimiter(rate.Limit(cfg.RateLimit), int(cfg.RateLimit))
	}

	return &node{
		name:      cfg.Name,
		wsURL:     cfg.WSURL,
		httpURL:   cfg.HTTPURL,
		timeout:   timeout,
		chainID:   chainID,
		client:    client,
		wsClient:  wsClient,
		rpcClient: rpcClient,
		breakers:  breakers,
		limiter:   limiter,
		metrics:   metrics,
		healthy:   true,
	}, nil
}

// isAvailable 检查节点是否可用（健康且未被熔断或限流暂停）。
func (n *node) isAvailable() bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	if !n.healthy {
		return false
	}
	if !n.rateLimitedUntil.IsZero() && time.Now().Before(n.rateLimitedUntil) {
		return false
	}
	return true
}

// is429Error 检查错误是否为 429 Too Many Requests。
func is429Error(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return strings.Contains(msg, "429") ||
		strings.Contains(strings.ToLower(msg), "too many requests") ||
		strings.Contains(strings.ToLower(msg), "rate limit")
}

// setRateLimited 设置节点限流暂停截止时间。
func (n *node) setRateLimited(until time.Time) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.rateLimitedUntil = until
}

// execute 执行 RPC 调用，包装限流逻辑。
func (n *node) execute(fn func(*node) (any, error)) (any, error) {
	if n.limiter != nil {
		if err := n.limiter.Wait(context.Background()); err != nil {
			return nil, fmt.Errorf("限流等待: %w", err)
		}
	}
	return fn(n)
}

// callContract 执行合约 view 调用。
func (n *node) callContract(ctx context.Context, msg ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, n.timeout)
	defer cancel()

	result, err := n.breakerExecute("eth_call", func() (any, error) {
		return n.client.CallContract(ctx, msg, blockNumber)
	})
	if err != nil {
		return nil, err
	}
	return result.([]byte), nil
}

// batchCallContract 批量执行合约 view 调用，内部自动分批（每批最多 50 个请求）。
func (n *node) batchCallContract(ctx context.Context, calls []ethereum.CallMsg) ([][]byte, error) {
	if len(calls) == 0 {
		return nil, nil
	}
	if n.rpcClient == nil {
		return nil, ErrAllNodesUnavailable
	}

	const batchSize = 50
	var allResults [][]byte

	for i := 0; i < len(calls); i += batchSize {
		end := i + batchSize
		if end > len(calls) {
			end = len(calls)
		}
		batch := calls[i:end]

		results, err := n.executeBatch(ctx, batch)
		if err != nil {
			return nil, err
		}
		allResults = append(allResults, results...)
	}

	return allResults, nil
}

// executeBatch 执行单批 eth_call 请求。
func (n *node) executeBatch(ctx context.Context, calls []ethereum.CallMsg) ([][]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, n.timeout)
	defer cancel()

	elems := make([]rpc.BatchElem, len(calls))
	for i, call := range calls {
		var to string
		if call.To != nil {
			to = call.To.Hex()
		}
		elems[i] = rpc.BatchElem{
			Method: "eth_call",
			Args:   []any{map[string]string{"to": to, "data": fmt.Sprintf("0x%x", call.Data)}, "latest"},
			Result: new(string),
		}
	}

	_, err := n.breakerExecute("eth_call", func() (any, error) {
		return nil, n.rpcClient.BatchCallContext(ctx, elems)
	})
	if err != nil {
		return nil, err
	}

	results := make([][]byte, len(elems))
	for i, elem := range elems {
		if elem.Error != nil {
			results[i] = nil
			continue
		}
		hexStr, ok := elem.Result.(*string)
		if !ok || hexStr == nil || *hexStr == "" || *hexStr == "0x" {
			results[i] = nil
			continue
		}
		data, err := hexutil.Decode(*hexStr)
		if err != nil {
			results[i] = nil
			continue
		}
		results[i] = data
	}

	return results, nil
}

// filterLogs 使用过滤器查询日志。
func (n *node) filterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	ctx, cancel := context.WithTimeout(ctx, n.timeout)
	defer cancel()

	result, err := n.breakerExecute("eth_getLogs", func() (any, error) {
		return n.client.FilterLogs(ctx, query)
	})
	if err != nil {
		return nil, err
	}
	return result.([]types.Log), nil
}

// transactionReceipt 查询交易 Receipt。
func (n *node) transactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	ctx, cancel := context.WithTimeout(ctx, n.timeout)
	defer cancel()

	result, err := n.breakerExecute("eth_getTransactionReceipt", func() (any, error) {
		return n.client.TransactionReceipt(ctx, txHash)
	})
	if err != nil {
		return nil, err
	}
	return result.(*types.Receipt), nil
}

// blockNumber 查询最新区块号。
func (n *node) blockNumber(ctx context.Context) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, n.timeout)
	defer cancel()

	result, err := n.breakerExecute("eth_blockNumber", func() (any, error) {
		return n.client.BlockNumber(ctx)
	})
	if err != nil {
		return 0, err
	}
	return result.(uint64), nil
}

// headerByNumber 按区块号查询区块头。
func (n *node) headerByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	ctx, cancel := context.WithTimeout(ctx, n.timeout)
	defer cancel()

	result, err := n.breakerExecute("eth_getBlockByNumber", func() (any, error) {
		return n.client.HeaderByNumber(ctx, number)
	})
	if err != nil {
		return nil, err
	}
	return result.(*types.Header), nil
}

// subscribeLogs 通过 WebSocket 客户端订阅实时日志（eth_subscribe 需要 WebSocket 连接）。
func (n *node) subscribeLogs(ctx context.Context, query ethereum.FilterQuery) (chan types.Log, ethereum.Subscription, error) {
	if n.wsClient == nil {
		return nil, nil, errors.New("节点未配置 WebSocket URL，无法订阅日志")
	}
	ch := make(chan types.Log, 128)
	sub, err := n.wsClient.SubscribeFilterLogs(ctx, query, ch)
	if err != nil {
		return nil, nil, err
	}
	return ch, sub, nil
}

// pendingNonceAt 查询指定地址的 pending nonce。
func (n *node) pendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, n.timeout)
	defer cancel()

	result, err := n.breakerExecute("eth_pendingNonce", func() (any, error) {
		return n.client.PendingNonceAt(ctx, account)
	})
	if err != nil {
		return 0, err
	}
	return result.(uint64), nil
}

// estimateGas 估算交易 gas 用量。
func (n *node) estimateGas(ctx context.Context, msg ethereum.CallMsg) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, n.timeout)
	defer cancel()

	result, err := n.breakerExecute("eth_estimateGas", func() (any, error) {
		return n.client.EstimateGas(ctx, msg)
	})
	if err != nil {
		return 0, err
	}
	return result.(uint64), nil
}

// sendTransaction 广播签名后的交易。
func (n *node) sendTransaction(ctx context.Context, tx *types.Transaction) error {
	ctx, cancel := context.WithTimeout(ctx, n.timeout)
	defer cancel()

	_, err := n.breakerExecute("eth_sendRawTransaction", func() (any, error) {
		return nil, n.client.SendTransaction(ctx, tx)
	})
	return err
}

// suggestGasPrice 获取当前链建议的 gas 价格。
func (n *node) suggestGasPrice(ctx context.Context) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(ctx, n.timeout)
	defer cancel()

	result, err := n.breakerExecute("eth_gasPrice", func() (any, error) {
		return n.client.SuggestGasPrice(ctx)
	})
	if err != nil {
		return nil, err
	}
	return result.(*big.Int), nil
}

// close 关闭节点连接。
func (n *node) close() {
	if n.client != nil {
		n.client.Close()
	}
	if n.wsClient != nil {
		n.wsClient.Close()
	}
	if n.rpcClient != nil {
		n.rpcClient.Close()
	}
}

// breakerExecute 通过熔断器执行函数。
func (n *node) breakerExecute(method string, fn func() (any, error)) (any, error) {
	cb, ok := n.breakers[method]
	if !ok {
		if n.metrics != nil {
			n.metrics.recordRequest(n.name, method)
		}
		result, err := fn()
		if err != nil && n.metrics != nil {
			n.metrics.recordError(n.name, method)
		}
		return result, err
	}

	if n.metrics != nil {
		n.metrics.recordRequest(n.name, method)
	}

	result, err := cb.Execute(fn)
	if err != nil {
		if errors.Is(err, gobreaker.ErrOpenState) || errors.Is(err, gobreaker.ErrTooManyRequests) {
			return nil, fmt.Errorf("熔断器拒绝 %s: %w", method, err)
		}
		if n.metrics != nil {
			n.metrics.recordError(n.name, method)
		}
		return nil, err
	}
	return result, nil
}

// setHealthy 设置节点健康状态。
func (n *node) setHealthy(healthy bool) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.healthy = healthy
	if healthy {
		n.failCount = 0
	}
}

// recordHealthResult 记录健康检查结果。
func (n *node) recordHealthResult(success bool) {
	n.mu.Lock()
	defer n.mu.Unlock()
	if success {
		n.failCount = 0
	} else {
		n.failCount++
	}
}
