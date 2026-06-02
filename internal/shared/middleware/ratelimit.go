// Package middleware 提供 Gin HTTP 中间件。
package middleware

import (
	"net/http"
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
