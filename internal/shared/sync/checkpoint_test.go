package sync

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCheckpoint_TableName 测试 Checkpoint 表名映射。
func TestCheckpoint_TableName(t *testing.T) {
	t.Parallel()

	var cp Checkpoint
	assert.Equal(t, "sync_checkpoints", cp.TableName())
}

// TestCheckpoint_Fields 测试 Checkpoint 结构体字段。
func TestCheckpoint_Fields(t *testing.T) {
	t.Parallel()

	cp := Checkpoint{
		ChainID:         1,
		ProcessorID:     "launchpad",
		LastSyncedBlock: 12345,
	}

	assert.Equal(t, int64(0), cp.ID)
	assert.Equal(t, 1, cp.ChainID)
	assert.Equal(t, "launchpad", cp.ProcessorID)
	assert.Equal(t, int64(12345), cp.LastSyncedBlock)
}

// TestNewCheckpointRepository 测试创建 CheckpointRepository 实例。
func TestNewCheckpointRepository(t *testing.T) {
	repo := NewCheckpointRepository(nil)
	assert.NotNil(t, repo)
}
