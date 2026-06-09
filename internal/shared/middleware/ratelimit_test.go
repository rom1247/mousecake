package middleware

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestAddressRateLimit_Pass(t *testing.T) {
	r := gin.New()
	r.Use(NewAddressRateLimit(10, 5))
	r.POST("/test", func(c *gin.Context) { c.Status(200) })

	body := `{"address":"0xabc1234567890abcdef1234567890abcdef12345678"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestAddressRateLimit_Reject(t *testing.T) {
	r := gin.New()
	r.Use(NewAddressRateLimit(1, 1))
	r.POST("/test", func(c *gin.Context) { c.Status(200) })

	body := `{"address":"0xdef1234567890abcdef1234567890abcdef123456"}`
	for i := 0; i < 3; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 429, w.Code)
}

func TestAddressRateLimit_IndependentCount(t *testing.T) {
	r := gin.New()
	r.Use(NewAddressRateLimit(1, 1))
	r.POST("/test", func(c *gin.Context) { c.Status(200) })

	bodyA := `{"address":"0xaaa1111111111111111111111111111111111111"}`
	wA1 := httptest.NewRecorder()
	reqA1, _ := http.NewRequest("POST", "/test", strings.NewReader(bodyA))
	reqA1.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(wA1, reqA1)
	assert.Equal(t, 200, wA1.Code)

	wA2 := httptest.NewRecorder()
	reqA2, _ := http.NewRequest("POST", "/test", strings.NewReader(bodyA))
	reqA2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(wA2, reqA2)
	assert.Equal(t, 429, wA2.Code)

	bodyB := `{"address":"0xbbb2222222222222222222222222222222222222"}`
	wB := httptest.NewRecorder()
	reqB, _ := http.NewRequest("POST", "/test", strings.NewReader(bodyB))
	reqB.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(wB, reqB)
	assert.Equal(t, 200, wB.Code)
}

func TestAddressRateLimit_EmptyAddress(t *testing.T) {
	r := gin.New()
	r.Use(NewAddressRateLimit(1, 1))
	r.POST("/test", func(c *gin.Context) { c.Status(200) })

	// 无 address 字段，应跳过限流
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/test", strings.NewReader(`{"other":"value"}`))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
	}
}

func TestAddressRateLimit_NonPost_NoAddress(t *testing.T) {
	r := gin.New()
	r.Use(NewAddressRateLimit(1, 1))
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	// GET 请求无地址来源，应跳过限流
	for i := 0; i < 5; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
	}
}

func TestAddressRateLimit_FromHeader(t *testing.T) {
	r := gin.New()
	r.Use(NewAddressRateLimit(1, 1))
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	// 从 X-Address 头获取地址，第一个请求通过
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/test", nil)
	req1.Header.Set("X-Address", "0xaaa1111111111111111111111111111111111111")
	r.ServeHTTP(w1, req1)
	assert.Equal(t, 200, w1.Code)

	// 同一地址第二个请求被限流
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/test", nil)
	req2.Header.Set("X-Address", "0xaaa1111111111111111111111111111111111111")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, 429, w2.Code)

	// 不同地址不受影响
	w3 := httptest.NewRecorder()
	req3, _ := http.NewRequest("GET", "/test", nil)
	req3.Header.Set("X-Address", "0xbbb2222222222222222222222222222222222222")
	r.ServeHTTP(w3, req3)
	assert.Equal(t, 200, w3.Code)
}

func TestAddressRateLimit_FromQuery(t *testing.T) {
	r := gin.New()
	r.Use(NewAddressRateLimit(1, 1))
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	// 从查询参数获取地址
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("GET", "/test?address=0xccc3333333333333333333333333333333333333", nil)
	r.ServeHTTP(w1, req1)
	assert.Equal(t, 200, w1.Code)

	// 同一地址第二个请求被限流
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/test?address=0xccc3333333333333333333333333333333333333", nil)
	r.ServeHTTP(w2, req2)
	assert.Equal(t, 429, w2.Code)
}

func TestAddressRateLimit_HeaderPriority(t *testing.T) {
	r := gin.New()
	r.Use(NewAddressRateLimit(1, 1))
	r.POST("/test", func(c *gin.Context) { c.Status(200) })

	// header 和 body 都有地址时，header 优先
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", strings.NewReader(`{"address":"0xbody0000000000000000000000000000000000body"}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Address", "0xheader00000000000000000000000000000000head")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	// 用 body 地址请求，应独立计数（不是同一个地址）
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", "/test", strings.NewReader(`{"address":"0xbody0000000000000000000000000000000000body"}`))
	req2.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w2, req2)
	assert.Equal(t, 200, w2.Code)
}

func TestAddressRateLimit_BodyPreserved(t *testing.T) {
	r := gin.New()
	r.Use(NewAddressRateLimit(10, 5))
	r.POST("/test", func(c *gin.Context) {
		body, _ := io.ReadAll(c.Request.Body)
		c.String(200, string(body))
	})

	body := `{"address":"0xabc1234567890abcdef1234567890abcdef12345678"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/test", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "0xabc123")
}
