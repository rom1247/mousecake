package chain

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/mousecake-go/mousecake-go/config"
)

// SubscriberMode 表示订阅模式。
type SubscriberMode string

const (
	// ModeWS WebSocket 实时推送模式。
	ModeWS SubscriberMode = "ws"
	// ModeHTTP HTTP 轮询降级模式。
	ModeHTTP SubscriberMode = "http"
)

// Subscriber 链上事件订阅器，支持 WS 推送和 HTTP 降级。
type Subscriber struct {
	pool               NodePool
	chainID            int
	addresses          []common.Address
	blockInterval      time.Duration
	confirmationBlocks int64
	contracts          config.SyncContractsConfig

	mu              sync.Mutex
	mode            SubscriberMode
	lastSeenBlock   int64
	lastMessageTime time.Time
	wsCh            chan types.Log
	wsSub           ethereum.Subscription
	stopCh          chan struct{}
	closeOnce       sync.Once
}

// NewSubscriber 创建事件订阅器。
func NewSubscriber(pool NodePool, chainID int, contracts config.SyncContractsConfig, confirmationBlocks int64, blockInterval time.Duration) *Subscriber {
	var addresses []common.Address
	if contracts.MouseTier != "" {
		addresses = append(addresses, common.HexToAddress(contracts.MouseTier))
	}
	if contracts.MousePadByTier != "" {
		addresses = append(addresses, common.HexToAddress(contracts.MousePadByTier))
	}

	return &Subscriber{
		pool:               pool,
		chainID:            chainID,
		addresses:          addresses,
		blockInterval:      blockInterval,
		confirmationBlocks: confirmationBlocks,
		contracts:          contracts,
		mode:               ModeWS,
		stopCh:             make(chan struct{}),
	}
}

// Start 启动订阅，返回日志通道。
func (s *Subscriber) Start(ctx context.Context) (<-chan types.Log, error) {
	outCh := make(chan types.Log, 256)

	query := ethereum.FilterQuery{
		Addresses: s.addresses,
	}

	// ctx 取消时安全关闭 outCh
	go func() {
		<-ctx.Done()
		s.closeOnce.Do(func() { close(outCh) })
	}()

	// 尝试 WS 连接
	ch, sub, err := s.pool.SubscribeLogs(ctx, query)
	if err != nil {
		slog.Warn("WS 订阅失败，降级为 HTTP 轮询", "chain_id", s.chainID, "error", err)
		s.mu.Lock()
		s.mode = ModeHTTP
		s.mu.Unlock()
		go s.httpPollLoop(ctx, outCh)
		return outCh, nil
	}

	s.mu.Lock()
	s.wsCh = ch
	s.wsSub = sub
	s.mode = ModeWS
	s.lastMessageTime = time.Now()
	s.mu.Unlock()

	go s.wsEventLoop(ctx, outCh)
	go s.heartbeatCheck(ctx)

	return outCh, nil
}

// Stop 停止订阅。
func (s *Subscriber) Stop() {
	close(s.stopCh)
	s.mu.Lock()
	if s.wsSub != nil {
		s.wsSub.Unsubscribe()
	}
	s.mu.Unlock()
}

// wsEventLoop 处理 WS 推送事件。
// 只有第一个 wsEventLoop（由 Start 直接启动的）负责关闭 outCh。
func (s *Subscriber) wsEventLoop(ctx context.Context, outCh chan<- types.Log) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stopCh:
			return
		case log, ok := <-s.wsCh:
			if !ok {
				s.handleWSDisconnect(ctx, outCh)
				return
			}
			s.mu.Lock()
			s.lastMessageTime = time.Now()
			s.mu.Unlock()

			select {
			case outCh <- log:
			case <-ctx.Done():
				return
			}
		case err := <-s.wsSub.Err():
			if err != nil {
				slog.Warn("WS 订阅错误", "chain_id", s.chainID, "error", err)
			}
			s.handleWSDisconnect(ctx, outCh)
			return
		}
	}
}

