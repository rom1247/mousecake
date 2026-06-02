package chain

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	chainMetricNamespace = "rpc"
)

var (
	// rpcNodeState RPC 节点健康状态（1=健康, 0=不健康）。
	rpcNodeState = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: chainMetricNamespace,
			Name:      "node_state",
			Help:      "RPC 节点健康状态",
		},
		[]string{"chain_id", "node"},
	)

	// rpcRequestsTotal RPC 请求计数。
	rpcRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: chainMetricNamespace,
			Name:      "requests_total",
			Help:      "RPC 请求总数",
		},
		[]string{"chain_id", "node", "method"},
	)

	// rpcErrorsTotal RPC 错误计数。
	rpcErrorsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: chainMetricNamespace,
			Name:      "errors_total",
			Help:      "RPC 错误总数",
		},
		[]string{"chain_id", "node", "method"},
	)
)

func init() {
	prometheus.MustRegister(rpcNodeState)
	prometheus.MustRegister(rpcRequestsTotal)
	prometheus.MustRegister(rpcErrorsTotal)
}

// rpcMetrics 记录 RPC 调用指标。
type rpcMetrics struct {
	chainID int
}

// newRPCMetrics 创建 RPC 指标实例。
func newRPCMetrics(chainID int) *rpcMetrics {
	return &rpcMetrics{chainID: chainID}
}

// recordRequest 记录 RPC 请求。
func (m *rpcMetrics) recordRequest(nodeName, method string) {
	rpcRequestsTotal.WithLabelValues(formatChainID(m.chainID), nodeName, method).Inc()
}

// recordError 记录 RPC 错误。
func (m *rpcMetrics) recordError(nodeName, method string) {
	rpcErrorsTotal.WithLabelValues(formatChainID(m.chainID), nodeName, method).Inc()
}

// setNodeState 更新节点健康状态。
func (m *rpcMetrics) setNodeState(nodeName string, healthy bool) {
	val := float64(0)
	if healthy {
		val = 1
	}
	rpcNodeState.WithLabelValues(formatChainID(m.chainID), nodeName).Set(val)
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
