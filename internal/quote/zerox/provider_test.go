package zerox

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/mousecake-go/mousecake-go/internal/quote/domain"
)

const (
	// testAPIKey 测试用 API Key。
	testAPIKey = "test-api-key-12345"
	// testFromToken 测试用源代币地址（WETH）。
	testFromToken = "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"
	// testToToken 测试用目标代币地址（DAI）。
	testToToken = "0x6b175474e89094c44da98b954eedeac495271d0f"
	// testAmount 测试用数量。
	testAmount = "100000000000000000000"
	// testTaker 测试用钱包地址。
	testTaker = "0x3f6a3f57569358a512ccc0e513f171516b0fd42a"
)

// newTestServer 创建测试用 HTTP 服务器，handler 函数用于自定义响应。
func newTestServer(handler http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(handler)
}

// newTestProvider 基于测试服务器 URL 创建 Provider。
func newTestProvider(serverURL string) *Provider {
	return NewProvider(Config{
		APIKey:  testAPIKey,
		BaseURL: serverURL,
		Timeout: 0, // 使用默认值
	})
}

// --- 认证头测试 ---

func TestGetQuote_包含认证Header(t *testing.T) {
	var gotAPIKey, gotVersion string
	server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		gotAPIKey = r.Header.Get(headerAPIKey)
		gotVersion = r.Header.Get(headerVersion)

		resp := priceResponse{
			SellAmount: testAmount,
			SellToken:  testFromToken,
			BuyAmount:  "458050129388884000000000",
			BuyToken:   testToToken,
			Gas:        "1116817",
			GasPrice:   "2558459858",
		}
		writeJSON(w, resp)
	})
	defer server.Close()

	provider := newTestProvider(server.URL)
	_, err := provider.GetQuote(context.Background(), domain.QuoteParams{
		ChainID:   1,
		FromToken: testFromToken,
		ToToken:   testToToken,
		Amount:    testAmount,
		SwapMode:  domain.SwapModeExactIn,
	})
	if err != nil {
		t.Fatalf("不期望错误: %v", err)
	}

	if gotAPIKey != testAPIKey {
		t.Errorf("期望 0x-api-key=%s, 得到 %s", testAPIKey, gotAPIKey)
	}
	if gotVersion != headerVersionValue {
		t.Errorf("期望 0x-version=%s, 得到 %s", headerVersionValue, gotVersion)
	}
}

// --- 参数映射测试 ---

func TestGetQuote_exactIn模式使用sellAmount(t *testing.T) {
	var gotSellAmount, gotBuyAmount, gotSellToken, gotBuyToken string
	server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		gotSellToken = r.URL.Query().Get("sellToken")
		gotBuyToken = r.URL.Query().Get("buyToken")
		gotSellAmount = r.URL.Query().Get("sellAmount")
		gotBuyAmount = r.URL.Query().Get("buyAmount")

		resp := priceResponse{
			SellAmount: testAmount,
			SellToken:  testFromToken,
			BuyAmount:  "458050129388884000000000",
			BuyToken:   testToToken,
			Gas:        "1116817",
		}
		writeJSON(w, resp)
	})
	defer server.Close()

	provider := newTestProvider(server.URL)
	_, err := provider.GetQuote(context.Background(), domain.QuoteParams{
		ChainID:   1,
		FromToken: testFromToken,
		ToToken:   testToToken,
		Amount:    testAmount,
		SwapMode:  domain.SwapModeExactIn,
	})
	if err != nil {
		t.Fatalf("不期望错误: %v", err)
	}

	if gotSellToken != testFromToken {
		t.Errorf("期望 sellToken=%s, 得到 %s", testFromToken, gotSellToken)
	}
	if gotBuyToken != testToToken {
		t.Errorf("期望 buyToken=%s, 得到 %s", testToToken, gotBuyToken)
	}
	if gotSellAmount != testAmount {
		t.Errorf("期望 sellAmount=%s, 得到 %s", testAmount, gotSellAmount)
	}
	if gotBuyAmount != "" {
		t.Errorf("exactIn 模式不应设置 buyAmount, 得到 %s", gotBuyAmount)
	}
}

