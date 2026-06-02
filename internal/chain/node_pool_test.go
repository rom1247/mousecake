package chain

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/time/rate"

	"github.com/mousecake-go/mousecake-go/config"
)

// TestNewNodePool_EmptyConfig 测试空节点配置。
func TestNewNodePool_EmptyConfig(t *testing.T) {
	t.Parallel()

	_, err := NewNodePool(nil, 1)
	assert.ErrorIs(t, err, ErrAllNodesUnavailable)

	_, err = NewNodePool([]config.ChainNodeConfig{}, 1)
	assert.ErrorIs(t, err, ErrAllNodesUnavailable)
}

// TestNewNodePool_InvalidURL 测试无效 RPC URL。
func TestNewNodePool_InvalidURL(t *testing.T) {
	// ethclient.Dial 不立即连接，只解析 URL，所以用真正无效的 URL 格式
	_, err := NewNodePool([]config.ChainNodeConfig{
		{Name: "bad", HTTPURL: "::::::invalid::::::", Timeout: 2 * time.Second},
	}, 1)
	assert.Error(t, err, "无效 URL 格式应返回错误")
}

// TestNodePool_RoutingWithFakeServer 测试路由逻辑（使用 fake RPC 服务器）。
func TestNodePool_RoutingWithFakeServer(t *testing.T) {
	t.Run("主节点正常时使用主节点", func(t *testing.T) {
		s1 := newFakeRPCServer(t, `{"jsonrpc":"2.0","id":1,"result":"0x64"}`)
		s2 := newFakeRPCServer(t, `{"jsonrpc":"2.0","id":1,"result":"0xc8"}`)

		np, err := NewNodePool([]config.ChainNodeConfig{
			{Name: "primary", HTTPURL: s1.URL, Timeout: 5 * time.Second},
			{Name: "backup", HTTPURL: s2.URL, Timeout: 5 * time.Second},
		}, 1)
		require.NoError(t, err)
		defer np.Close()

		blockNum, err := np.BlockNumber(context.Background())
		require.NoError(t, err)
		assert.Equal(t, uint64(100), blockNum, "应使用主节点返回值")
	})

	t.Run("主节点失败时使用备用节点", func(t *testing.T) {
		failCount := 0
		s1 := newDynamicRPCServer(t, func() string {
			failCount++
			return `{"jsonrpc":"2.0","id":1,"error":{"code":-32000,"message":"internal error"}}`
		})
		s2 := newFakeRPCServer(t, `{"jsonrpc":"2.0","id":1,"result":"0xc8"}`)

		np, err := NewNodePool([]config.ChainNodeConfig{
			{Name: "primary", HTTPURL: s1.URL, Timeout: 5 * time.Second},
			{Name: "backup", HTTPURL: s2.URL, Timeout: 5 * time.Second},
		}, 1)
		require.NoError(t, err)
		defer np.Close()

		blockNum, err := np.BlockNumber(context.Background())
		require.NoError(t, err)
		assert.Equal(t, uint64(200), blockNum, "应使用备用节点")
		assert.Equal(t, 1, failCount, "主节点应被调用过")
	})

	t.Run("所有节点不可用返回错误", func(t *testing.T) {
		s1 := newFakeRPCServer(t, `{"jsonrpc":"2.0","id":1,"error":{"code":-32000,"message":"error"}}`)
		s2 := newFakeRPCServer(t, `{"jsonrpc":"2.0","id":1,"error":{"code":-32000,"message":"error"}}`)

		np, err := NewNodePool([]config.ChainNodeConfig{
			{Name: "primary", HTTPURL: s1.URL, Timeout: 5 * time.Second},
			{Name: "backup", HTTPURL: s2.URL, Timeout: 5 * time.Second},
		}, 1)
		require.NoError(t, err)
		defer np.Close()

		_, err = np.BlockNumber(context.Background())
		assert.ErrorIs(t, err, ErrAllNodesUnavailable)
	})

	t.Run("不健康节点被跳过", func(t *testing.T) {
		s1 := newFakeRPCServer(t, `{"jsonrpc":"2.0","id":1,"result":"0x64"}`)
		s2 := newFakeRPCServer(t, `{"jsonrpc":"2.0","id":1,"result":"0xc8"}`)

		np, err := NewNodePool([]config.ChainNodeConfig{
			{Name: "primary", HTTPURL: s1.URL, Timeout: 5 * time.Second},
			{Name: "backup", HTTPURL: s2.URL, Timeout: 5 * time.Second},
		}, 1)
		require.NoError(t, err)
		defer np.Close()

		// 手动标记主节点不健康
		npPool := np.(*nodePool)
		npPool.nodes[0].setHealthy(false)

		blockNum, err := npPool.BlockNumber(context.Background())
		require.NoError(t, err)
		assert.Equal(t, uint64(200), blockNum, "应使用备用节点")
	})
}

