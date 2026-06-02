package chain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFormatChainID 测试链 ID 格式化为字符串标签。
func TestFormatChainID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		chainID int
		want    string
	}{
		{
			name:    "以太坊主网",
			chainID: 1,
			want:    "ethereum",
		},
		{
			name:    "Sepolia 测试网",
			chainID: 11155111,
			want:    "sepolia",
		},
		{
			name:    "未知链",
			chainID: 999,
			want:    "unknown",
		},
		{
			name:    "零值链 ID",
			chainID: 0,
			want:    "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := formatChainID(tt.chainID)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestNewRPCMetrics 测试创建 RPC 指标实例。
func TestNewRPCMetrics(t *testing.T) {
	t.Parallel()

	m := newRPCMetrics(1)
	assert.NotNil(t, m)
	assert.Equal(t, 1, m.chainID)
}

// TestRPCMetrics_RecordRequest 测试记录 RPC 请求不 panic。
func TestRPCMetrics_RecordRequest(t *testing.T) {
	t.Parallel()

	m := newRPCMetrics(1)

	assert.NotPanics(t, func() {
		m.recordRequest("node1", "eth_blockNumber")
	})

	assert.NotPanics(t, func() {
		m.recordRequest("node2", "eth_call")
	})
}

// TestRPCMetrics_RecordError 测试记录 RPC 错误不 panic。
func TestRPCMetrics_RecordError(t *testing.T) {
	t.Parallel()

	m := newRPCMetrics(11155111)

	assert.NotPanics(t, func() {
		m.recordError("node1", "eth_call")
	})

	assert.NotPanics(t, func() {
		m.recordError("node2", "eth_getLogs")
	})
}

// TestRPCMetrics_SetNodeState 测试更新节点健康状态不 panic。
func TestRPCMetrics_SetNodeState(t *testing.T) {
	t.Parallel()

	m := newRPCMetrics(1)

	assert.NotPanics(t, func() {
		m.setNodeState("node1", true)
	})

	assert.NotPanics(t, func() {
		m.setNodeState("node1", false)
	})

	assert.NotPanics(t, func() {
		m.setNodeState("node2", true)
		m.setNodeState("node2", false)
	})
}
