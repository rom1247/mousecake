package quote

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mousecake-go/mousecake-go/internal/quote/domain"
)

func setupTestRouter(svc *QuoteService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	handler := NewHandler(svc)
	v1 := r.Group("/api/v1")
	handler.RegisterRoutes(v1)
	return r
}

func TestValidateEVMAddress(t *testing.T) {
	tests := []struct {
		name    string
		addr    string
		isValid bool
	}{
		{"合法地址", "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48", true},
		{"原生代币地址", "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE", true},
		{"无效地址-太短", "0x1234", false},
		{"无效地址-无0x前缀", "A0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48", false},
		{"空字符串", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidEVMAddress(tt.addr)
			assert.Equal(t, tt.isValid, result)
		})
	}
}

func TestValidateAmount(t *testing.T) {
	tests := []struct {
		name    string
		amount  string
		isValid bool
	}{
		{"合法数量", "1000000", true},
		{"非法字符", "abc", false},
		{"空字符串", "", false},
		{"含小数点", "1.5", false},
		{"零", "0", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidAmount(tt.amount)
			assert.Equal(t, tt.isValid, result)
		})
	}
}

func TestValidateTxHash(t *testing.T) {
	tests := []struct {
		name    string
		hash    string
		isValid bool
	}{
		{"合法哈希", "0xabc123def456abc123def456abc123def456abc123def456abc123def456abcd", true},
		{"无效格式", "not_a_hash", false},
		{"太短", "0xabc", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidTxHash(tt.hash)
			assert.Equal(t, tt.isValid, result)
		})
	}
}

func TestHandler_GetProviders(t *testing.T) {
	t.Run("成功获取供应商列表", func(t *testing.T) {
		registry := NewProviderRegistry()
		registry.Register(&mockQuoteProvider{})
		cache := NewMemoryCache(10 * time.Second)
		repo := newMockRepo()
		svc := NewQuoteService(registry, cache, repo, 1)

		router := setupTestRouter(svc)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/providers", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, float64(0), resp["code"])
	})

	t.Run("无已注册供应商", func(t *testing.T) {
		registry := NewProviderRegistry()
		cache := NewMemoryCache(10 * time.Second)
		repo := newMockRepo()
		svc := NewQuoteService(registry, cache, repo, 1)

		router := setupTestRouter(svc)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/providers", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestHandler_GetQuote(t *testing.T) {
	setup := func() (*gin.Engine, *mockRepo) {
		registry := NewProviderRegistry()
		registry.Register(&mockQuoteProvider{
			quoteResult: &domain.QuoteResult{
				Provider:   "test-provider",
				FromAmount: "1000",
				ToAmount:   "2000",
			},
		})
		cache := NewMemoryCache(10 * time.Second)
		repo := newMockRepo()
		svc := NewQuoteService(registry, cache, repo, 1)
		return setupTestRouter(svc), repo
	}

	t.Run("缺少必填参数", func(t *testing.T) {
		router, _ := setup()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/quote?provider=test-provider", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Equal(t, float64(ErrCodeInvalidParam), resp["code"])
	})

	t.Run("无效代币地址格式", func(t *testing.T) {
		router, _ := setup()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/quote?provider=test-provider&chain_id=1&from_token=invalid_address&to_token=0xB&amount=1000&swap_mode=exactIn", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("无效 amount 格式", func(t *testing.T) {
		router, _ := setup()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/quote?provider=test-provider&chain_id=1&from_token=0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48&to_token=0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2&amount=abc&swap_mode=exactIn", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("原生代币地址合法", func(t *testing.T) {
		router, _ := setup()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/api/v1/quote?provider=test-provider&chain_id=1&from_token=0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE&to_token=0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2&amount=1000000&swap_mode=exactIn", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestHandler_GetSwap(t *testing.T) {
	setup := func() *gin.Engine {
		registry := NewProviderRegistry()
		registry.Register(&mockQuoteProvider{
			swapResult: &domain.SwapResult{
				QuoteResult: domain.QuoteResult{
					Provider:   "test-provider",
					FromAmount: "1000",
					ToAmount:   "2000",
				},
				TxData: domain.TxData{To: "0xContract", Data: "0x"},
			},
		})
		cache := NewMemoryCache(10 * time.Second)
		repo := newMockRepo()
		svc := NewQuoteService(registry, cache, repo, 1)
		return setupTestRouter(svc)
	}

	makeBody := func(overrides map[string]interface{}) *bytes.Buffer {
		body := map[string]interface{}{
			"provider":            "test-provider",
			"chain_id":            1,
			"from_token":          "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
			"to_token":            "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
			"amount":              "1000",
			"swap_mode":           "exactIn",
			"slippage_percent":    0.5,
			"user_wallet_address": "0x1234567890123456789012345678901234567890",
		}
		for k, v := range overrides {
			if v == nil {
				delete(body, k)
			} else {
				body[k] = v
			}
		}
		b, _ := json.Marshal(body)
		return bytes.NewBuffer(b)
	}

	t.Run("缺少 user_wallet_address", func(t *testing.T) {
		router := setup()
		w := httptest.NewRecorder()
		body := makeBody(map[string]interface{}{"user_wallet_address": nil})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/swap", body)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("slippage_percent 超出范围", func(t *testing.T) {
		router := setup()
		w := httptest.NewRecorder()
		body := makeBody(map[string]interface{}{"slippage_percent": 101.0})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/swap", body)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("无效钱包地址", func(t *testing.T) {
		router := setup()
		w := httptest.NewRecorder()
		body := makeBody(map[string]interface{}{"user_wallet_address": "0x1234"})
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/swap", body)
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestHandler_SubmitSwap(t *testing.T) {
	t.Run("无效 tx_hash 格式", func(t *testing.T) {
		registry := NewProviderRegistry()
		cache := NewMemoryCache(10 * time.Second)
		repo := newMockRepo()
		svc := NewQuoteService(registry, cache, repo, 1)
		router := setupTestRouter(svc)

		body, _ := json.Marshal(map[string]string{"tx_hash": "not_a_hash"})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodPost, "/api/v1/swap/123/submit", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestValidationRegexes(t *testing.T) {
	// 测试 evmAddressRegex
	assert.True(t, regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`).MatchString("0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"))
	assert.False(t, regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`).MatchString("0x1234"))

	// 测试 txHashRegex
	assert.True(t, regexp.MustCompile(`^0x[0-9a-fA-F]{64}$`).MatchString("0xabc123def456abc123def456abc123def456abc123def456abc123def456abcd"))
	assert.False(t, regexp.MustCompile(`^0x[0-9a-fA-F]{64}$`).MatchString("not_a_hash"))
}
