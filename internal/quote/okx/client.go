package okx

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	// defaultBaseURL OKX API 默认基础 URL。
	defaultBaseURL = "https://www.okx.com"
	// defaultTimeout HTTP 请求默认超时时间。
	defaultTimeout = 10 * time.Second
	// quotePath 获取报价的 API 路径。
	quotePath = "/api/v6/dex/aggregator/quote"
	// swapPath 获取交易数据的 API 路径。
	swapPath = "/api/v6/dex/aggregator/swap"
)

// Config OKX 供应商配置。
type Config struct {
	// APIKey OKX API 密钥。
	APIKey string
	// SecretKey OKX API 密钥对应的 Secret。
	SecretKey string
	// Passphrase OKX API 密钥对应的密码短语。
	Passphrase string
	// BaseURL OKX API 基础 URL，默认为 "https://www.okx.com"。
	BaseURL string
	// Timeout HTTP 请求超时时间，默认为 10 秒。
	Timeout time.Duration
}

// httpClient 封装 HTTP 客户端和签名逻辑。
type httpClient struct {
	cfg     Config
	client  *http.Client
	baseURL string
}

// newHTTPClient 根据 Config 创建 HTTP 客户端。
func newHTTPClient(cfg Config) *httpClient {
	if cfg.BaseURL == "" {
		cfg.BaseURL = defaultBaseURL
	}
	if cfg.Timeout <= 0 {
		cfg.Timeout = defaultTimeout
	}
	return &httpClient{
		cfg:     cfg,
		client:  &http.Client{Timeout: cfg.Timeout},
		baseURL: cfg.BaseURL,
	}
}

// sign 生成 OKX API 签名。
// stringToSign = timestamp + method + requestPath + ("?" + queryString if queryString not empty)
// signature = Base64(HMAC-SHA256(stringToSign, secretKey))
func sign(timestamp, method, requestPath, queryString, secretKey string) string {
	var message string
	if queryString != "" {
		message = timestamp + method + requestPath + "?" + queryString
	} else {
		message = timestamp + method + requestPath
	}
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(message))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

// formatTimestamp 生成 OKX API 要求的 ISO 8601 毫秒精度时间戳。
func formatTimestamp(t time.Time) string {
	return t.UTC().Format("2006-01-02T15:04:05.000Z")
}

// buildAuthHeaders 构建带签名的认证 HTTP 头。
func (c *httpClient) buildAuthHeaders(requestPath, queryString string) http.Header {
	timestamp := formatTimestamp(time.Now())
	signature := sign(timestamp, "GET", requestPath, queryString, c.cfg.SecretKey)

	header := http.Header{}
	header.Set("OK-ACCESS-KEY", c.cfg.APIKey)
	header.Set("OK-ACCESS-SIGN", signature)
	header.Set("OK-ACCESS-TIMESTAMP", timestamp)
	header.Set("OK-ACCESS-PASSPHRASE", c.cfg.Passphrase)
	header.Set("Content-Type", "application/json")
	return header
}

// get 发送带签名的 GET 请求。
func (c *httpClient) get(ctx context.Context, requestPath string, params url.Values) (*okxResponse, error) {
	queryString := params.Encode()
	fullURL := c.baseURL + requestPath
	if queryString != "" {
		fullURL += "?" + queryString
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header = c.buildAuthHeaders(requestPath, queryString)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("okx api 返回 HTTP %d: %s", resp.StatusCode, string(body))
	}

	var okxResp okxResponse
	if err := json.Unmarshal(body, &okxResp); err != nil {
		return nil, fmt.Errorf("解析响应 JSON 失败: %w", err)
	}

	if okxResp.Code != "0" {
		return nil, fmt.Errorf("okx api 业务错误 code=%s msg=%s", okxResp.Code, okxResp.Msg)
	}

	return &okxResp, nil
}

// buildQuoteParams 构造 quote 接口的查询参数。
func buildQuoteParams(chainID int, fromToken, toToken, amount, swapMode string) url.Values {
	params := url.Values{}
	params.Set("chainIndex", strconv.Itoa(chainID))
	params.Set("fromTokenAddress", fromToken)
	params.Set("toTokenAddress", toToken)
	params.Set("amount", amount)
	params.Set("swapMode", swapMode)
	return params
}

// buildSwapParams 构造 swap 接口的查询参数。
func buildSwapParams(chainID int, fromToken, toToken, amount, swapMode string, slippagePercent float64, userWalletAddress string) url.Values {
	params := url.Values{}
	params.Set("chainIndex", strconv.Itoa(chainID))
	params.Set("fromTokenAddress", fromToken)
	params.Set("toTokenAddress", toToken)
	params.Set("amount", amount)
	params.Set("swapMode", swapMode)
	params.Set("slippagePercent", strconv.FormatFloat(slippagePercent, 'f', -1, 64))
	params.Set("userWalletAddress", userWalletAddress)
	return params
}