// handleWSDisconnect 处理 WS 断开，降级为 HTTP 轮询。
func (s *Subscriber) handleWSDisconnect(ctx context.Context, outCh chan<- types.Log) {
	s.mu.Lock()
	s.mode = ModeHTTP
	s.mu.Unlock()

	slog.Warn("WS 断开，降级为 HTTP 轮询", "chain_id", s.chainID)

	go s.httpPollLoop(ctx, outCh)
	go s.wsReconnectLoop(ctx, outCh)
}

// heartbeatCheck 检测 WS 僵死。
func (s *Subscriber) heartbeatCheck(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	heartbeatTimeout := s.blockInterval * 3

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.mu.Lock()
			mode := s.mode
			lastMsg := s.lastMessageTime
			s.mu.Unlock()

			if mode == ModeWS && time.Since(lastMsg) > heartbeatTimeout {
				slog.Warn("WS 僵死检测触发", "chain_id", s.chainID, "last_msg_ago", time.Since(lastMsg))
				s.mu.Lock()
				if s.wsSub != nil {
					s.wsSub.Unsubscribe()
				}
				s.mu.Unlock()
				return
			}
		}
	}
}

// httpPollLoop HTTP 轮询补漏。
func (s *Subscriber) httpPollLoop(ctx context.Context, outCh chan<- types.Log) {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.pollOnce(ctx, outCh)
		}
	}
}

// pollOnce 执行一次 HTTP 轮询。
func (s *Subscriber) pollOnce(ctx context.Context, outCh chan<- types.Log) {
	currentBlock, err := s.pool.BlockNumber(ctx)
	if err != nil {
		slog.Warn("HTTP 轮询获取区块号失败", "chain_id", s.chainID, "error", err)
		return
	}

	s.mu.Lock()
	fromBlock := s.lastSeenBlock + 1
	s.mu.Unlock()

	toBlock := int64(currentBlock) - s.confirmationBlocks
	if toBlock <= fromBlock {
		return
	}

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(fromBlock),
		ToBlock:   big.NewInt(toBlock),
		Addresses: s.addresses,
	}

	logs, err := s.pool.FilterLogs(ctx, query)
	if err != nil {
		slog.Warn("HTTP 轮询 FilterLogs 失败", "chain_id", s.chainID, "error", err)
		return
	}

	for _, log := range logs {
		select {
		case outCh <- log:
		case <-ctx.Done():
			return
		}
	}

	s.mu.Lock()
	s.lastSeenBlock = toBlock
	s.mu.Unlock()
}

// wsReconnectLoop 使用指数退避重连 WS。
func (s *Subscriber) wsReconnectLoop(ctx context.Context, outCh chan<- types.Log) {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = 1 * time.Second
	b.MaxInterval = 30 * time.Second
	b.Multiplier = 2
	b.MaxElapsedTime = 0

	operation := func() error {
		query := ethereum.FilterQuery{
			Addresses: s.addresses,
		}

		ch, sub, err := s.pool.SubscribeLogs(ctx, query)
		if err != nil {
			return fmt.Errorf("WS 重连失败: %w", err)
		}

		s.mu.Lock()
		s.wsCh = ch
		s.wsSub = sub
		s.mode = ModeWS
		s.lastMessageTime = time.Now()
		s.mu.Unlock()

		slog.Info("WS 重连成功", "chain_id", s.chainID)

		go s.wsEventLoop(ctx, outCh)
		go s.heartbeatCheck(ctx)

		return nil
	}

	if err := backoff.Retry(operation, backoff.WithContext(b, ctx)); err != nil {
		slog.Warn("WS 重连最终放弃", "chain_id", s.chainID, "error", err)
	}
}

// SetLastSeenBlock 设置最后看到的区块号。
func (s *Subscriber) SetLastSeenBlock(block int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lastSeenBlock = block
}

// Mode 返回当前订阅模式。
func (s *Subscriber) Mode() SubscriberMode {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.mode
}
