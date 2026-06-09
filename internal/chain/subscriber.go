package chain

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Subscriber 链上事件订阅器，基于 WS 推送 + event.Resubscribe 自动重连。
type Subscriber struct {
	pool      NodePool
	chainID   int
	addresses []common.Address

	stopCh    chan struct{}
	closeOnce sync.Once
	scope     event.SubscriptionScope
}

// NewSubscriber 创建事件订阅器。
func NewSubscriber(pool NodePool, chainID int, addresses []common.Address) *Subscriber {
	return &Subscriber{
		pool:      pool,
		chainID:   chainID,
		addresses: addresses,
		stopCh:    make(chan struct{}),
	}
}

// Start 启动订阅，返回日志通道。
// 内部使用 event.Resubscribe 实现 WS 断开后的指数退避自动重连。
func (s *Subscriber) Start(ctx context.Context) (<-chan types.Log, error) {
	outCh := make(chan types.Log, 256)

	query := ethereum.FilterQuery{
		Addresses: s.addresses,
	}

	resub := event.ResubscribeErr(2*time.Minute, func(subCtx context.Context, lastErr error) (event.Subscription, error) {
		if lastErr != nil {
			slog.Warn("WS 订阅断开，正在重连",
				"chain_id", s.chainID,
				"last_error", lastErr)
		}

		ch, sub, err := s.pool.SubscribeLogs(subCtx, query)
		if err != nil {
			return nil, fmt.Errorf("WS 订阅失败: %w", err)
		}

		return s.bridgeSubscription(ctx, ch, sub, outCh), nil
	})

	s.scope.Track(resub)

	// ctx 取消或显式停止时关闭 outCh
	go func() {
		select {
		case <-ctx.Done():
		case <-s.stopCh:
		}
		s.closeOnce.Do(func() { close(outCh) })
	}()

	return outCh, nil
}

// bridgeSubscription 将 pool.SubscribeLogs 的结果桥接为 event.Subscription。
// 从源通道读取日志转发到 outCh，当源订阅出错时通知 Resubscribe 触发重连。
func (s *Subscriber) bridgeSubscription(
	ctx context.Context,
	srcCh <-chan types.Log,
	srcSub ethereum.Subscription,
	outCh chan<- types.Log,
) event.Subscription {
	return event.NewSubscription(func(quit <-chan struct{}) error {
		if srcSub != nil {
			defer srcSub.Unsubscribe()
		}

		if srcCh == nil || srcSub == nil {
			return fmt.Errorf("WS 订阅返回 nil 通道或订阅")
		}

		for {
			select {
			case <-quit:
				return nil
			case <-ctx.Done():
				return nil
			case <-s.stopCh:
				return nil
			case err := <-srcSub.Err():
				if err != nil {
					return fmt.Errorf("WS 订阅断开: %w", err)
				}
				return fmt.Errorf("WS 订阅已关闭")
			case log, ok := <-srcCh:
				if !ok {
					return fmt.Errorf("WS 日志通道已关闭")
				}
				slog.InfoContext(ctx, "收到 WS 推送日志",
					"chain_id", s.chainID,
					"block", log.BlockNumber,
					"address", log.Address.Hex(),
					"tx", log.TxHash.Hex(),
					"index", log.Index,
				)

				select {
				case outCh <- log:
				case <-quit:
					return nil
				case <-ctx.Done():
					return nil
				}
			}
		}
	})
}

// Stop 停止订阅。
func (s *Subscriber) Stop() {
	close(s.stopCh)
	s.scope.Close()
}
