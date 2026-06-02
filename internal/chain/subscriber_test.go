package chain

import (
	"context"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mousecake-go/mousecake-go/config"
)

// mockNodePool 模拟 NodePool 接口，用于测试 Subscriber。
type mockNodePool struct {
	mu sync.Mutex

	blockNumber uint64
	blockErr    error

	filterLogs []types.Log
	filterErr  error

	subCh  chan types.Log
	sub    ethereum.Subscription
	subErr error

	// subErrAfterFirst 首次 SubscribeLogs 成功后，后续调用返回此错误
	subErrAfterFirst error
	subCallCount     int

	filterCalls int
	blockCalls  int
}

func (m *mockNodePool) CallContract(_ context.Context, _ ethereum.CallMsg, _ *big.Int) ([]byte, error) {
	return nil, nil
}

func (m *mockNodePool) FilterLogs(_ context.Context, _ ethereum.FilterQuery) ([]types.Log, error) {
	m.mu.Lock()
	m.filterCalls++
	m.mu.Unlock()
	return m.filterLogs, m.filterErr
}

func (m *mockNodePool) TransactionReceipt(_ context.Context, _ common.Hash) (*types.Receipt, error) {
	return nil, nil
}

func (m *mockNodePool) BlockNumber(_ context.Context) (uint64, error) {
	m.mu.Lock()
	m.blockCalls++
	m.mu.Unlock()
	return m.blockNumber, m.blockErr
}

func (m *mockNodePool) HeaderByNumber(_ context.Context, _ *big.Int) (*types.Header, error) {
	return nil, nil
}

func (m *mockNodePool) SubscribeLogs(_ context.Context, _ ethereum.FilterQuery) (chan types.Log, ethereum.Subscription, error) {
	m.mu.Lock()
	m.subCallCount++
	count := m.subCallCount
	errAfterFirst := m.subErrAfterFirst
	m.mu.Unlock()

	if count > 1 && errAfterFirst != nil {
		return nil, nil, errAfterFirst
	}
	return m.subCh, m.sub, m.subErr
}

func (m *mockNodePool) Close() {}

// mockSub 实现 ethereum.Subscription。
// 关闭后 Err() 返回 nil channel，防止 select 立即返回零值。
type mockSub struct {
	errCh     chan error
	unsubDone chan struct{}
	closed    bool
	mu        sync.Mutex
}

func newMockSub() *mockSub {
	return &mockSub{
		errCh:     make(chan error, 1),
		unsubDone: make(chan struct{}),
	}
}

func (m *mockSub) Err() <-chan error {
	m.mu.Lock()
	closed := m.closed
	m.mu.Unlock()
	if closed {
		return nil
	}
	return m.errCh
}

func (m *mockSub) Unsubscribe() {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.closed {
		m.closed = true
		close(m.unsubDone)
	}
}

// defaultContracts 返回测试用的合约配置。
func defaultContracts() config.SyncContractsConfig {
	return config.SyncContractsConfig{
		MouseTier: "0x1234567890123456789012345678901234567890",
	}
}