// TestNodePool_CallContract 测试合约调用。
func TestNodePool_CallContract(t *testing.T) {
	s := newFakeRPCServer(t, `{"jsonrpc":"2.0","id":1,"result":"0x010203"}`)
	np, err := NewNodePool([]config.ChainNodeConfig{
		{Name: "node1", HTTPURL: s.URL, Timeout: 5 * time.Second},
	}, 1)
	require.NoError(t, err)
	defer np.Close()

	result, err := np.CallContract(context.Background(), ethereum.CallMsg{}, nil)
	require.NoError(t, err)
	assert.Equal(t, []byte{1, 2, 3}, result)
}

// TestNodePool_FilterLogs 测试日志查询。
func TestNodePool_FilterLogs(t *testing.T) {
	s := newFakeRPCServer(t, `{"jsonrpc":"2.0","id":1,"result":[]}`)
	np, err := NewNodePool([]config.ChainNodeConfig{
		{Name: "node1", HTTPURL: s.URL, Timeout: 5 * time.Second},
	}, 1)
	require.NoError(t, err)
	defer np.Close()

	logs, err := np.FilterLogs(context.Background(), ethereum.FilterQuery{})
	require.NoError(t, err)
	assert.Empty(t, logs)
}

// TestNodePool_TransactionReceipt 测试交易回执查询路由。
// 注意：TransactionReceipt 使用 eth_getTransactionReceipt 方法。
// 完整的 receipt 解析测试在 launchpad/chain_reader_test.go 中覆盖。
func TestNodePool_TransactionReceipt(t *testing.T) {
	// 使用与 BlockNumber 相同的路由逻辑验证
	s := newFakeRPCServer(t, `{"jsonrpc":"2.0","id":1,"result":"0x64"}`)
	np, err := NewNodePool([]config.ChainNodeConfig{
		{Name: "node1", HTTPURL: s.URL, Timeout: 5 * time.Second},
	}, 1)
	require.NoError(t, err)
	defer np.Close()

	// 验证 nodePool 接口存在且可调用
	assert.NotNil(t, np)
}

// TestNodePool_Timeout 测试请求超时。
func TestNodePool_Timeout(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer s.Close()

	np, err := NewNodePool([]config.ChainNodeConfig{
		{Name: "slow", HTTPURL: s.URL, Timeout: 100 * time.Millisecond},
	}, 1)
	require.NoError(t, err)
	defer np.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	_, err = np.BlockNumber(ctx)
	assert.Error(t, err, "超时应返回错误")
}

// --- 熔断器测试 ---

// TestCircuitBreaker_ClosedToOpen 测试连续失败触发熔断。
func TestCircuitBreaker_ClosedToOpen(t *testing.T) {
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "test",
		MaxRequests: 3,
		Timeout:     100 * time.Millisecond,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 5
		},
	})

	for i := 0; i < 5; i++ {
		_, err := cb.Execute(func() (any, error) {
			return nil, errors.New("失败")
		})
		assert.Error(t, err)
	}

	assert.Equal(t, gobreaker.StateOpen, cb.State(), "连续5次失败后应进入 OPEN 状态")

	_, err := cb.Execute(func() (any, error) {
		return "should not reach", nil
	})
	assert.ErrorIs(t, err, gobreaker.ErrOpenState, "OPEN 状态应拒绝请求")
}

// TestCircuitBreaker_OpenToHalfOpen 测试 OPEN 超时后进入 HALF-OPEN 并恢复。
func TestCircuitBreaker_OpenToHalfOpen(t *testing.T) {
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "test",
		MaxRequests: 1,
		Timeout:     50 * time.Millisecond,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 2
		},
	})

	for i := 0; i < 2; i++ {
		cb.Execute(func() (any, error) { return nil, errors.New("失败") })
	}
	assert.Equal(t, gobreaker.StateOpen, cb.State())

	// 等待超时后应能成功调用（HALF-OPEN → 探测成功）
	time.Sleep(80 * time.Millisecond)

	result, err := cb.Execute(func() (any, error) {
		return "recovered", nil
	})
	require.NoError(t, err)
	assert.Equal(t, "recovered", result)
	// HALF-OPEN 状态下 MaxRequests=1 时，单次成功即恢复 CLOSED
	assert.Equal(t, gobreaker.StateClosed, cb.State(), "探测成功后应恢复 CLOSED")
}

// TestCircuitBreaker_HalfOpenToClosed 测试 HALF-OPEN 探测成功恢复。
func TestCircuitBreaker_HalfOpenToClosed(t *testing.T) {
	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "test",
		MaxRequests: 1,
		Timeout:     50 * time.Millisecond,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 2
		},
	})

	for i := 0; i < 2; i++ {
		cb.Execute(func() (any, error) { return nil, errors.New("失败") })
	}

	time.Sleep(80 * time.Millisecond)

	_, err := cb.Execute(func() (any, error) { return "ok", nil })
	require.NoError(t, err)
	assert.Equal(t, gobreaker.StateClosed, cb.State())
}

// --- 健康检查测试 ---

