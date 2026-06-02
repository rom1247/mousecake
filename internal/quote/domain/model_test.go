package domain

import (
	"strings"
	"testing"
)

func TestQuoteParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  QuoteParams
		wantErr error
	}{
		{
			name: "QuoteParams 字段完整",
			params: QuoteParams{
				ChainID:   1,
				FromToken: "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
				ToToken:   "0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2",
				Amount:    "1000000",
				SwapMode:  SwapModeExactIn,
			},
			wantErr: nil,
		},
		{
			name: "缺少 chain_id",
			params: QuoteParams{
				FromToken: "0xA",
				ToToken:   "0xB",
				Amount:    "1000",
				SwapMode:  SwapModeExactIn,
			},
			wantErr: ErrInvalidChainID,
		},
		{
			name: "缺少 from_token",
			params: QuoteParams{
				ChainID:  1,
				ToToken:  "0xB",
				Amount:   "1000",
				SwapMode: SwapModeExactIn,
			},
			wantErr: ErrFromTokenEmpty,
		},
		{
			name: "缺少 to_token",
			params: QuoteParams{
				ChainID:   1,
				FromToken: "0xA",
				Amount:    "1000",
				SwapMode:  SwapModeExactIn,
			},
			wantErr: ErrToTokenEmpty,
		},
		{
			name: "缺少 amount",
			params: QuoteParams{
				ChainID:   1,
				FromToken: "0xA",
				ToToken:   "0xB",
				SwapMode:  SwapModeExactIn,
			},
			wantErr: ErrAmountEmpty,
		},
		{
			name: "无效 swap_mode",
			params: QuoteParams{
				ChainID:   1,
				FromToken: "0xA",
				ToToken:   "0xB",
				Amount:    "1000",
				SwapMode:  "invalid",
			},
			wantErr: ErrInvalidSwapMode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("期望错误 %v, 得到 nil", tt.wantErr)
				}
				if err != tt.wantErr {
					t.Fatalf("期望错误 %v, 得到 %v", tt.wantErr, err)
				}
			} else if err != nil {
				t.Fatalf("不期望错误, 得到 %v", err)
			}
		})
	}
}

func TestSwapParams_Validate(t *testing.T) {
	tests := []struct {
		name    string
		params  SwapParams
		wantErr error
	}{
		{
			name: "SwapParams 包含钱包地址和滑点",
			params: SwapParams{
				QuoteParams: QuoteParams{
					ChainID:   1,
					FromToken: "0xA",
					ToToken:   "0xB",
					Amount:    "1000",
					SwapMode:  SwapModeExactIn,
				},
				SlippagePercent:   0.5,
				UserWalletAddress: "0xUser12345678901234567890123456789012345678",
			},
			wantErr: nil,
		},
		{
			name: "缺少钱包地址",
			params: SwapParams{
				QuoteParams: QuoteParams{
					ChainID:   1,
					FromToken: "0xA",
					ToToken:   "0xB",
					Amount:    "1000",
					SwapMode:  SwapModeExactIn,
				},
				SlippagePercent:   0.5,
				UserWalletAddress: "",
			},
			wantErr: ErrWalletAddressEmpty,
		},
		{
			name: "滑点为负数",
			params: SwapParams{
				QuoteParams: QuoteParams{
					ChainID:   1,
					FromToken: "0xA",
					ToToken:   "0xB",
					Amount:    "1000",
					SwapMode:  SwapModeExactIn,
				},
				SlippagePercent:   -1,
				UserWalletAddress: "0xUser",
			},
			wantErr: ErrInvalidSlippage,
		},
		{
			name: "滑点超过 100",
			params: SwapParams{
				QuoteParams: QuoteParams{
					ChainID:   1,
					FromToken: "0xA",
					ToToken:   "0xB",
					Amount:    "1000",
					SwapMode:  SwapModeExactIn,
				},
				SlippagePercent:   101,
				UserWalletAddress: "0xUser",
			},
			wantErr: ErrInvalidSlippage,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.params.Validate()
			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("期望错误 %v, 得到 nil", tt.wantErr)
				}
				if err != tt.wantErr {
					t.Fatalf("期望错误 %v, 得到 %v", tt.wantErr, err)
				}
			} else if err != nil {
				t.Fatalf("不期望错误, 得到 %v", err)
			}
		})
	}
}

