package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit_JSONHandler(t *testing.T) {
	var buf bytes.Buffer
	cfg := LogConfig{Level: "info", Format: "json", AddSource: false}
	Init(cfg, &buf)

	slog.Info("测试消息", "key", "value")

	var entry map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	assert.Equal(t, "测试消息", entry["msg"])
	assert.Equal(t, "INFO", entry["level"])
	assert.Equal(t, "value", entry["key"])
	assert.Contains(t, entry, "time")
}

func TestInit_TextHandler(t *testing.T) {
	var buf bytes.Buffer
	cfg := LogConfig{Level: "debug", Format: "text", AddSource: true}
	Init(cfg, &buf)

	slog.Debug("调试消息")

	output := buf.String()
	assert.Contains(t, output, "调试消息")
	assert.Contains(t, output, "level=DEBUG")
}

func TestInit_LevelFilter(t *testing.T) {
	var buf bytes.Buffer
	cfg := LogConfig{Level: "warn", Format: "json", AddSource: false}
	Init(cfg, &buf)

	slog.Info("不应出现的信息")
	slog.Warn("应出现的警告")

	output := buf.String()
	assert.NotContains(t, output, "不应出现的信息")
	assert.Contains(t, output, "应出现的警告")
}

func TestSetLevel(t *testing.T) {
	var buf bytes.Buffer
	cfg := LogConfig{Level: "info", Format: "json", AddSource: false}
	Init(cfg, &buf)

	slog.Debug("过滤掉的调试消息")
	assert.Empty(t, buf.String())

	SetLevel(slog.LevelDebug)
	slog.Debug("可见的调试消息")

	output := buf.String()
	assert.Contains(t, output, "可见的调试消息")
}

func TestReplaceAttr_BigInt(t *testing.T) {
	var buf bytes.Buffer
	cfg := LogConfig{Level: "info", Format: "json", AddSource: false}
	Init(cfg, &buf)

	balance := big.NewInt(1000000000000000000)
	slog.Info("余额", "balance", balance)

	var entry map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	assert.Equal(t, "1000000000000000000", fmt.Sprintf("%v", entry["balance"]))
}

func TestReplaceAttr_NilStringer(t *testing.T) {
	var buf bytes.Buffer
	cfg := LogConfig{Level: "info", Format: "json", AddSource: false}
	Init(cfg, &buf)

	var s *strings.Builder // nil fmt.Stringer
	slog.Info("测试", "value", s)

	var entry map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	assert.Equal(t, "<nil>", fmt.Sprintf("%v", entry["value"]))
}

func TestReplaceAttr_Error(t *testing.T) {
	var buf bytes.Buffer
	cfg := LogConfig{Level: "info", Format: "json", AddSource: false}
	Init(cfg, &buf)

	err := fmt.Errorf("连接失败: %w", fmt.Errorf("timeout"))
	slog.Error("数据库错误", "error", err)

	var entry map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	errMsg := fmt.Sprintf("%v", entry["error"])
	assert.Contains(t, errMsg, "连接失败")
	assert.Contains(t, errMsg, "timeout")
}
