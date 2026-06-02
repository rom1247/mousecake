package chain

import (
	"context"
	"log/slog"
	"time"
)

// healthChecker 定期检查节点健康状态。
type healthChecker struct {
	nodes   []*node
	metrics *rpcMetrics
}

// newHealthChecker 创建健康检查器。
func newHealthChecker(nodes []*node, metrics *rpcMetrics) *healthChecker {
	return &healthChecker{
		nodes:   nodes,
		metrics: metrics,
	}
}

// start 启动后台健康检查 goroutine。
func (h *healthChecker) start(ctx context.Context) {
	go h.run(ctx)
}

// run 定期执行健康检查。
func (h *healthChecker) run(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			h.checkAll(ctx)
		}
	}
}

// checkAll 检查所有节点。
func (h *healthChecker) checkAll(ctx context.Context) {
	for _, n := range h.nodes {
		h.checkNode(ctx, n)
	}
}

// checkNode 检查单个节点。
func (h *healthChecker) checkNode(ctx context.Context, n *node) {
	checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := n.blockNumber(checkCtx)
	success := err == nil

	n.recordHealthResult(success)

	n.mu.Lock()
	wasHealthy := n.healthy
	failCount := n.failCount
	n.mu.Unlock()

	if wasHealthy && failCount >= 3 {
		n.setHealthy(false)
		if h.metrics != nil {
			h.metrics.setNodeState(n.name, false)
		}
		slog.Warn("节点标记为不健康", "node", n.name, "fail_count", failCount)
	} else if !wasHealthy && failCount == 0 {
		n.setHealthy(true)
		if h.metrics != nil {
			h.metrics.setNodeState(n.name, true)
		}
		slog.Info("节点恢复健康", "node", n.name)
	}
}
