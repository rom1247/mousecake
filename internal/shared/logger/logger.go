// Package logger 提供全局 slog Logger 初始化和动态级别控制。
package logger

import (
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"os"

	"github.com/holiman/uint256"
)

// LogConfig 日志配置，与 config.LogConfig 保持字段一致。
// 此处独立定义避免循环导入（logger 包不依赖 config 包）。
type LogConfig struct {
	Level     string
	Format    string
	AddSource bool
}

// levelVar 全局动态日志级别变量。
var levelVar = new(slog.LevelVar)

// Init 根据配置初始化全局 slog Logger 并设为默认。
// 可选传入自定义输出目标（用于测试），为 nil 时使用 os.Stdout(os.Stderr for text)。
func Init(cfg LogConfig, w ...io.Writer) {
	levelVar.Set(parseLevel(cfg.Level))

	opts := &slog.HandlerOptions{
		Level:       levelVar,
		AddSource:   cfg.AddSource,
		ReplaceAttr: replaceAttr,
	}

	var writer io.Writer = os.Stdout
	if len(w) > 0 && w[0] != nil {
		writer = w[0]
	}

	var handler slog.Handler
	switch cfg.Format {
	case "text":
		if len(w) == 0 || w[0] == nil {
			writer = os.Stderr
		}
		handler = slog.NewTextHandler(writer, opts)
	default:
		handler = slog.NewJSONHandler(writer, opts)
	}

	slog.SetDefault(slog.New(handler))
}

// SetLevel 动态调整全局日志级别（线程安全）。
func SetLevel(level slog.Level) {
	levelVar.Set(level)
}

// parseLevel 将字符串级别转换为 slog.Level。
func parseLevel(s string) slog.Level {
	switch s {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// replaceAttr 为区块链类型提供自定义格式化。
func replaceAttr(groups []string, a slog.Attr) slog.Attr {
	// 处理 big.Int 和 uint256.Int
	if val := a.Value.Any(); val != nil {
		switch v := val.(type) {
		case *big.Int:
			return slog.String(a.Key, v.String())
		case *uint256.Int:
			return slog.String(a.Key, v.String())
		case fmt.Stringer:
			// nil fmt.Stringer 指针由下面的 nil 检查处理
			return a
		}
	}

	// 处理 nil fmt.Stringer：当 Any() 返回 nil 但 Kind 是 LogValuer 等特殊类型
	if a.Value.Kind() == slog.KindAny {
		if a.Value.Any() == nil {
			return slog.String(a.Key, "<nil>")
		}
	}

	return a
}
