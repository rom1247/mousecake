package zerox

import (
	"context"
	"fmt"

	"github.com/mousecake-go/mousecake-go/internal/quote/domain"
)

const (
	// providerName 供应商名称。
	providerName = "zerox"
	// pricePath 0x /price 接口路径模板，{chainId} 将被替换为实际链 ID。
	pricePath = "/swap/allowance-holder/price"
	// quotePath 0x /quote 接口路径模板。
	quotePath = "/swap/allowance-holder/quote"
)

// Provider 0x Swap API 报价供应商，实现 domain.QuoteProvider 接口。
type Provider struct {
	// client HTTP 客户端。
	client *httpClient
}

// NewProvider 创建 0x 供应商实例。
func NewProvider(cfg Config) *Provider {
	return &Provider{
		client: newHTTPClient(cfg),
	}
}

// Name 返回供应商名称 "zerox"。
func (p *Provider) Name() string {
	return providerName
}

// GetQuote 获取指示性报价，调用 0x /swap/allowance-holder/price 接口。
// exactIn 模式使用 sellAmount，exactOut 模式使用 buyAmount。
func (p *Provider) GetQuote(ctx context.Context, params domain.QuoteParams) (*domain.QuoteResult, error) {
	isExactOut := params.SwapMode == domain.SwapModeExactOut

	queryParams := buildPriceParams(
		params.ChainID,
		params.FromToken,
		params.ToToken,
		params.Amount,
		isExactOut,
	)

	// 构造带 chainId 的 URL path
	path := fmt.Sprintf("/%d%s", params.ChainID, pricePath)

	var resp priceResponse
	if err := p.client.doGET(ctx, path, queryParams, &resp); err != nil {
		return nil, fmt.Errorf("0x 获取报价失败: %w", err)
	}

	result := &domain.QuoteResult{
		Provider: providerName,
		FromToken: domain.TokenInfo{
			Address: resp.SellToken,
		},
		ToToken: domain.TokenInfo{
			Address: resp.BuyToken,
		},
		FromAmount:   resp.SellAmount,
		ToAmount:     resp.BuyAmount,
		EstimatedGas: resp.Gas,
		Route:        convertRoute(resp.Route),
	}

	return result, nil
}

// GetSwap 获取确定性报价和交易数据，调用 0x /swap/allowance-holder/quote 接口。
// slippagePercent 转换为 slippageBps（1% = 100bps）。
func (p *Provider) GetSwap(ctx context.Context, params domain.SwapParams) (*domain.SwapResult, error) {
	isExactOut := params.SwapMode == domain.SwapModeExactOut
	slippageBps := int(params.SlippagePercent * 100)

	queryParams := buildQuoteParams(
		params.ChainID,
		params.FromToken,
		params.ToToken,
		params.Amount,
		isExactOut,
		params.UserWalletAddress,
		slippageBps,
	)

	path := fmt.Sprintf("/%d%s", params.ChainID, quotePath)

	var resp quoteResponse
	if err := p.client.doGET(ctx, path, queryParams, &resp); err != nil {
		return nil, fmt.Errorf("0x 获取 swap 失败: %w", err)
	}

	result := &domain.SwapResult{
		QuoteResult: domain.QuoteResult{
			Provider: providerName,
			FromToken: domain.TokenInfo{
				Address: resp.SellToken,
			},
			ToToken: domain.TokenInfo{
				Address: resp.BuyToken,
			},
			FromAmount:   resp.SellAmount,
			ToAmount:     resp.BuyAmount,
			EstimatedGas: resp.Gas,
			Route:        convertRoute(resp.Route),
		},
	}

	if resp.Transaction != nil {
		result.TxData = domain.TxData{
			To:       resp.Transaction.To,
			Data:     resp.Transaction.Data,
			Value:    resp.Transaction.Value,
			Gas:      resp.Transaction.Gas,
			GasPrice: resp.Transaction.GasPrice,
		}
	}

	return result, nil
}

// convertRoute 将 0x 路由信息转换为领域 RouteHop 列表。
// 0x 的 route.fills 描述了路由中每一跳，每跳的 source 对应流动性池名称。
func convertRoute(route routeInfo) []domain.RouteHop {
	if len(route.Fills) == 0 {
		return nil
	}

	hops := make([]domain.RouteHop, 0, len(route.Fills))
	for _, fill := range route.Fills {
		hops = append(hops, domain.RouteHop{
			PoolAddress: fill.Source,
			TokenAddresses: []string{
				fill.From,
				fill.To,
			},
		})
	}
	return hops
}
