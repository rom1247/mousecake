// Package middleware 提供 Gin HTTP 中间件。
package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"

	"github.com/mousecake-go/mousecake-go/internal/shared/errs"
	"github.com/mousecake-go/mousecake-go/internal/shared/response"
)

const (
	cleanupInterval = 3 * time.Minute
	maxIdleTime     = 3 * time.Minute
)

// NewIPRateLimit 创建按 IP 维度的令牌桶限流中间件。
func NewIPRateLimit(r float64, burst int) gin.HandlerFunc {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		t := time.NewTicker(cleanupInterval)
		defer t.Stop()
		for range t.C {
			mu.Lock()
			now := time.Now()
			for ip, c := range clients {
				if now.Sub(c.lastSeen) > maxIdleTime {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()

		mu.Lock()
		cl, exists := clients[ip]
		if !exists {
			cl = &client{limiter: rate.NewLimiter(rate.Limit(r), burst)}
			clients[ip] = cl
		}
		cl.lastSeen = time.Now()
		mu.Unlock()

		if !cl.limiter.Allow() {
			response.Error(c, http.StatusTooManyRequests, errs.CodeRateLimited, errs.GetErrorMessage(errs.CodeRateLimited))
			c.Abort()
			return
		}

		c.Next()
	}
}

// NewAccountRateLimit 创建按账号维度的令牌桶限流中间件。
// 从 gin.Context 的 "account_key" 中读取账号标识。
func NewAccountRateLimit(r float64, burst int) gin.HandlerFunc {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		t := time.NewTicker(cleanupInterval)
		defer t.Stop()
		for range t.C {
			mu.Lock()
			now := time.Now()
			for key, c := range clients {
				if now.Sub(c.lastSeen) > maxIdleTime {
					delete(clients, key)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		key, exists := c.Get("account_key")
		if !exists {
			c.Next()
			return
		}

		accountKey, ok := key.(string)
		if !ok {
			c.Next()
			return
		}

		mu.Lock()
		cl, exists := clients[accountKey]
		if !exists {
			cl = &client{limiter: rate.NewLimiter(rate.Limit(r), burst)}
			clients[accountKey] = cl
		}
		cl.lastSeen = time.Now()
		mu.Unlock()

		if !cl.limiter.Allow() {
			response.Error(c, http.StatusTooManyRequests, errs.CodeRateLimited, errs.GetErrorMessage(errs.CodeRateLimited))
			c.Abort()
			return
		}

		c.Next()
	}
}

// addressRequest 从请求体中提取 address 字段的最小结构。
type addressRequest struct {
	Address string `json:"address"`
}

// extractAddress 从请求中提取钱包地址，按优先级依次检查：
// 1. X-Address 请求头
// 2. address 查询参数
// 3. POST 请求体中的 address JSON 字段
// 从 body 读取后会将 body 缓存回 c.Request.Body，确保后续 handler 可再次读取。
func extractAddress(c *gin.Context) string {
	// 优先从请求头获取
	if addr := strings.TrimSpace(c.GetHeader("X-Address")); addr != "" {
		return addr
	}

	// 其次从查询参数获取
	if addr := strings.TrimSpace(c.Query("address")); addr != "" {
		return addr
	}

	// 最后从 POST body 获取
	if c.Request.Method != http.MethodPost || c.Request.Body == nil {
		return ""
	}

	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return ""
	}
	// 恢复 body 供后续 handler 读取
	c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	var req addressRequest
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		return ""
	}
	return strings.TrimSpace(req.Address)
}

// NewAddressRateLimit 创建按钱包地址维度的令牌桶限流中间件。
// 从请求头（X-Address）、查询参数（address）或 POST 请求体中提取地址作为限流键。
// 如果无法提取地址，跳过限流。
func NewAddressRateLimit(r float64, burst int) gin.HandlerFunc {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		t := time.NewTicker(cleanupInterval)
		defer t.Stop()
		for range t.C {
			mu.Lock()
			now := time.Now()
			for addr, c := range clients {
				if now.Sub(c.lastSeen) > maxIdleTime {
					delete(clients, addr)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		address := extractAddress(c)
		if address == "" {
			c.Next()
			return
		}
		// 地址标准化为小写，避免大小写差异导致绕过
		address = strings.ToLower(address)

		mu.Lock()
		cl, exists := clients[address]
		if !exists {
			cl = &client{limiter: rate.NewLimiter(rate.Limit(r), burst)}
			clients[address] = cl
		}
		cl.lastSeen = time.Now()
		mu.Unlock()

		if !cl.limiter.Allow() {
			response.Error(c, http.StatusTooManyRequests, errs.CodeRateLimited, errs.GetErrorMessage(errs.CodeRateLimited))
			c.Abort()
			return
		}

		c.Next()
	}
}
