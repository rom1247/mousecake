package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// skipPaths 不记录 AccessLog 的路径集合。
var skipPaths = map[string]bool{
	"/healthz": true,
	"/readyz":  true,
	"/metrics": true,
}

// statusWriter 包装 gin.ResponseWriter，拦截 WriteHeader 捕获 status code。
type statusWriter struct {
	gin.ResponseWriter
	status int
}

// WriteHeader 拦截状态码写入。
func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// AccessLog 是 Gin 中间件，自动记录 HTTP 请求处理结果。
func AccessLog(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if skipPaths[c.Request.URL.Path] {
			c.Next()
			return
		}

		start := time.Now()
		w := &statusWriter{ResponseWriter: c.Writer, status: 200}
		c.Writer = w

		c.Next()

		duration := time.Since(start)
		status := w.status

		attrs := []any{
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", status),
			slog.Int64("duration_ms", duration.Milliseconds()),
			slog.String("client_ip", c.ClientIP()),
		}

		if rid := GetRequestID(c); rid != "" {
			attrs = append(attrs, slog.String("request_id", rid))
		}

		msg := "HTTP 请求完成"
		switch {
		case status >= 500:
			logger.ErrorContext(c.Request.Context(), msg, attrs...)
		case status >= 400:
			logger.WarnContext(c.Request.Context(), msg, attrs...)
		default:
			logger.InfoContext(c.Request.Context(), msg, attrs...)
		}
	}
}
