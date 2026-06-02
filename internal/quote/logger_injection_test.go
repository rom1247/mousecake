package quote

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestLogger() (*slog.Logger, *bytes.Buffer) {
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)
	return logger, &buf
}

func TestNewQuoteService_LoggerInjection(t *testing.T) {
	logger, buf := setupTestLogger()
	defer func() { slog.SetDefault(slog.Default()) }()

	registry := NewProviderRegistry()
	cache := NewMemoryCache(0)
	repo := NewSwapRecordRepository(setupTestDB(t))
	svc := NewQuoteService(registry, cache, repo, 1)

	assert.NotNil(t, svc.log)
	svc.log.Info("测试模块日志")

	var entry map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	assert.Equal(t, "quote", entry["module"])
	assert.Equal(t, "service", entry["layer"])
	_ = logger
}

func TestNewHandler_LoggerInjection(t *testing.T) {
	logger, buf := setupTestLogger()
	defer func() { slog.SetDefault(slog.Default()) }()

	registry := NewProviderRegistry()
	cache := NewMemoryCache(0)
	repo := NewSwapRecordRepository(setupTestDB(t))
	svc := NewQuoteService(registry, cache, repo, 1)
	handler := NewHandler(svc)

	assert.NotNil(t, handler.log)
	handler.log.Info("测试模块日志")

	var entry map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	assert.Equal(t, "quote", entry["module"])
	assert.Equal(t, "handler", entry["layer"])
	_ = logger
}

func TestNewSwapRecordRepository_LoggerInjection(t *testing.T) {
	logger, buf := setupTestLogger()
	defer func() { slog.SetDefault(slog.Default()) }()

	repo := NewSwapRecordRepository(setupTestDB(t))

	assert.NotNil(t, repo.log)
	repo.log.Info("测试模块日志")

	var entry map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	assert.Equal(t, "quote", entry["module"])
	assert.Equal(t, "repository", entry["layer"])
	_ = logger
}
