// Package middleware 提供 Gin HTTP 中间件。
package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/mousecake-go/mousecake-go/config"
)

// NewCORS 根据配置创建 CORS 中间件。
// 生产环境禁止 AllowOrigins 包含 "*"（SEC-21），遇到时打印警告日志。
func NewCORS(cfg config.CORSConfig) gin.HandlerFunc {
	for _, o := range cfg.AllowOrigins {
		if o == "*" {
			slog.Warn("CORS AllowOrigins 包含 '*'，生产环境应使用白名单模式（SEC-21）")
			break
		}
	}

	maxAge := cfg.MaxAge
	if maxAge == 0 {
		maxAge = 12 * time.Hour
	}

	return cors.New(cors.Config{
		AllowOrigins:     cfg.AllowOrigins,
		AllowMethods:     cfg.AllowMethods,
		AllowHeaders:     cfg.AllowHeaders,
		ExposeHeaders:    cfg.ExposeHeaders,
		AllowCredentials: cfg.AllowCredentials,
		MaxAge:           maxAge,
	})
}
