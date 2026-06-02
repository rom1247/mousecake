package okx

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/mousecake-go/mousecake-go/internal/quote/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testConfig 返回用于测试的 OKX 配置，BaseURL 将由测试服务器替换。
func testConfig() Config {
	return Config{
		APIKey:     "test-api-key",
		SecretKey:  "test-secret-key",
		Passphrase: "test-passphrase",
		BaseURL:    "https://www.okx.com",
		Timeout:    5 * time.Second,
	}
}

// TestSign 验证 HMAC-SHA256 签名生成结果。
func TestSign(t *testing.T) {
	t.Parallel()

	// 使用固定参数验证签名算法正确性
	timestamp := "2023-10-18T12:21:41.274Z"
	method := "GET"
	requestPath := "/api/v6/dex/aggregator/quote"
	queryString := "amount=10000000000000000000&chainIndex=1&fromTokenAddress=0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee&toTokenAddress=0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"
	secretKey := "test-secret-key"

	result := sign(timestamp, method, requestPath, queryString, secretKey)

	// 签名应是非空的 Base64 字符串
	assert.NotEmpty(t, result, "签名不应为空")

	// 相同输入应产生相同输出
	result2 := sign(timestamp, method, requestPath, queryString, secretKey)
	assert.Equal(t, result, result2, "相同输入应产生相同签名")

	// 不同输入应产生不同输出
	result3 := sign("2023-10-18T12:21:42.274Z", method, requestPath, queryString, secretKey)
	assert.NotEqual(t, result, result3, "不同时间戳应产生不同签名")
}

// TestFormatTimestamp 验证时间戳格式为 ISO 8601 毫秒精度。
func TestFormatTimestamp(t *testing.T) {
	t.Parallel()

	// 解析固定时间并验证格式化输出
	now := time.Date(2023, 10, 18, 12, 21, 41, 274e6, time.UTC)
	ts := formatTimestamp(now)
	assert.Equal(t, "2023-10-18T12:21:41.274Z", ts)
}

// TestBuildQuoteParams 验证 quote 参数映射正确。
func TestBuildQuoteParams(t *testing.T) {
	t.Parallel()

	params := buildQuoteParams(1, "0xFrom", "0xTo", "1000000", "exactIn")

	assert.Equal(t, "1", params.Get("chainIndex"))
	assert.Equal(t, "0xFrom", params.Get("fromTokenAddress"))
	assert.Equal(t, "0xTo", params.Get("toTokenAddress"))
	assert.Equal(t, "1000000", params.Get("amount"))
	assert.Equal(t, "exactIn", params.Get("swapMode"))
}

// TestBuildSwapParams 验证 swap 参数映射正确。
func TestBuildSwapParams(t *testing.T) {
	t.Parallel()

	params := buildSwapParams(
		1, "0xFrom", "0xTo", "1000000", "exactIn",
		0.5, "0xUserWallet",
	)

	assert.Equal(t, "1", params.Get("chainIndex"))
	assert.Equal(t, "0xFrom", params.Get("fromTokenAddress"))
	assert.Equal(t, "0xTo", params.Get("toTokenAddress"))
	assert.Equal(t, "1000000", params.Get("amount"))
	assert.Equal(t, "exactIn", params.Get("swapMode"))
	assert.Equal(t, "0.5", params.Get("slippagePercent"))
	assert.Equal(t, "0xUserWallet", params.Get("userWalletAddress"))
}

