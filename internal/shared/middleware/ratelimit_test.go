package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestIPRateLimit_Pass(t *testing.T) {
	r := gin.New()
	r.Use(NewIPRateLimit(20, 2))
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestIPRateLimit_Reject(t *testing.T) {
	r := gin.New()
	r.Use(NewIPRateLimit(1, 1))
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "192.168.1.1:1234"
		r.ServeHTTP(w, req)
	}

	// 最后一次应该被限流
	lastW := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "192.168.1.1:1234"
	r.ServeHTTP(lastW, req)
	assert.Equal(t, 429, lastW.Code)
}

func TestAccountRateLimit_Pass(t *testing.T) {
	r := gin.New()
	r.Use(NewAccountRateLimit(10, 5))
	r.POST("/test", func(c *gin.Context) {
		c.Set("account_key", "0xabc1234567890abcdef1234567890abcdef12345678")
		c.Status(200)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

// W5: 验证限流拒绝时返回错误码 40106

func TestIPRateLimit_Reject_ErrorCode(t *testing.T) {
	r := gin.New()
	r.Use(NewIPRateLimit(1, 1))
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.RemoteAddr = "10.0.0.1:1234"
		r.ServeHTTP(w, req)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.RemoteAddr = "10.0.0.1:1234"
	r.ServeHTTP(w, req)

	assert.Equal(t, 429, w.Code)

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(40106), resp["code"])
}

// W6: 不同账号独立计数

func TestAccountRateLimit_IndependentCount(t *testing.T) {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		if key := c.Query("account"); key != "" {
			c.Set("account_key", key)
		}
		c.Next()
	})
	r.Use(NewAccountRateLimit(1, 1))
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	wA1 := httptest.NewRecorder()
	reqA1, _ := http.NewRequest("GET", "/test?account=account_a", nil)
	r.ServeHTTP(wA1, reqA1)
	assert.Equal(t, 200, wA1.Code)

	wA2 := httptest.NewRecorder()
	reqA2, _ := http.NewRequest("GET", "/test?account=account_a", nil)
	r.ServeHTTP(wA2, reqA2)
	assert.Equal(t, 429, wA2.Code)

	wB := httptest.NewRecorder()
	reqB, _ := http.NewRequest("GET", "/test?account=account_b", nil)
	r.ServeHTTP(wB, reqB)
	assert.Equal(t, 200, wB.Code)
}
