package okx

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/mousecake-go/mousecake-go/internal/quote/domain"
)

// Provider 实现 OKX DEX 聚合器报价供应商。
type Provider struct {
	client *httpClient
}

// NewProvider 根据 Config 创建 OKX 供应商实例。
func NewProvider(cfg Config) *Provider {
	return &Provider{
		client: newHTTPClient(cfg),
	}
}

// Name 返回供应商名称 "okx"。
func (p *Provider) Name() string {
	return "okx"
}

// GetQuote 调用 OKX DEX 聚合器获取报价数据。
func (p *Provider) GetQuote(ctx context.Context, params domain.QuoteParams) (*domain.QuoteResult, error) {
	queryParams := buildQuoteParams(
		params.ChainID,
		params.FromToken,
		params.ToToken,
		params.Amount,
		string(params.SwapMode),
	)

	okxResp, err := p.client.get(ctx, quotePath, queryParams)
	if err != nil {
		return nil, fmt.Errorf("okx quote 请求失败: %w", err)
	}

	if len(okxResp.Data) == 0 {
		return nil, fmt.Errorf("okx quote 响应数据为空")
	}

	var data quoteData
	if err := json.Unmarshal(okxResp.Data[0], &data); err != nil {
		return nil, fmt.Errorf("okx quote 解析 data 失败: %w", err)
	}

	return mapQuoteDataToResult(&data), nil
}

// GetSwap 调用 OKX DEX 聚合器获取交易数据。
func (p *Provider) GetSwap(ctx context.Context, params domain.SwapParams) (*domain.SwapResult, error) {
	queryParams := buildSwapParams(
		params.ChainID,
		params.FromToken,
		params.ToToken,
		params.Amount,
		string(params.SwapMode),
		params.SlippagePercent,
		params.UserWalletAddress,
	)

	okxResp, err := p.client.get(ctx, swapPath, queryParams)
	if err != nil {
		return nil, fmt.Errorf("okx swap 请求失败: %w", err)
	}

	if len(okxResp.Data) == 0 {
		return nil, fmt.Errorf("okx swap 响应数据为空")
	}

	var data swapData
	if err := json.Unmarshal(okxResp.Data[0], &data); err != nil {
		return nil, fmt.Errorf("okx swap 解析 data 失败: %w", err)
	}

	return mapSwapDataToResult(&data), nil
}

// mapQuoteDataToResult 将 OKX quote 响应数据转换为统一的 QuoteResult。
func mapQuoteDataToResult(data *quoteData) *domain.QuoteResult {
	fromDecimals, _ := strconv.Atoi(data.FromToken.Decimal)
	toDecimals, _ := strconv.Atoi(data.ToToken.Decimal)

	return &domain.QuoteResult{
		Provider: "okx",
		FromToken: domain.TokenInfo{
			Address:  data.FromToken.TokenContractAddress,
			Decimals: fromDecimals,
			Symbol:   data.FromToken.TokenSymbol,
		},
		ToToken: domain.TokenInfo{
			Address:  data.ToToken.TokenContractAddress,
			Decimals: toDecimals,
			Symbol:   data.ToToken.TokenSymbol,
		},
		FromAmount:   data.FromTokenAmount,
		ToAmount:     data.ToTokenAmount,
		PriceImpact:  data.PriceImpactPercent,
		EstimatedGas: data.EstimateGasFee,
	}
}

// mapSwapDataToResult 将 OKX swap 响应数据转换为统一的 SwapResult。
func mapSwapDataToResult(data *swapData) *domain.SwapResult {
	rr := &data.RouterResult
	fromDecimals, _ := strconv.Atoi(rr.FromToken.Decimal)
	toDecimals, _ := strconv.Atoi(rr.ToToken.Decimal)

	return &domain.SwapResult{
		QuoteResult: domain.QuoteResult{
			Provider: "okx",
			FromToken: domain.TokenInfo{
				Address:  rr.FromToken.TokenContractAddress,
				Decimals: fromDecimals,
				Symbol:   rr.FromToken.TokenSymbol,
			},
			ToToken: domain.TokenInfo{
				Address:  rr.ToToken.TokenContractAddress,
				Decimals: toDecimals,
				Symbol:   rr.ToToken.TokenSymbol,
			},
			FromAmount:   rr.FromTokenAmount,
			ToAmount:     rr.ToTokenAmount,
			PriceImpact:  rr.PriceImpactPercent,
			EstimatedGas: rr.EstimateGasFee,
		},
		TxData: domain.TxData{
			To:       data.Tx.To,
			Data:     data.Tx.Data,
			Value:    data.Tx.Value,
			Gas:      data.Tx.Gas,
			GasPrice: data.Tx.GasPrice,
		},
	}
}