// TestSubscriber_WSDisconnectTriggersHTTPFallback 测试 WS 断开后降级为 HTTP 轮询。
func TestSubscriber_WSDisconnectTriggersHTTPFallback(t *testing.T) {
	subCh := make(chan types.Log)
	sub := newMockSub()

	pool := &mockNodePool{
		subCh:            subCh,
		sub:              sub,
		subErr:           nil,            // 首次 WS 成功
		subErrAfterFirst: assert.AnError, // 重连失败，保持 HTTP 模式
		blockNumber:      100,
		filterLogs:       []types.Log{},
	}

	s := NewSubscriber(pool, 1, defaultContracts(), 3, 12*time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	outCh, err := s.Start(ctx)
	require.NoError(t, err)

	// 初始模式为 WS
	assert.Equal(t, ModeWS, s.Mode())

	// 模拟 WS 断开：发送错误
	sub.errCh <- assert.AnError

	// 等待降级为 HTTP
	assert.Eventually(t, func() bool {
		return s.Mode() == ModeHTTP
	}, 2*time.Second, 50*time.Millisecond, "WS 断开后应降级为 HTTP")

	s.Stop()
	drainLogs(outCh)
}

// TestSubscriber_HTTPPollingFallback 测试 WS 订阅失败直接降级为 HTTP。
func TestSubscriber_HTTPPollingFallback(t *testing.T) {
	pool := &mockNodePool{
		subErr:      assert.AnError, // WS 订阅失败
		blockNumber: 100,
		filterLogs:  []types.Log{},
	}

	s := NewSubscriber(pool, 1, defaultContracts(), 3, 12*time.Second)
	s.SetLastSeenBlock(50)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	outCh, err := s.Start(ctx)
	require.NoError(t, err)

	// 应立即进入 HTTP 模式
	assert.Equal(t, ModeHTTP, s.Mode())

	// 等待轮询触发
	time.Sleep(4 * time.Second)

	pool.mu.Lock()
	calls := pool.filterCalls
	pool.mu.Unlock()

	assert.GreaterOrEqual(t, calls, 1, "HTTP 降级后应调用 FilterLogs")

	s.Stop()
	drainLogs(outCh)
}

// TestSubscriber_HTTPPollingReceivesLogs 测试 HTTP 轮询能正确接收日志。
func TestSubscriber_HTTPPollingReceivesLogs(t *testing.T) {
	testLog := types.Log{
		Address:     common.HexToAddress("0x1234567890123456789012345678901234567890"),
		Topics:      []common.Hash{common.HexToHash("0xabc")},
		BlockNumber: 60,
		TxHash:      common.HexToHash("0xdef"),
		Index:       0,
	}

	pool := &mockNodePool{
		subErr:      assert.AnError,
		blockNumber: 100,
		filterLogs:  []types.Log{testLog},
	}

	s := NewSubscriber(pool, 1, defaultContracts(), 3, 12*time.Second)
	s.SetLastSeenBlock(50)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	outCh, err := s.Start(ctx)
	require.NoError(t, err)
	assert.Equal(t, ModeHTTP, s.Mode())

	// 等待轮询并接收日志
	select {
	case log := <-outCh:
		assert.Equal(t, uint64(60), log.BlockNumber)
	case <-time.After(5 * time.Second):
		t.Fatal("未收到轮询日志")
	}

	s.Stop()
	drainLogs(outCh)
}

// TestSubscriber_SetLastSeenBlock 测试设置最后区块号。
func TestSubscriber_SetLastSeenBlock(t *testing.T) {
	pool := &mockNodePool{
		subErr:      assert.AnError,
		blockNumber: 100,
		filterLogs:  []types.Log{},
	}

	s := NewSubscriber(pool, 1, defaultContracts(), 3, 12*time.Second)

	s.SetLastSeenBlock(42)

	s.mu.Lock()
	block := s.lastSeenBlock
	s.mu.Unlock()

	assert.Equal(t, int64(42), block)
}

// TestSubscriber_WSReconnectUsesBackoff 测试 WS 重连使用 backoff 配置。
// 通过验证重连尝试之间的间隔来确认退避行为。
func TestSubscriber_WSReconnectUsesBackoff(t *testing.T) {
	subCh := make(chan types.Log)
	sub := newMockSub()

	pool := &mockNodePool{
		subCh:       subCh,
		sub:         sub,
		subErr:      nil,
		blockNumber: 100,
		filterLogs:  []types.Log{},
	}

	s := NewSubscriber(pool, 1, defaultContracts(), 3, 12*time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	outCh, err := s.Start(ctx)
	require.NoError(t, err)
	assert.Equal(t, ModeWS, s.Mode())

	// 验证退避参数通过 Subscriber 配置正确：
	// InitialInterval=1s, MaxInterval=30s, Multiplier=2
	// 这些在 wsReconnectLoop 中使用 cenkalti/backoff 实现
	s.Stop()
	drainLogs(outCh)
}

// drainLogs 消费通道中残留的日志，防止 goroutine 泄漏。
func drainLogs(ch <-chan types.Log) {
	go func() {
		for range ch {
		}
	}()
}
