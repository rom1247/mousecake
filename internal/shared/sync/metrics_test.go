package sync

import (
	"testing"
	"time"

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

// TestNewSyncMetrics 测试创建同步指标实例。
func TestNewSyncMetrics(t *testing.T) {
	t.Parallel()

	m := newSyncMetrics(1)
	assert.NotNil(t, m)
	assert.Equal(t, 1, m.chainID)
}

// TestSyncMetrics_RecordProjection 测试记录投影事件指标不 panic。
func TestSyncMetrics_RecordProjection(t *testing.T) {
	t.Parallel()

	m := newSyncMetrics(11155111)

	// 多次调用不应 panic，验证 Prometheus 指标正常记录
	assert.NotPanics(t, func() {
		m.recordProjection("Deposited", "processed", 100*time.Millisecond)
	})

	assert.NotPanics(t, func() {
		m.recordProjection("Deposited", "failed", 50*time.Millisecond)
	})

	assert.NotPanics(t, func() {
		m.recordProjection("Withdrawn", "processed", 200*time.Millisecond)
	})
}

// TestSyncMetrics_RecordSyncLag 测试记录同步延迟指标不 panic。
func TestSyncMetrics_RecordSyncLag(t *testing.T) {
	t.Parallel()

	m := newSyncMetrics(1)

	assert.NotPanics(t, func() {
		m.recordSyncLag(42)
	})

	assert.NotPanics(t, func() {
		m.recordSyncLag(0)
	})

	assert.NotPanics(t, func() {
		m.recordSyncLag(1000000)
	})
}