func TestGetQuote_exactOut模式使用buyAmount(t *testing.T) {
	var gotSellAmount, gotBuyAmount string
	server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		gotSellAmount = r.URL.Query().Get("sellAmount")
		gotBuyAmount = r.URL.Query().Get("buyAmount")

		resp := priceResponse{
			SellAmount: "500000000000000000000",
			SellToken:  testFromToken,
			BuyAmount:  testAmount,
			BuyToken:   testToToken,
			Gas:        "1116817",
		}
		writeJSON(w, resp)
	})
	defer server.Close()

	provider := newTestProvider(server.URL)
	_, err := provider.GetQuote(context.Background(), domain.QuoteParams{
		ChainID:   1,
		FromToken: testFromToken,
		ToToken:   testToToken,
		Amount:    testAmount,
		SwapMode:  domain.SwapModeExactOut,
	})
	if err != nil {
		t.Fatalf("不期望错误: %v", err)
	}

	if gotSellAmount != "" {
		t.Errorf("exactOut 模式不应设置 sellAmount, 得到 %s", gotSellAmount)
	}
	if gotBuyAmount != testAmount {
		t.Errorf("期望 buyAmount=%s, 得到 %s", testAmount, gotBuyAmount)
	}
}

// --- chain_id 嵌入 URL path 测试 ---

func TestGetQuote_chainID嵌入URLPath(t *testing.T) {
	var gotPath string
	server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path

		resp := priceResponse{
			SellAmount: testAmount,
			SellToken:  testFromToken,
			BuyAmount:  "458050129388884000000000",
			BuyToken:   testToToken,
			Gas:        "1116817",
		}
		writeJSON(w, resp)
	})
	defer server.Close()

	provider := newTestProvider(server.URL)
	chainID := 137
	_, err := provider.GetQuote(context.Background(), domain.QuoteParams{
		ChainID:   chainID,
		FromToken: testFromToken,
		ToToken:   testToToken,
		Amount:    testAmount,
		SwapMode:  domain.SwapModeExactIn,
	})
	if err != nil {
		t.Fatalf("不期望错误: %v", err)
	}

	expected := "/137/swap/allowance-holder/price"
	if gotPath != expected {
		t.Errorf("期望 path=%s, 得到 %s", expected, gotPath)
	}
}

// --- 成功获取报价测试 ---

func TestGetQuote_成功获取报价(t *testing.T) {
	server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		resp := priceResponse{
			SellAmount: "100000000000000000000",
			SellToken:  testFromToken,
			BuyAmount:  "458050129388884000000000",
			BuyToken:   testToToken,
			Gas:        "1116817",
			GasPrice:   "2558459858",
			Route: routeInfo{
				Fills: []routeFill{
					{
						From:          testFromToken,
						To:            "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
						Source:        "Uniswap_V3",
						ProportionBps: "500",
					},
					{
						From:          "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
						To:            testToToken,
						Source:        "Maker_PSM",
						ProportionBps: "10000",
					},
				},
			},
		}
		writeJSON(w, resp)
	})
	defer server.Close()

	provider := newTestProvider(server.URL)
	result, err := provider.GetQuote(context.Background(), domain.QuoteParams{
		ChainID:   1,
		FromToken: testFromToken,
		ToToken:   testToToken,
		Amount:    testAmount,
		SwapMode:  domain.SwapModeExactIn,
	})
	if err != nil {
		t.Fatalf("不期望错误: %v", err)
	}

	if result.Provider != "zerox" {
		t.Errorf("期望 provider=zerox, 得到 %s", result.Provider)
	}
	if result.FromToken.Address != testFromToken {
		t.Errorf("期望 fromToken=%s, 得到 %s", testFromToken, result.FromToken.Address)
	}
	if result.ToToken.Address != testToToken {
		t.Errorf("期望 toToken=%s, 得到 %s", testToToken, result.ToToken.Address)
	}
	if result.FromAmount != "100000000000000000000" {
		t.Errorf("期望 fromAmount=100000000000000000000, 得到 %s", result.FromAmount)
	}
	if result.ToAmount != "458050129388884000000000" {
		t.Errorf("期望 toAmount=458050129388884000000000, 得到 %s", result.ToAmount)
	}
	if result.EstimatedGas != "1116817" {
		t.Errorf("期望 estimatedGas=1116817, 得到 %s", result.EstimatedGas)
	}
	if len(result.Route) != 2 {
		t.Fatalf("期望 2 个路由跳, 得到 %d", len(result.Route))
	}
	if result.Route[0].PoolAddress != "Uniswap_V3" {
		t.Errorf("期望 route[0].pool=Uniswap_V3, 得到 %s", result.Route[0].PoolAddress)
	}
}

