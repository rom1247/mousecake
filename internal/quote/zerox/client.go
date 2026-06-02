package zerox

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	// defaultBaseURL 0x API 默认基础地址。
	defaultBaseURL = "https://api.0x.org"
	// defaultTimeout HTTP 请求默认超时时间。
	defaultTimeout = 10 * time.Second
	// headerAPIKey 0x API Key 请求头名称。
	headerAPIKey = "0x-api-key"
	// headerVersion 0x API 版本请求头名称。
	headerVersion = "0x-version"
	// headerVersionValue 0x API 版本值。
	headerVersionValue = "v2"
)

// Config 0x 供应商配置。
type Config struct {
	// APIKey 0x API 密钥。
	APIKey string
	// BaseURL 0x API 基础地址，默认 "https://api.0x.org"。
	BaseURL string
	// Timeout HTTP 请求超时时间，默认 10 秒。
	Timeout time.Duration
}

// httpClient 封装了 0x API 的 HTTP 请求逻辑，负责认证头注入和响应解析。
type httpClient struct {
	// baseURL 0x API 基础地址。
	baseURL string
	// apiKey 0x API 密钥。
	apiKey string
	// client 底层 HTTP 客户端。
	client *http.Client
}

// newHTTPClient 根据配置创建 HTTP 客户端实例。
// 使用独立的 http.Client 和独立超时，不直接传播用户请求 context。
func newHTTPClient(cfg Config) *httpClient {
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = defaultTimeout
	}
	return &httpClient{
		baseURL: baseURL,
		apiKey:  cfg.APIKey,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

// doGET 执行 GET 请求，path 格式为 "/{chainId}/swap/allowance-holder/price"。
// 使用独立 context 和超时，不直接传播用户请求 context。
func (c *httpClient) doGET(ctx context.Context, path string, params url.Values, result interface{}) error {
	u := c.baseURL + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return fmt.Errorf("构造请求失败: %w", err)
	}

	// 注入认证头
	req.Header.Set(headerAPIKey, c.apiKey)
	req.Header.Set(headerVersion, headerVersionValue)

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应体失败: %w", err)
	}

	// 处理 HTTP 错误
	if resp.StatusCode == http.StatusTooManyRequests {
		return &RateLimitError{
			StatusCode: resp.StatusCode,
			Body:       string(body),
		}
	}
	if resp.StatusCode >= http.StatusBadRequest {
		var apiErr zeroxError
		if jsonErr := json.Unmarshal(body, &apiErr); jsonErr == nil && apiErr.Reason != "" {
			return &APIError{
				StatusCode: resp.StatusCode,
				Reason:     apiErr.Reason,
				Code:       apiErr.Code,
			}
		}
		return &APIError{
			StatusCode: resp.StatusCode,
			Reason:     string(body),
		}
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("解析响应 JSON 失败: %w", err)
	}
	return nil
}

// --- 错误类型定义 ---

// APIError 0x API 业务错误，包含 HTTP 状态码和错误原因。
type APIError struct {
	// StatusCode HTTP 状态码。
	StatusCode int
	// Reason 错误原因描述。
	Reason string
	// Code 0x 错误码。
	Code int
}

// Error 实现 error 接口。
func (e *APIError) Error() string {
	if e.Code > 0 {
		return fmt.Sprintf("0x api 错误: status=%d, code=%d, reason=%s", e.StatusCode, e.Code, e.Reason)
	}
	return fmt.Sprintf("0x api 错误: status=%d, reason=%s", e.StatusCode, e.Reason)
}

// RateLimitError 0x API 限流错误（HTTP 429）。
type RateLimitError struct {
	// StatusCode HTTP 状态码。
	StatusCode int
	// Body 原始响应体。
	Body string
}

// Error 实现 error 接口。
func (e *RateLimitError) Error() string {
	return fmt.Sprintf("0x api 限流: status=%d", e.StatusCode)
}

// buildPriceParams 构建 /price 接口的查询参数。
func buildPriceParams(chainID int, sellToken, buyToken, amount string, isExactOut bool) url.Values {
	params := url.Values{}
	params.Set("chainId", strconv.Itoa(chainID))
	params.Set("sellToken", sellToken)
	params.Set("buyToken", buyToken)

	if isExactOut {
		params.Set("buyAmount", amount)
	} else {
		params.Set("sellAmount", amount)
	}
	return params
}

// buildQuoteParams 构建 /quote 接口的查询参数，在 price 参数基础上增加 taker 和 slippageBps。
func buildQuoteParams(chainID int, sellToken, buyToken, amount string, isExactOut bool, taker string, slippageBps int) url.Values {
	params := buildPriceParams(chainID, sellToken, buyToken, amount, isExactOut)
	params.Set("taker", taker)
	params.Set("slippageBps", strconv.Itoa(slippageBps))
	return params
}
