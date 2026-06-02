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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"

	"github.com/mousecake-go/mousecake-go/config"
)

// node 表示单个 RPC 节点，持有 ethclient、熔断器、限流器。
type node struct {
	name     string
	wsURL    string
	httpURL  string
	timeout  time.Duration
	chainID  int
	client   *ethclient.Client
	breakers map[string]*gobreaker.CircuitBreaker
	limiter  *rate.Limiter
	metrics  *rpcMetrics
	mu       sync.Mutex

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
	methods := []string{"eth_call", "eth_getLogs", "eth_getTransactionReceipt", "eth_blockNumber", "eth_subscribe", "eth_getBlockByNumber"}
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

	var limiter *rate.Limiter
	if cfg.RateLimit > 0 {
		limiter = rate.NewLimiter(rate.Limit(cfg.RateLimit), int(cfg.RateLimit))
	}

	return &node{
		name:     cfg.Name,
		wsURL:    cfg.WSURL,
		httpURL:  cfg.HTTPURL,
		timeout:  timeout,
		chainID:  chainID,
		client:   client,
		breakers: breakers,
		limiter:  limiter,
		metrics:  metrics,
		healthy:  true,
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

// subscribeLogs 订阅实时日志。
func (n *node) subscribeLogs(ctx context.Context, query ethereum.FilterQuery) (chan types.Log, ethereum.Subscription, error) {
	ch := make(chan types.Log, 128)
	sub, err := n.client.SubscribeFilterLogs(ctx, query, ch)
	if err != nil {
		return nil, nil, err
	}
	return ch, sub, nil
}

// close 关闭节点连接。
func (n *node) close() {
	if n.client != nil {
		n.client.Close()
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