// --- 0x API 错误处理测试 ---

func TestGetQuote_API错误处理(t *testing.T) {
	server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		errResp := zeroxError{
			Code:   400,
			Reason: "Invalid token",
		}
		writeJSON(w, errResp)
	})
	defer server.Close()

	provider := newTestProvider(server.URL)
	_, err := provider.GetQuote(context.Background(), domain.QuoteParams{
		ChainID:   1,
		FromToken: "invalid",
		ToToken:   testToToken,
		Amount:    testAmount,
		SwapMode:  domain.SwapModeExactIn,
	})
	if err == nil {
		t.Fatal("期望错误, 得到 nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("期望 *APIError, 得到 %T: %v", err, err)
	}
	if apiErr.StatusCode != 400 {
		t.Errorf("期望 statusCode=400, 得到 %d", apiErr.StatusCode)
	}
	if apiErr.Reason != "Invalid token" {
		t.Errorf("期望 reason=Invalid token, 得到 %s", apiErr.Reason)
	}
}

// --- 限流处理测试 ---

func TestGetQuote_限流处理(t *testing.T) {
	server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`{"code":429,"reason":"Rate limit exceeded"}`))
	})
	defer server.Close()

	provider := newTestProvider(server.URL)
	_, err := provider.GetQuote(context.Background(), domain.QuoteParams{
		ChainID:   1,
		FromToken: testFromToken,
		ToToken:   testToToken,
		Amount:    testAmount,
		SwapMode:  domain.SwapModeExactIn,
	})
	if err == nil {
		t.Fatal("期望错误, 得到 nil")
	}

	var rateLimitErr *RateLimitError
	if !errors.As(err, &rateLimitErr) {
		t.Fatalf("期望 *RateLimitError, 得到 %T: %v", err, err)
	}
	if rateLimitErr.StatusCode != 429 {
		t.Errorf("期望 statusCode=429, 得到 %d", rateLimitErr.StatusCode)
	}
}

// --- Swap 参数映射测试 ---

func TestGetSwap_slippageBps转换(t *testing.T) {
	var gotSlippageBps, gotTaker string
	server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		gotSlippageBps = r.URL.Query().Get("slippageBps")
		gotTaker = r.URL.Query().Get("taker")

		resp := quoteResponse{
			SellAmount: "100000000000000000000",
			SellToken:  testFromToken,
			BuyAmount:  "458050129388884000000000",
			BuyToken:   testToToken,
			Gas:        "832055",
			GasPrice:   "1676486955",
			Transaction: &transactionData{
				To:       "0x0000000000001ff3684f28c67538d4d072c22734",
				Data:     "0x2213bc0b",
				Value:    "0",
				Gas:      "832055",
				GasPrice: "1676486955",
			},
		}
		writeJSON(w, resp)
	})
	defer server.Close()

	provider := newTestProvider(server.URL)
	_, err := provider.GetSwap(context.Background(), domain.SwapParams{
		QuoteParams: domain.QuoteParams{
			ChainID:   1,
			FromToken: testFromToken,
			ToToken:   testToToken,
			Amount:    testAmount,
			SwapMode:  domain.SwapModeExactIn,
		},
		SlippagePercent:   0.5,
		UserWalletAddress: testTaker,
	})
	if err != nil {
		t.Fatalf("不期望错误: %v", err)
	}

	// 0.5% → 50 bps
	if gotSlippageBps != "50" {
		t.Errorf("期望 slippageBps=50, 得到 %s", gotSlippageBps)
	}
	if gotTaker != testTaker {
		t.Errorf("期望 taker=%s, 得到 %s", testTaker, gotTaker)
	}
}