// TestGetQuote_Success 验证成功获取报价。
func TestGetQuote_Success(t *testing.T) {
	t.Parallel()

	// 准备 mock 响应
	quoteResp := okxResponse{
		Code: "0",
		Msg:  "",
	}
	quoteDataItem := quoteData{
		ChainIndex:         "1",
		FromTokenAmount:    "1000000000000000000",
		ToTokenAmount:      "500000000",
		PriceImpactPercent: "0.01",
		EstimateGasFee:     "21000",
		Router:             "0xeeee--0xa0b8",
		FromToken: tokenData{
			TokenContractAddress: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
			TokenSymbol:          "ETH",
			Decimal:              "18",
		},
		ToToken: tokenData{
			TokenContractAddress: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
			TokenSymbol:          "USDC",
			Decimal:              "6",
		},
		SwapMode: "exactIn",
	}
	dataBytes, err := json.Marshal(quoteDataItem)
	require.NoError(t, err)
	quoteResp.Data = []json.RawMessage{dataBytes}
	respBody, err := json.Marshal(quoteResp)
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求方法和路径
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, quotePath, r.URL.Path)

		// 验证认证头存在
		assert.NotEmpty(t, r.Header.Get("OK-ACCESS-KEY"))
		assert.NotEmpty(t, r.Header.Get("OK-ACCESS-SIGN"))
		assert.NotEmpty(t, r.Header.Get("OK-ACCESS-TIMESTAMP"))
		assert.NotEmpty(t, r.Header.Get("OK-ACCESS-PASSPHRASE"))

		// 验证查询参数
		assert.Equal(t, "1", r.URL.Query().Get("chainIndex"))
		assert.Equal(t, "0xFrom", r.URL.Query().Get("fromTokenAddress"))
		assert.Equal(t, "0xTo", r.URL.Query().Get("toTokenAddress"))
		assert.Equal(t, "1000000", r.URL.Query().Get("amount"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(respBody)
	}))
	defer server.Close()

	cfg := testConfig()
	cfg.BaseURL = server.URL
	provider := NewProvider(cfg)

	result, err := provider.GetQuote(context.Background(), domain.QuoteParams{
		ChainID:   1,
		FromToken: "0xFrom",
		ToToken:   "0xTo",
		Amount:    "1000000",
		SwapMode:  domain.SwapModeExactIn,
	})
	require.NoError(t, err)

	assert.Equal(t, "okx", result.Provider)
	assert.Equal(t, "1000000000000000000", result.FromAmount)
	assert.Equal(t, "500000000", result.ToAmount)
	assert.Equal(t, "0.01", result.PriceImpact)
	assert.Equal(t, "21000", result.EstimatedGas)
	assert.Equal(t, "ETH", result.FromToken.Symbol)
	assert.Equal(t, 18, result.FromToken.Decimals)
	assert.Equal(t, "USDC", result.ToToken.Symbol)
	assert.Equal(t, 6, result.ToToken.Decimals)
}

// TestGetQuote_BusinessError 验证 OKX 业务错误（code != "0"）处理。
func TestGetQuote_BusinessError(t *testing.T) {
	t.Parallel()

	errorResp := okxResponse{
		Code: "51001",
		Msg:  "Invalid parameter",
		Data: []json.RawMessage{},
	}
	respBody, err := json.Marshal(errorResp)
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(respBody)
	}))
	defer server.Close()

	cfg := testConfig()
	cfg.BaseURL = server.URL
	provider := NewProvider(cfg)

	_, err = provider.GetQuote(context.Background(), domain.QuoteParams{
		ChainID:   1,
		FromToken: "0xFrom",
		ToToken:   "0xTo",
		Amount:    "1000000",
		SwapMode:  domain.SwapModeExactIn,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "51001")
	assert.Contains(t, err.Error(), "Invalid parameter")
}

