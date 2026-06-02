package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	svc := NewJWTService("test-secret", 4*time.Hour, 8*time.Hour)
	token, _ := svc.IssueToken(TokenTypeWallet, 1, false)

	r := gin.New()
	r.Use(NewAuthMiddleware(svc))
	r.GET("/test", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		assert.Equal(t, int64(1), userID)
		c.Status(200)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	svc := NewJWTService("test-secret", 4*time.Hour, 8*time.Hour)

	r := gin.New()
	r.Use(NewAuthMiddleware(svc))
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	svc := NewJWTService("test-secret", 4*time.Hour, 8*time.Hour)

	r := gin.New()
	r.Use(NewAuthMiddleware(svc))
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Basic xxx")
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	svc := NewJWTService("test-secret", -1*time.Hour, -1*time.Hour)
	token, _ := svc.IssueToken(TokenTypeWallet, 1, false)

	// 用新的 svc 做验证（secret 相同）
	validSvc := NewJWTService("test-secret", 4*time.Hour, 8*time.Hour)

	r := gin.New()
	r.Use(NewAuthMiddleware(validSvc))
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestAdminMiddleware_AdminPass(t *testing.T) {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("is_admin", true)
		c.Next()
	})
	r.Use(NewAdminMiddleware())
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestAdminMiddleware_NonAdminReject(t *testing.T) {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("is_admin", false)
		c.Next()
	})
	r.Use(NewAdminMiddleware())
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 403, w.Code)
}