func TestGetSwap_slippageBps整数百分比(t *testing.T) {
	tests := []struct {
		name            string
		slippagePercent float64
		expectedBps     string
	}{
		{"0.1%", 0.1, "10"},
		{"0.5%", 0.5, "50"},
		{"1%", 1.0, "100"},
		{"3%", 3.0, "300"},
		{"5%", 5.0, "500"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotSlippageBps string
			server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
				gotSlippageBps = r.URL.Query().Get("slippageBps")

				resp := quoteResponse{
					SellAmount: testAmount,
					SellToken:  testFromToken,
					BuyAmount:  "458050129388884000000000",
					BuyToken:   testToToken,
					Gas:        "832055",
					Transaction: &transactionData{
						To:   "0x1234",
						Data: "0xabc",
					},
				}
				writeJSON(w, resp)
			})
			defer server.Close()

			provider := newTestProvider(server.URL)
			_, err := provider.GetSwap(context.Background(), domain.SwapParams{
				QuoteParams: domain.QuoteParams{
					ChainID:   1,
					FromToken: testFromToken,
					ToToken:   testToToken,
					Amount:    testAmount,
					SwapMode:  domain.SwapModeExactIn,
				},
				SlippagePercent:   tt.slippagePercent,
				UserWalletAddress: testTaker,
			})
			if err != nil {
				t.Fatalf("不期望错误: %v", err)
			}

			if gotSlippageBps != tt.expectedBps {
				t.Errorf("期望 slippageBps=%s, 得到 %s", tt.expectedBps, gotSlippageBps)
			}
		})
	}
}

// --- 成功获取 Swap 测试 ---

func TestGetSwap_成功获取Swap数据(t *testing.T) {
	server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		// 验证路径包含 quote
		if r.URL.Path != "/1/swap/allowance-holder/quote" {
			t.Errorf("期望 path=/1/swap/allowance-holder/quote, 得到 %s", r.URL.Path)
		}

		resp := quoteResponse{
			SellAmount: "100000000000000000000",
			SellToken:  testFromToken,
			BuyAmount:  "458050129388884000000000",
			BuyToken:   testToToken,
			Gas:        "832055",
			GasPrice:   "1676486955",
			Route: routeInfo{
				Fills: []routeFill{
					{
						From:          testFromToken,
						To:            testToToken,
						Source:        "Uniswap_V3",
						ProportionBps: "10000",
					},
				},
			},
			Transaction: &transactionData{
				To:       "0x0000000000001ff3684f28c67538d4d072c22734",
				Data:     "0x2213bc0b000000000000000000000000df31a70a21a1931e02033dbba7deace6c45",
				Value:    "0",
				Gas:      "832055",
				GasPrice: "1676486955",
			},
		}
		writeJSON(w, resp)
	})
	defer server.Close()

	provider := newTestProvider(server.URL)
	result, err := provider.GetSwap(context.Background(), domain.SwapParams{
		QuoteParams: domain.QuoteParams{
			ChainID:   1,
			FromToken: testFromToken,
			ToToken:   testToToken,
			Amount:    testAmount,
			SwapMode:  domain.SwapModeExactIn,
		},
		SlippagePercent:   0.5,
		UserWalletAddress: testTaker,
	})
	if err != nil {
		t.Fatalf("不期望错误: %v", err)
	}

	// 验证 QuoteResult
	if result.Provider != "zerox" {
		t.Errorf("期望 provider=zerox, 得到 %s", result.Provider)
	}
	if result.FromAmount != "100000000000000000000" {
		t.Errorf("期望 fromAmount=100000000000000000000, 得到 %s", result.FromAmount)
	}
	if result.ToAmount != "458050129388884000000000" {
		t.Errorf("期望 toAmount=458050129388884000000000, 得到 %s", result.ToAmount)
	}

	// 验证 TxData
	if result.TxData.To != "0x0000000000001ff3684f28c67538d4d072c22734" {
		t.Errorf("期望 tx.to=0x0000000000001ff3684f28c67538d4d072c22734, 得到 %s", result.TxData.To)
	}
	if result.TxData.Data != "0x2213bc0b000000000000000000000000df31a70a21a1931e02033dbba7deace6c45" {
		t.Errorf("期望 tx.data 匹配, 得到 %s", result.TxData.Data)
	}
	if result.TxData.Value != "0" {
		t.Errorf("期望 tx.value=0, 得到 %s", result.TxData.Value)
	}
	if result.TxData.Gas != "832055" {
		t.Errorf("期望 tx.gas=832055, 得到 %s", result.TxData.Gas)
	}
	if result.TxData.GasPrice != "1676486955" {
		t.Errorf("期望 tx.gasPrice=1676486955, 得到 %s", result.TxData.GasPrice)
	}

	// 验证路由
	if len(result.Route) != 1 {
		t.Fatalf("期望 1 个路由跳, 得到 %d", len(result.Route))
	}
}

// --- Name 测试 ---