// TestGetQuote_CodeIsString 验证 OKX code 字段为 string 类型。
// 如果 code 使用 int 接收，"0" 会被解析为 0，而非零错误码如 "51001" 也能正常解析。
// 此测试确保 code 字段始终作为 string 处理。
func TestGetQuote_CodeIsString(t *testing.T) {
	t.Parallel()

	// 构造一个包含 "0"（字符串零）的响应
	resp := map[string]interface{}{
		"code": "0",
		"msg":  "",
		"data": []interface{}{
			map[string]interface{}{
				"chainIndex":         "1",
				"fromTokenAmount":    "1000",
				"toTokenAmount":      "2000",
				"priceImpactPercent": "0.05",
				"estimateGasFee":     "50000",
				"fromToken": map[string]interface{}{
					"tokenContractAddress": "0xFrom",
					"tokenSymbol":          "TK1",
					"decimal":              "18",
				},
				"toToken": map[string]interface{}{
					"tokenContractAddress": "0xTo",
					"tokenSymbol":          "TK2",
					"decimal":              "6",
				},
			},
		},
	}
	respBody, err := json.Marshal(resp)
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(respBody)
	}))
	defer server.Close()

	cfg := testConfig()
	cfg.BaseURL = server.URL
	provider := NewProvider(cfg)

	result, err := provider.GetQuote(context.Background(), domain.QuoteParams{
		ChainID:   1,
		FromToken: "0xFrom",
		ToToken:   "0xTo",
		Amount:    "1000",
		SwapMode:  domain.SwapModeExactIn,
	})
	require.NoError(t, err)
	assert.Equal(t, "okx", result.Provider)
}

// TestGetQuote_EmptyData 验证空 data 数组的错误处理。
func TestGetQuote_EmptyData(t *testing.T) {
	t.Parallel()

	resp := okxResponse{Code: "0", Msg: "", Data: []json.RawMessage{}}
	respBody, err := json.Marshal(resp)
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(respBody)
	}))
	defer server.Close()

	cfg := testConfig()
	cfg.BaseURL = server.URL
	provider := NewProvider(cfg)

	_, err = provider.GetQuote(context.Background(), domain.QuoteParams{
		ChainID:   1,
		FromToken: "0xFrom",
		ToToken:   "0xTo",
		Amount:    "1000",
		SwapMode:  domain.SwapModeExactIn,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "响应数据为空")
}

// TestGetQuote_HTTPError 验证 HTTP 非 200 的错误处理。
func TestGetQuote_HTTPError(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal server error"))
	}))
	defer server.Close()

	cfg := testConfig()
	cfg.BaseURL = server.URL
	provider := NewProvider(cfg)

	_, err := provider.GetQuote(context.Background(), domain.QuoteParams{
		ChainID:   1,
		FromToken: "0xFrom",
		ToToken:   "0xTo",
		Amount:    "1000",
		SwapMode:  domain.SwapModeExactIn,
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "HTTP 500")
}

// TestGetSwap_Success 验证成功获取 swap 数据。
func TestGetSwap_Success(t *testing.T) {
	t.Parallel()

	swapResp := okxResponse{Code: "0", Msg: ""}
	swapDataItem := swapData{
		RouterResult: routerResult{
			ChainIndex:         "1",
			FromTokenAmount:    "100000000",
			ToTokenAmount:      "50000000",
			PriceImpactPercent: "0.02",
			EstimateGasFee:     "150000",
			Router:             "0xa0b8--0x2260",
			FromToken: tokenData{
				TokenContractAddress: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
				TokenSymbol:          "USDC",
				Decimal:              "6",
			},
			ToToken: tokenData{
				TokenContractAddress: "0x2260fac5e5542a773aa44fbcfedf7c193bc2c599",
				TokenSymbol:          "WBTC",
				Decimal:              "8",
			},
		},
		Tx: txData{
			From:     "0x77660f108043c9e300b4e30a35a61dd19f5ae28a",
			To:       "0x5E1f62Dac767b0491e3CE72469C217365D5B48cC",
			Data:     "0xf2c42696...",
			Value:    "0",
			Gas:      "1248837",
			GasPrice: "557703374",
		},
	}
	dataBytes, err := json.Marshal(swapDataItem)
	require.NoError(t, err)
	swapResp.Data = []json.RawMessage{dataBytes}
	respBody, err := json.Marshal(swapResp)
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, swapPath, r.URL.Path)
		assert.Equal(t, "0xUserWallet", r.URL.Query().Get("userWalletAddress"))
		assert.Equal(t, "0.5", r.URL.Query().Get("slippagePercent"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(respBody)
	}))
	defer server.Close()

	cfg := testConfig()
	cfg.BaseURL = server.URL
	provider := NewProvider(cfg)

	result, err := provider.GetSwap(context.Background(), domain.SwapParams{
		QuoteParams: domain.QuoteParams{
			ChainID:   1,
			FromToken: "0xFrom",
			ToToken:   "0xTo",
			Amount:    "100000000",
			SwapMode:  domain.SwapModeExactIn,
		},
		SlippagePercent:   0.5,
		UserWalletAddress: "0xUserWallet",
	})
	require.NoError(t, err)

	// 验证 QuoteResult 部分
	assert.Equal(t, "okx", result.Provider)
	assert.Equal(t, "100000000", result.FromAmount)
	assert.Equal(t, "50000000", result.ToAmount)
	assert.Equal(t, "0.02", result.PriceImpact)
	assert.Equal(t, "150000", result.EstimatedGas)
	assert.Equal(t, "USDC", result.FromToken.Symbol)
	assert.Equal(t, "WBTC", result.ToToken.Symbol)

	// 验证 TxData 部分
	assert.Equal(t, "0x5E1f62Dac767b0491e3CE72469C217365D5B48cC", result.TxData.To)
	assert.Equal(t, "0xf2c42696...", result.TxData.Data)
	assert.Equal(t, "0", result.TxData.Value)
	assert.Equal(t, "1248837", result.TxData.Gas)
	assert.Equal(t, "557703374", result.TxData.GasPrice)
}

