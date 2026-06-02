package quote

import (
	"context"
	"testing"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/mousecake-go/mousecake-go/internal/quote/domain"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	// 创建表结构
	sql := `
	CREATE TABLE swap_records (
		id BIGINT PRIMARY KEY,
		provider TEXT NOT NULL,
		chain_id INTEGER NOT NULL,
		from_token TEXT NOT NULL,
		to_token TEXT NOT NULL,
		from_amount TEXT NOT NULL,
		to_amount TEXT NOT NULL,
		slippage_percent REAL NOT NULL DEFAULT 0,
		swap_mode TEXT NOT NULL,
		status TEXT NOT NULL DEFAULT 'pending',
		tx_hash TEXT,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`
	require.NoError(t, db.Exec(sql).Error)
	return db
}

func TestSwapRecordRepository_Save(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSwapRecordRepository(db)
	ctx := context.Background()

	record, err := domain.NewSwapRecord(1, domain.NewSwapRecordOpts{
		Provider:        "okx",
		ChainID:         1,
		FromToken:       "0xA",
		ToToken:         "0xB",
		FromAmount:      "1000",
		ToAmount:        "2000",
		SlippagePercent: 0.5,
		SwapMode:        domain.SwapModeExactIn,
	})
	require.NoError(t, err)

	err = repo.Save(ctx, record)
	assert.NoError(t, err)

	// 验证可以查到
	found, err := repo.FindByID(ctx, record.ID)
	require.NoError(t, err)
	assert.Equal(t, "okx", found.Provider)
	assert.Equal(t, domain.SwapStatusPending, found.Status)
}

func TestSwapRecordRepository_FindByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSwapRecordRepository(db)
	ctx := context.Background()

	t.Run("查找存在的记录", func(t *testing.T) {
		record := createTestRecord(t)
		require.NoError(t, repo.Save(ctx, record))

		found, err := repo.FindByID(ctx, record.ID)
		require.NoError(t, err)
		assert.Equal(t, record.ID, found.ID)
		assert.Equal(t, record.Provider, found.Provider)
	})

	t.Run("查找不存在的 ID", func(t *testing.T) {
		_, err := repo.FindByID(ctx, 999)
		assert.ErrorIs(t, err, domain.ErrSwapRecordNotFound)
	})
}

func TestSwapRecordRepository_UpdateStatus(t *testing.T) {
	db := setupTestDB(t)
	repo := NewSwapRecordRepository(db)
	ctx := context.Background()

	t.Run("首次提交成功", func(t *testing.T) {
		record := createTestRecord(t)
		require.NoError(t, repo.Save(ctx, record))

		err := repo.UpdateStatus(ctx, record.ID, "0xabc123")
		assert.NoError(t, err)

		found, err := repo.FindByID(ctx, record.ID)
		require.NoError(t, err)
		assert.Equal(t, domain.SwapStatusSubmitted, found.Status)
		assert.Equal(t, "0xabc123", found.TxHash)
	})

	t.Run("重复提交被拒绝", func(t *testing.T) {
		record := createTestRecord(t)
		require.NoError(t, repo.Save(ctx, record))

		// 首次提交
		require.NoError(t, repo.UpdateStatus(ctx, record.ID, "0xabc"))

		// 再次提交
		err := repo.UpdateStatus(ctx, record.ID, "0xdef")
		assert.ErrorIs(t, err, domain.ErrAlreadySubmitted)
	})
}

func createTestRecord(t *testing.T) *domain.SwapRecord {
	t.Helper()
	record, err := domain.NewSwapRecord(1, domain.NewSwapRecordOpts{
		Provider:        "okx",
		ChainID:         1,
		FromToken:       "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
		ToToken:         "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
		FromAmount:      "1000000",
		ToAmount:        "2000000",
		SlippagePercent: 0.5,
		SwapMode:        domain.SwapModeExactIn,
	})
	require.NoError(t, err)
	return record
}

func TestSwapRecordPOConversion(t *testing.T) {
	record := domain.ReconstructSwapRecord(domain.SwapRecordSnapshot{
		ID:              17369650207862784,
		Provider:        "okx",
		ChainID:         1,
		FromToken:       "0xA",
		ToToken:         "0xB",
		FromAmount:      "1000",
		ToAmount:        "2000",
		SlippagePercent: 0.5,
		SwapMode:        domain.SwapModeExactIn,
		Status:          domain.SwapStatusPending,
		TxHash:          "",
	})

	po := toSwapRecordPO(record)
	assert.Equal(t, int64(17369650207862784), po.ID)
	assert.Equal(t, "okx", po.Provider)
	assert.Equal(t, "pending", po.Status)

	entity := toSwapRecordEntity(&po)
	assert.Equal(t, record.ID, entity.ID)
	assert.Equal(t, record.Provider, entity.Provider)
}
