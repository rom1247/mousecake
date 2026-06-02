package sync

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/mousecake-go/mousecake-go/config"
)

// EventService 定义事件处理接口，由业务模块实现。
type EventService interface {
	// HandleEvent 处理单条链上事件。
	HandleEvent(ctx context.Context, event *ChainEvent) error
}

// Projector 异步消费 pending 事件，调用 EventService 写入业务投影表。
type Projector struct {
	store   *EventStore
	svc     EventService
	chainID int
	cfg     config.ProjectorConfig
	metrics *syncMetrics

	// 区块级分区：按区块号分到不同 worker
	partitions map[int64]chan *ChainEvent
	mu         sync.Mutex
	stopCh     chan struct{}
	wg         sync.WaitGroup
}

// NewProjector 创建 Projector 实例。
func NewProjector(
	store *EventStore,
	svc EventService,
	chainID int,
	cfg config.ProjectorConfig,
	metrics *syncMetrics,
) *Projector {
	return &Projector{
		store:      store,
		svc:        svc,
		chainID:    chainID,
		cfg:        cfg,
		metrics:    metrics,
		partitions: make(map[int64]chan *ChainEvent),
		stopCh:     make(chan struct{}),
	}
}

// Start 启动 Projector。
func (p *Projector) Start(ctx context.Context) {
	p.wg.Add(1)
	go p.pollLoop(ctx)

	p.wg.Add(1)
	go p.timeoutScanLoop(ctx)
}

// Stop 停止 Projector，等待在途事件处理完成。
func (p *Projector) Stop() {
	close(p.stopCh)
	p.wg.Wait()

	p.mu.Lock()
	for _, ch := range p.partitions {
		close(ch)
	}
	p.partitions = make(map[int64]chan *ChainEvent)
	p.mu.Unlock()
}

// pollLoop 轮询 pending 事件并分发。
func (p *Projector) pollLoop(ctx context.Context) {
	defer p.wg.Done()

	ticker := time.NewTicker(p.cfg.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-p.stopCh:
			return
		case <-ticker.C:
			p.pollOnce(ctx)
		}
	}
}

// pollOnce 执行一次轮询。
func (p *Projector) pollOnce(ctx context.Context) {
	events, err := p.store.ListPending(ctx, p.chainID, "", p.cfg.MaxWorkers*10)
	if err != nil {
		slog.Warn("Projector 查询 pending 事件失败", "chain_id", p.chainID, "error", err)
		return
	}

	for i := range events {
		ev := &events[i]

		// 行级乐观锁获取
		claimed, err := p.store.ClaimPending(ctx, ev.ID, "")
		if err != nil {
			slog.Warn("Projector 获取事件失败", "event_id", ev.ID, "error", err)
			continue
		}
		if !claimed {
			continue
		}

		p.dispatch(ctx, ev)
	}
}

// dispatch 将事件分发到对应区块的分区。
func (p *Projector) dispatch(ctx context.Context, event *ChainEvent) {
	blockNum := event.BlockNumber

	p.mu.Lock()
	ch, ok := p.partitions[blockNum]
	if !ok {
		ch = make(chan *ChainEvent, 64)
		p.partitions[blockNum] = ch
		go p.partitionWorker(ctx, blockNum, ch)
	}
	p.mu.Unlock()

	select {
	case ch <- event:
	case <-ctx.Done():
	}
}

// partitionWorker 区块级分区 worker，同区块内串行处理。
func (p *Projector) partitionWorker(ctx context.Context, blockNum int64, ch chan *ChainEvent) {
	for event := range ch {
		p.processEvent(ctx, event)
	}

	// 分区消费完后清理
	p.mu.Lock()
	delete(p.partitions, blockNum)
	p.mu.Unlock()
}

// processEvent 处理单条事件。
func (p *Projector) processEvent(ctx context.Context, event *ChainEvent) {
	start := time.Now()

	err := p.svc.HandleEvent(ctx, event)
	elapsed := time.Since(start)

	if err != nil {
		slog.Warn("事件处理失败",
			"event_id", event.ID,
			"event_name", event.EventName,
			"retry_count", event.RetryCount+1,
			"error", err,
			"elapsed", elapsed)

		if p.metrics != nil {
			p.metrics.recordProjection(event.EventName, "failed", elapsed)
		}

		// 判断是否需要进入死信队列
		if event.RetryCount+1 >= p.cfg.MaxRetries {
			if markErr := p.store.MarkDeadLetter(ctx, event.ID, err.Error()); markErr != nil {
				slog.Error("标记死信失败", "event_id", event.ID, "error", markErr)
			}
			return
		}

		if markErr := p.store.MarkFailed(ctx, event.ID, err.Error()); markErr != nil {
			slog.Error("标记失败", "event_id", event.ID, "error", markErr)
		}
		return
	}

	if markErr := p.store.MarkProcessed(ctx, event.ID); markErr != nil {
		slog.Error("标记已处理失败", "event_id", event.ID, "error", markErr)
	}

	if p.metrics != nil {
		p.metrics.recordProjection(event.EventName, "processed", elapsed)
	}
}

// timeoutScanLoop 定期扫描超时的 processing 事件。
func (p *Projector) timeoutScanLoop(ctx context.Context) {
	defer p.wg.Done()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-p.stopCh:
			return
		case <-ticker.C:
			reset, err := p.store.ResetProcessingTimeout(ctx, p.cfg.ProcessingTimeout)
			if err != nil {
				slog.Warn("重置超时事件失败", "chain_id", p.chainID, "error", err)
				continue
			}
			if reset > 0 {
				slog.Info("重置超时 processing 事件", "chain_id", p.chainID, "count", reset)
			}
		}
	}
}