// TestName 验证供应商名称。
func TestName(t *testing.T) {
	t.Parallel()

	provider := NewProvider(testConfig())
	assert.Equal(t, "okx", provider.Name())
}

// TestNewProvider_DefaultConfig 验证默认配置填充。
func TestNewProvider_DefaultConfig(t *testing.T) {
	t.Parallel()

	cfg := Config{
		APIKey:     "key",
		SecretKey:  "secret",
		Passphrase: "pass",
	}
	provider := NewProvider(cfg)
	require.NotNil(t, provider)
	assert.Equal(t, "okx", provider.Name())
}

// TestBuildAuthHeaders 验证认证头包含所有必需字段。
func TestBuildAuthHeaders(t *testing.T) {
	t.Parallel()

	cfg := testConfig()
	client := newHTTPClient(cfg)
	headers := client.buildAuthHeaders("/api/v6/dex/aggregator/quote", "chainIndex=1")

	assert.Equal(t, "test-api-key", headers.Get("OK-ACCESS-KEY"))
	assert.NotEmpty(t, headers.Get("OK-ACCESS-SIGN"))
	assert.NotEmpty(t, headers.Get("OK-ACCESS-TIMESTAMP"))
	assert.Equal(t, "test-passphrase", headers.Get("OK-ACCESS-PASSPHRASE"))
}

// TestSign_Deterministic 验证签名对同一输入是确定性的。
func TestSign_Deterministic(t *testing.T) {
	t.Parallel()

	s1 := sign("2023-10-18T12:21:41.274Z", "GET", "/path", "key=val", "secret")
	s2 := sign("2023-10-18T12:21:41.274Z", "GET", "/path", "key=val", "secret")
	assert.Equal(t, s1, s2, "相同输入的签名必须一致")
}

// TestBuildQuoteParams_Encoding 验证参数编码使用 url.Values（按字母序排列）。
func TestBuildQuoteParams_Encoding(t *testing.T) {
	t.Parallel()

	params := buildQuoteParams(1, "0xA", "0xB", "1000", "exactIn")
	encoded := params.Encode()

	// url.Values.Encode() 保证参数按 key 字母序排列
	assert.Contains(t, encoded, "amount=1000")
	assert.Contains(t, encoded, "chainIndex=1")
	assert.Contains(t, encoded, "fromTokenAddress=0xA")
	assert.Contains(t, encoded, "toTokenAddress=0xB")
	assert.Contains(t, encoded, "swapMode=exactIn")
}
