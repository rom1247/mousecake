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
)

// mockSub 实现 ethereum.Subscription，用于测试。
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
	if m == nil {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if !m.closed {
		m.closed = true
		close(m.unsubDone)
	}
}

// testNodePool mocks NodePool for Subscriber tests.
type testNodePool struct {
	mu sync.Mutex

	blockNumber uint64
	blockErr    error

	filterLogs []types.Log
	filterErr  error

	subCh    chan types.Log
	sub      ethereum.Subscription
	subErr   error
	subCalls int
}

func (m *testNodePool) CallContract(_ context.Context, _ ethereum.CallMsg, _ *big.Int) ([]byte, error) {
	return nil, nil
}

func (m *testNodePool) FilterLogs(_ context.Context, _ ethereum.FilterQuery) ([]types.Log, error) {
	return m.filterLogs, m.filterErr
}

func (m *testNodePool) TransactionReceipt(_ context.Context, _ common.Hash) (*types.Receipt, error) {
	return nil, nil
}

func (m *testNodePool) BlockNumber(_ context.Context) (uint64, error) {
	return m.blockNumber, m.blockErr
}

func (m *testNodePool) HeaderByNumber(_ context.Context, _ *big.Int) (*types.Header, error) {
	return nil, nil
}

func (m *testNodePool) SubscribeLogs(_ context.Context, _ ethereum.FilterQuery) (chan types.Log, ethereum.Subscription, error) {
	m.mu.Lock()
	m.subCalls++
	m.mu.Unlock()

	return m.subCh, m.sub, m.subErr
}

func (m *testNodePool) BatchCallContract(_ context.Context, _ []ethereum.CallMsg) ([][]byte, error) {
	return nil, nil
}

func (m *testNodePool) PendingNonceAt(_ context.Context, _ common.Address) (uint64, error) {
	return 0, nil
}

func (m *testNodePool) EstimateGas(_ context.Context, _ ethereum.CallMsg) (uint64, error) {
	return 0, nil
}

func (m *testNodePool) SendTransaction(_ context.Context, _ *types.Transaction) error {
	return nil
}

func (m *testNodePool) SuggestGasPrice(_ context.Context) (*big.Int, error) {
	return big.NewInt(0), nil
}

func (m *testNodePool) Close() {}

// testAddresses 返回测试用合约地址列表。
var testAddresses = []common.Address{
	common.HexToAddress("0x1234567890123456789012345678901234567890"),
}

// TestSubscriber_StartReturnsLogChannel 测试 Start 返回可用的日志通道。
func TestSubscriber_StartReturnsLogChannel(t *testing.T) {
	subCh := make(chan types.Log)
	mock := newMockSub()

	pool := &testNodePool{
		subCh: subCh,
		sub:   mock,
	}

	s := NewSubscriber(pool, 1, testAddresses)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	outCh, err := s.Start(ctx)
	require.NoError(t, err)
	require.NotNil(t, outCh, "Start 应返回非 nil 日志通道")

	s.Stop()
	drainLogs(outCh)
}

// TestSubscriber_LogsForwarded 测试 WS 推送的日志正确转发到输出通道。
func TestSubscriber_LogsForwarded(t *testing.T) {
	subCh := make(chan types.Log, 1)
	mock := newMockSub()

	pool := &testNodePool{
		subCh: subCh,
		sub:   mock,
	}

	s := NewSubscriber(pool, 1, testAddresses)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	outCh, err := s.Start(ctx)
	require.NoError(t, err)

	// 发送测试日志到源通道
	testLog := types.Log{
		Address:     testAddresses[0],
		Topics:      []common.Hash{common.HexToHash("0xabc")},
		BlockNumber: 100,
		TxHash:      common.HexToHash("0xdef"),
		Index:       0,
	}
	subCh <- testLog

	// 验证输出通道收到日志
	select {
	case log := <-outCh:
		assert.Equal(t, uint64(100), log.BlockNumber)
		assert.Equal(t, testAddresses[0], log.Address)
	case <-time.After(2 * time.Second):
		t.Fatal("未收到转发的日志")
	}

	s.Stop()
	drainLogs(outCh)
}

// TestSubscriber_AutoReconnectOnSubError 测试 sub.Err() 触发自动重连。
func TestSubscriber_AutoReconnectOnSubError(t *testing.T) {
	// 第一次订阅成功，收到错误后重连
	firstSubCh := make(chan types.Log, 1)
	firstSub := newMockSub()

	// 使用 atomicallySwappablePool 支持动态切换订阅
	pool := &swappableSubPool{
		subCh: firstSubCh,
		sub:   firstSub,
	}

	s := NewSubscriber(pool, 1, testAddresses)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	outCh, err := s.Start(ctx)
	require.NoError(t, err)

	// 在第二次订阅之前先准备重连后的订阅
	secondSubCh := make(chan types.Log, 1)
	secondSub := newMockSub()

	// 切换到第二次订阅
	pool.swap(secondSubCh, secondSub)

	// 发送错误触发第一次 bridgeSubscription 返回错误
	firstSub.errCh <- assert.AnError

	// 给 Resubscribe 时间重连
	time.Sleep(500 * time.Millisecond)

	// 在第二次订阅通道发送日志验证重连成功
	testLog := types.Log{
		Address:     testAddresses[0],
		BlockNumber: 200,
		TxHash:      common.HexToHash("0x123"),
	}
	secondSubCh <- testLog

	select {
	case log := <-outCh:
		assert.Equal(t, uint64(200), log.BlockNumber)
	case <-time.After(3 * time.Second):
		t.Fatal("重连后未收到日志")
	}

	s.Stop()
	drainLogs(outCh)
}

