package sync

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	metricNamespace = "sync"
)

var (
	// syncLagBlocks 同步延迟（落后区块数）。
	syncLagBlocks = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: metricNamespace,
			Name:      "lag_blocks",
			Help:      "同步延迟的区块数",
		},
		[]string{"chain_id"},
	)

	// projectorEventsTotal 投影事件计数。
	projectorEventsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: metricNamespace,
			Name:      "projector_events_total",
			Help:      "投影事件处理总数",
		},
		[]string{"chain_id", "event_name", "status"},
	)

	// projectorProcessingSeconds 投影处理耗时。
	projectorProcessingSeconds = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: metricNamespace,
			Name:      "projector_processing_seconds",
			Help:      "投影事件处理耗时",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"chain_id", "event_name"},
	)
)

func init() {
	prometheus.MustRegister(syncLagBlocks)
	prometheus.MustRegister(projectorEventsTotal)
	prometheus.MustRegister(projectorProcessingSeconds)
}

// syncMetrics 记录同步框架指标。
type syncMetrics struct {
	chainID int
}

// newSyncMetrics 创建同步指标实例。
func newSyncMetrics(chainID int) *syncMetrics {
	return &syncMetrics{chainID: chainID}
}

// recordProjection 记录投影事件指标。
func (m *syncMetrics) recordProjection(eventName, status string, elapsed time.Duration) {
	chainIDLabel := formatChainID(m.chainID)

	projectorEventsTotal.WithLabelValues(chainIDLabel, eventName, status).Inc()
	projectorProcessingSeconds.WithLabelValues(chainIDLabel, eventName).Observe(elapsed.Seconds())
}

// recordSyncLag 记录同步延迟指标。
func (m *syncMetrics) recordSyncLag(lagBlocks int64) {
	syncLagBlocks.WithLabelValues(formatChainID(m.chainID)).Set(float64(lagBlocks))
}

// formatChainID 将链 ID 格式化为字符串标签。
func formatChainID(chainID int) string {
	switch chainID {
	case 1:
		return "ethereum"
	case 11155111:
		return "sepolia"
	default:
		return "unknown"
	}
}