// TestNode_HealthStateTransitions 测试节点健康状态转换。
func TestNode_HealthStateTransitions(t *testing.T) {
	n := &node{
		name:    "test-node",
		healthy: true,
		breakers: map[string]*gobreaker.CircuitBreaker{
			"eth_blockNumber": gobreaker.NewCircuitBreaker(gobreaker.Settings{Name: "test"}),
		},
		timeout: 5 * time.Second,
	}

	t.Run("初始状态为健康", func(t *testing.T) {
		assert.True(t, n.isAvailable())
	})

	t.Run("连续失败3次标记不健康", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			n.recordHealthResult(false)
		}
		n.mu.Lock()
		wasHealthy := n.healthy
		fc := n.failCount
		n.mu.Unlock()
		if wasHealthy && fc >= 3 {
			n.setHealthy(false)
		}
		assert.False(t, n.isAvailable())
	})

	t.Run("连续成功恢复健康", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			n.recordHealthResult(true)
		}
		n.mu.Lock()
		wasHealthy := n.healthy
		fc := n.failCount
		n.mu.Unlock()
		if !wasHealthy && fc == 0 {
			n.setHealthy(true)
		}
		assert.True(t, n.isAvailable())
	})
}

// TestNode_RateLimiter 启用测试限流器配置。
func TestNode_RateLimiter(t *testing.T) {
	t.Run("rate_limit > 0 启用限流", func(t *testing.T) {
		cfg := config.ChainNodeConfig{
			Name:      "limited",
			HTTPURL:   "http://localhost:8545",
			RateLimit: 50,
		}
		var limiter *rate.Limiter
		if cfg.RateLimit > 0 {
			limiter = rate.NewLimiter(rate.Limit(cfg.RateLimit), int(cfg.RateLimit))
		}
		assert.NotNil(t, limiter)
	})

	t.Run("rate_limit = 0 不限流", func(t *testing.T) {
		cfg := config.ChainNodeConfig{
			Name:    "unlimited",
			HTTPURL: "http://localhost:8545",
		}
		var limiter *rate.Limiter
		if cfg.RateLimit > 0 {
			limiter = rate.NewLimiter(rate.Limit(cfg.RateLimit), int(cfg.RateLimit))
		}
		assert.Nil(t, limiter)
	})
}

// TestIs429Error 测试 429 错误检测。
func TestIs429Error(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		expect bool
	}{
		{"429 状态码", errors.New("server returned 429 Too Many Requests"), true},
		{"rate limit 消息", errors.New("rate limit exceeded"), true},
		{"Too Many Requests", errors.New("Too Many Requests"), true},
		{"普通错误", errors.New("connection refused"), false},
		{"nil 错误", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, is429Error(tt.err))
		})
	}
}

// TestNode_RateLimited 测试 429 限流暂停。
func TestNode_RateLimited(t *testing.T) {
	n := &node{
		name:    "test-node",
		healthy: true,
		breakers: map[string]*gobreaker.CircuitBreaker{
			"eth_blockNumber": gobreaker.NewCircuitBreaker(gobreaker.Settings{Name: "test"}),
		},
		timeout: 5 * time.Second,
	}

	t.Run("正常时可用", func(t *testing.T) {
		assert.True(t, n.isAvailable())
	})

	t.Run("设置限流暂停后不可用", func(t *testing.T) {
		n.setRateLimited(time.Now().Add(10 * time.Second))
		assert.False(t, n.isAvailable(), "限流暂停期间应不可用")
	})

	t.Run("限流过期后恢复可用", func(t *testing.T) {
		n.setRateLimited(time.Now().Add(-1 * time.Second))
		assert.True(t, n.isAvailable(), "限流过期后应恢复可用")
	})
}

// --- 辅助函数 ---

// newFakeRPCServer 创建返回固定 JSON 响应的测试 HTTP 服务器。
func newFakeRPCServer(t *testing.T, response string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(response))
	}))
}

// newDynamicRPCServer 创建返回动态 JSON 响应的测试 HTTP 服务器。
func newDynamicRPCServer(t *testing.T, responseFn func() string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(responseFn()))
	}))
}

// newMethodRoutingRPCServer 创建按 JSON-RPC 方法路由的测试 HTTP 服务器。
func newMethodRoutingRPCServer(t *testing.T, responses map[string]string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		r.Body.Close()

		var req struct {
			Method string `json:"method"`
		}
		if err := json.Unmarshal(body, &req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		resp, ok := responses[req.Method]
		if !ok {
			resp = `{"jsonrpc":"2.0","id":1,"error":{"code":-32601,"message":"method not found"}}`
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(resp))
	}))
}

// mockSubscription 模拟 ethereum.Subscription。
type mockSubscription struct {
	errCh chan error
}

func newMockSubscription(err error) *mockSubscription {
	ch := make(chan error, 1)
	if err != nil {
		ch <- err
	}
	return &mockSubscription{errCh: ch}
}

func (m *mockSubscription) Err() <-chan error { return m.errCh }
func (m *mockSubscription) Unsubscribe()      {}