// TestSubscriber_AutoReconnectOnChannelClose 测试源通道关闭触发自动重连。
func TestSubscriber_AutoReconnectOnChannelClose(t *testing.T) {
	firstSubCh := make(chan types.Log, 1)
	firstSub := newMockSub()

	secondSubCh := make(chan types.Log, 1)
	secondSub := newMockSub()

	pool := &swappableSubPool{
		subCh: firstSubCh,
		sub:   firstSub,
	}

	s := NewSubscriber(pool, 1, testAddresses)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	outCh, err := s.Start(ctx)
	require.NoError(t, err)

	// 先切换到第二次订阅，然后关闭第一次通道
	pool.swap(secondSubCh, secondSub)
	close(firstSubCh)

	// 给 Resubscribe 时间重连
	time.Sleep(500 * time.Millisecond)

	// 在第二次订阅通道发送日志验证重连
	testLog := types.Log{
		Address:     testAddresses[0],
		BlockNumber: 300,
	}
	secondSubCh <- testLog

	select {
	case log := <-outCh:
		assert.Equal(t, uint64(300), log.BlockNumber)
	case <-time.After(3 * time.Second):
		t.Fatal("源通道关闭重连后未收到日志")
	}

	s.Stop()
	drainLogs(outCh)
}

// TestSubscriber_StopClosesScope 测试 Stop 正确关闭订阅。
func TestSubscriber_StopClosesScope(t *testing.T) {
	subCh := make(chan types.Log)
	mock := newMockSub()

	pool := &testNodePool{
		subCh: subCh,
		sub:   mock,
	}

	s := NewSubscriber(pool, 1, testAddresses)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	outCh, err := s.Start(ctx)
	require.NoError(t, err)

	s.Stop()

	// 验证输出通道关闭
	select {
	case _, ok := <-outCh:
		assert.False(t, ok, "Stop 后输出通道应关闭")
	case <-time.After(2 * time.Second):
		t.Fatal("Stop 后输出通道未在预期时间内关闭")
	}
}

// swappableSubPool 支持动态切换订阅的 mock pool，用于重连测试。
type swappableSubPool struct {
	mu    sync.Mutex
	subCh chan types.Log
	sub   ethereum.Subscription
}

func (p *swappableSubPool) swap(subCh chan types.Log, sub ethereum.Subscription) {
	p.mu.Lock()
	p.subCh = subCh
	p.sub = sub
	p.mu.Unlock()
}

func (p *swappableSubPool) SubscribeLogs(_ context.Context, _ ethereum.FilterQuery) (chan types.Log, ethereum.Subscription, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.subCh, p.sub, nil
}

func (p *swappableSubPool) CallContract(_ context.Context, _ ethereum.CallMsg, _ *big.Int) ([]byte, error) {
	return nil, nil
}
func (p *swappableSubPool) FilterLogs(_ context.Context, _ ethereum.FilterQuery) ([]types.Log, error) {
	return nil, nil
}
func (p *swappableSubPool) TransactionReceipt(_ context.Context, _ common.Hash) (*types.Receipt, error) {
	return nil, nil
}
func (p *swappableSubPool) BlockNumber(_ context.Context) (uint64, error) { return 0, nil }
func (p *swappableSubPool) HeaderByNumber(_ context.Context, _ *big.Int) (*types.Header, error) {
	return nil, nil
}
func (p *swappableSubPool) BatchCallContract(_ context.Context, _ []ethereum.CallMsg) ([][]byte, error) {
	return nil, nil
}
func (p *swappableSubPool) PendingNonceAt(_ context.Context, _ common.Address) (uint64, error) {
	return 0, nil
}
func (p *swappableSubPool) EstimateGas(_ context.Context, _ ethereum.CallMsg) (uint64, error) {
	return 0, nil
}
func (p *swappableSubPool) SendTransaction(_ context.Context, _ *types.Transaction) error {
	return nil
}
func (p *swappableSubPool) SuggestGasPrice(_ context.Context) (*big.Int, error) {
	return big.NewInt(0), nil
}
func (p *swappableSubPool) Close() {}

// drainLogs 消费通道中残留的日志，防止 goroutine 泄漏。
func drainLogs(ch <-chan types.Log) {
	go func() {
		for range ch {
		}
	}()
}
