package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDKey = "request_id"

// RequestID 是 Gin 中间件，为每个请求生成或复用 X-Request-ID。
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader("X-Request-ID")
		if rid == "" {
			rid = uuid.New().String()
		}

		c.Set(requestIDKey, rid)
		c.Header("X-Request-ID", rid)
		c.Next()
	}
}

// GetRequestID 从 gin.Context 中获取 request_id，未经过中间件时返回空字符串。
func GetRequestID(c *gin.Context) string {
	rid, _ := c.Get(requestIDKey)
	if s, ok := rid.(string); ok {
		return s
	}
	return ""
}
