package middleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatusWriter_Default200(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	sw := &statusWriter{ResponseWriter: c.Writer, status: 200}

	assert.Equal(t, 200, sw.status)
}

func TestStatusWriter_ExplicitStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	sw := &statusWriter{ResponseWriter: c.Writer, status: 200}

	sw.WriteHeader(http.StatusNotFound)
	assert.Equal(t, 404, sw.status)
}

func TestAccessLog_SuccessRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))

	r := gin.New()
	r.Use(RequestID())
	r.Use(AccessLog(logger))
	r.GET("/api/v1/quote", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/quote?from=ETH&to=USDT", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var entry map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	assert.Equal(t, "INFO", entry["level"])
	assert.Equal(t, "GET", entry["method"])
	assert.Equal(t, "/api/v1/quote", entry["path"])
	assert.Equal(t, float64(200), entry["status"])
	assert.Contains(t, entry, "duration_ms")
	assert.Contains(t, entry, "client_ip")
	assert.Contains(t, entry, "request_id")
}

func TestAccessLog_ClientError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))

	r := gin.New()
	r.Use(AccessLog(logger))
	r.POST("/api/v1/swap", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/swap", nil)
	r.ServeHTTP(w, req)

	var entry map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	assert.Equal(t, "WARN", entry["level"])
	assert.Equal(t, float64(400), entry["status"])
}

func TestAccessLog_ServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))

	r := gin.New()
	r.Use(AccessLog(logger))
	r.GET("/api/v1/quote", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "fail"})
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/quote", nil)
	r.ServeHTTP(w, req)

	var entry map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	assert.Equal(t, "ERROR", entry["level"])
	assert.Equal(t, float64(500), entry["status"])
}

func TestAccessLog_SkipHealthPaths(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))

	r := gin.New()
	r.Use(AccessLog(logger))
	r.GET("/healthz", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	r.GET("/readyz", func(c *gin.Context) { c.String(http.StatusOK, "ok") })
	r.GET("/metrics", func(c *gin.Context) { c.String(http.StatusOK, "ok") })

	for _, path := range []string{"/healthz", "/readyz", "/metrics"} {
		buf.Reset()
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, path, nil)
		r.ServeHTTP(w, req)
		assert.Empty(t, buf.String(), "path %s should not produce access log", path)
	}
}

func TestAccessLog_BusinessPathNotSkipped(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))

	r := gin.New()
	r.Use(AccessLog(logger))
	r.GET("/api/v1/quote", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/quote", nil)
	r.ServeHTTP(w, req)

	assert.NotEmpty(t, buf.String())
}