func TestSwapResult_ContainsTxData(t *testing.T) {
	result := SwapResult{
		QuoteResult: QuoteResult{
			Provider:   "okx",
			FromAmount: "1000",
			ToAmount:   "2000",
		},
		TxData: TxData{
			To:       "0xContract",
			Data:     "0xabcd",
			Value:    "0",
			Gas:      "21000",
			GasPrice: "1000000000",
		},
	}
	if result.TxData.To != "0xContract" {
		t.Errorf("期望 TxData.To=0xContract, 得到 %s", result.TxData.To)
	}
	if result.TxData.Data != "0xabcd" {
		t.Errorf("期望 TxData.Data=0xabcd, 得到 %s", result.TxData.Data)
	}
}

func TestCacheKey(t *testing.T) {
	params := QuoteParams{
		ChainID:   1,
		FromToken: "0xA",
		ToToken:   "0xB",
		Amount:    "1000",
		SwapMode:  SwapModeExactIn,
	}
	key := CacheKey("okx", params)
	expected := "quote:okx:1:0xA:0xB:1000:exactIn"
	if key != expected {
		t.Errorf("期望 key=%s, 得到 %s", expected, key)
	}

	// 不同 amount 应产生不同 key
	params2 := QuoteParams{
		ChainID:   1,
		FromToken: "0xA",
		ToToken:   "0xB",
		Amount:    "2000",
		SwapMode:  SwapModeExactIn,
	}
	key2 := CacheKey("okx", params2)
	if key2 == key {
		t.Error("不同 amount 应产生不同的 key")
	}

	// 不同 swapMode 应产生不同 key（回归：CacheKey 曾缺少 swapMode）
	params3 := QuoteParams{
		ChainID:   1,
		FromToken: "0xA",
		ToToken:   "0xB",
		Amount:    "1000",
		SwapMode:  SwapModeExactOut,
	}
	key3 := CacheKey("okx", params3)
	if key3 == key {
		t.Error("不同 swapMode 应产生不同的 key")
	}
}

func TestTokenInfo(t *testing.T) {
	info := TokenInfo{
		Address:  "0xA",
		Decimals: 18,
		Symbol:   "USDT",
	}
	if info.Symbol != "USDT" {
		t.Errorf("期望 Symbol=USDT, 得到 %s", info.Symbol)
	}
}

func TestRouteHop(t *testing.T) {
	hop := RouteHop{
		PoolAddress:    "0xPool",
		TokenAddresses: []string{"0xA", "0xB"},
	}
	if len(hop.TokenAddresses) != 2 {
		t.Errorf("期望 2 个 token, 得到 %d", len(hop.TokenAddresses))
	}
}

func TestSwapResult_Fields(t *testing.T) {
	result := SwapResult{
		QuoteResult: QuoteResult{
			Provider:   "zerox",
			FromToken:  TokenInfo{Address: "0xA", Decimals: 18, Symbol: "WETH"},
			ToToken:    TokenInfo{Address: "0xB", Decimals: 6, Symbol: "USDT"},
			FromAmount: "1000000000000000000",
			ToAmount:   "3000000000",
			Route: []RouteHop{
				{PoolAddress: "0xPool1", TokenAddresses: []string{"0xA", "0xB"}},
			},
		},
		TxData: TxData{
			To:       "0xRouter",
			Data:     "0xdeadbeef",
			Value:    "0",
			Gas:      "150000",
			GasPrice: "2000000000",
		},
	}

	if result.Provider != "zerox" {
		t.Errorf("期望 Provider=zerox, 得到 %s", result.Provider)
	}
	if result.FromToken.Symbol != "WETH" {
		t.Errorf("期望 FromToken.Symbol=WETH, 得到 %s", result.FromToken.Symbol)
	}
	if len(result.Route) != 1 {
		t.Errorf("期望 1 个 route hop, 得到 %d", len(result.Route))
	}
	if !strings.HasPrefix(result.TxData.To, "0x") {
		t.Errorf("期望 TxData.To 以 0x 开头, 得到 %s", result.TxData.To)
	}
}
