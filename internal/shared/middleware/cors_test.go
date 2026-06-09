package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/mousecake-go/mousecake-go/config"
)

func TestCORS_AllowedOrigin(t *testing.T) {
	cfg := config.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	r := gin.New()
	r.Use(NewCORS(cfg))
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	tests := []struct {
		name          string
		origin        string
		expectOrigin  string
		expectVary    bool
		method        string
		expectStatus  int
		expectMethods string
	}{
		{
			name:         "白名单 Origin GET 请求",
			origin:       "http://localhost:3000",
			expectOrigin: "http://localhost:3000",
			expectVary:   true,
			method:       "GET",
			expectStatus: 200,
		},
		{
			name:          "白名单 Origin OPTIONS 预检",
			origin:        "http://localhost:5173",
			expectOrigin:  "http://localhost:5173",
			expectVary:    true,
			method:        "OPTIONS",
			expectStatus:  204,
			expectMethods: "GET,POST,OPTIONS",
		},
		{
			name:         "非白名单 Origin 被拒绝",
			origin:       "http://evil.example.com",
			expectOrigin: "",
			expectVary:   false,
			method:       "GET",
			expectStatus: 403,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(tt.method, "/test", nil)
			req.Header.Set("Origin", tt.origin)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectStatus, w.Code)

			if tt.expectOrigin != "" {
				assert.Equal(t, tt.expectOrigin, w.Header().Get("Access-Control-Allow-Origin"))
			} else {
				assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
			}

			if tt.expectVary {
				assert.Contains(t, w.Header().Get("Vary"), "Origin")
			}

			if tt.expectMethods != "" {
				assert.Equal(t, tt.expectMethods, w.Header().Get("Access-Control-Allow-Methods"))
			}
		})
	}
}

func TestCORS_AllowCredentials(t *testing.T) {
	cfg := config.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	r := gin.New()
	r.Use(NewCORS(cfg))
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	r.ServeHTTP(w, req)

	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
}

func TestCORS_WildcardWarning(t *testing.T) {
	cfg := config.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET"},
		AllowHeaders: []string{"Content-Type"},
		MaxAge:       12 * time.Hour,
	}

	r := gin.New()
	r.Use(NewCORS(cfg))
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestCORS_MaxAgeOnPreflight(t *testing.T) {
	cfg := config.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{"GET", "POST"},
		AllowHeaders: []string{"Content-Type"},
		MaxAge:       6 * time.Hour,
	}

	r := gin.New()
	r.Use(NewCORS(cfg))
	r.GET("/test", func(c *gin.Context) { c.Status(200) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	r.ServeHTTP(w, req)

	assert.Equal(t, "21600", w.Header().Get("Access-Control-Max-Age"))
}