func TestName_返回zerox(t *testing.T) {
	provider := NewProvider(Config{APIKey: testAPIKey})
	if provider.Name() != "zerox" {
		t.Errorf("期望 name=zerox, 得到 %s", provider.Name())
	}
}

// --- 辅助函数 ---

// writeJSON 将对象序列化为 JSON 并写入响应。
func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, _ := json.Marshal(v)
	w.Write(data)
}

// --- buildPriceParams 单元测试 ---

func TestBuildPriceParams_exactIn(t *testing.T) {
	params := buildPriceParams(1, testFromToken, testToToken, testAmount, false)

	if params.Get("chainId") != "1" {
		t.Errorf("期望 chainId=1, 得到 %s", params.Get("chainId"))
	}
	if params.Get("sellToken") != testFromToken {
		t.Errorf("期望 sellToken=%s, 得到 %s", testFromToken, params.Get("sellToken"))
	}
	if params.Get("buyToken") != testToToken {
		t.Errorf("期望 buyToken=%s, 得到 %s", testToToken, params.Get("buyToken"))
	}
	if params.Get("sellAmount") != testAmount {
		t.Errorf("期望 sellAmount=%s, 得到 %s", testAmount, params.Get("sellAmount"))
	}
	if params.Get("buyAmount") != "" {
		t.Errorf("exactIn 不应设置 buyAmount, 得到 %s", params.Get("buyAmount"))
	}
}

func TestBuildPriceParams_exactOut(t *testing.T) {
	params := buildPriceParams(1, testFromToken, testToToken, testAmount, true)

	if params.Get("buyAmount") != testAmount {
		t.Errorf("期望 buyAmount=%s, 得到 %s", testAmount, params.Get("buyAmount"))
	}
	if params.Get("sellAmount") != "" {
		t.Errorf("exactOut 不应设置 sellAmount, 得到 %s", params.Get("sellAmount"))
	}
}

func TestBuildQuoteParams_包含Taker和Slippage(t *testing.T) {
	params := buildQuoteParams(1, testFromToken, testToToken, testAmount, false, testTaker, 50)

	if params.Get("taker") != testTaker {
		t.Errorf("期望 taker=%s, 得到 %s", testTaker, params.Get("taker"))
	}
	if params.Get("slippageBps") != "50" {
		t.Errorf("期望 slippageBps=50, 得到 %s", params.Get("slippageBps"))
	}
}

// --- 错误类型测试 ---

func TestAPIError_错误消息(t *testing.T) {
	err := &APIError{StatusCode: 400, Code: 400, Reason: "Invalid token"}
	msg := err.Error()
	if msg == "" {
		t.Error("期望非空错误消息")
	}
}

func TestRateLimitError_错误消息(t *testing.T) {
	err := &RateLimitError{StatusCode: 429}
	msg := err.Error()
	if msg == "" {
		t.Error("期望非空错误消息")
	}
}

// --- 非标准 JSON 错误响应测试 ---

func TestGetQuote_非JSON错误响应(t *testing.T) {
	server := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	})
	defer server.Close()

	provider := newTestProvider(server.URL)
	_, err := provider.GetQuote(context.Background(), domain.QuoteParams{
		ChainID:   1,
		FromToken: testFromToken,
		ToToken:   testToToken,
		Amount:    testAmount,
		SwapMode:  domain.SwapModeExactIn,
	})
	if err == nil {
		t.Fatal("期望错误, 得到 nil")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("期望 *APIError, 得到 %T: %v", err, err)
	}
	if apiErr.StatusCode != 500 {
		t.Errorf("期望 statusCode=500, 得到 %d", apiErr.StatusCode)
	}
}

// --- slippageBps 转换边界值测试 ---

func TestSlippageBps_边界值转换(t *testing.T) {
	tests := []struct {
		slippagePercent float64
		expectedBps     int
	}{
		{0.0, 0},
		{0.01, 1},
		{0.5, 50},
		{1.0, 100},
		{50.0, 5000},
		{100.0, 10000},
	}

	for _, tt := range tests {
		t.Run(strconv.FormatFloat(tt.slippagePercent, 'f', -1, 64)+"%", func(t *testing.T) {
			result := int(tt.slippagePercent * 100)
			if result != tt.expectedBps {
				t.Errorf("期望 %f%% → %d bps, 得到 %d", tt.slippagePercent, tt.expectedBps, result)
			}
		})
	}
}
