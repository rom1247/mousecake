// Package okx 实现 OKX DEX 聚合器报价供应商适配器，
// 通过 OKX DEX API v6 获取代币兑换报价和交易数据。
package okx

import "encoding/json"

// okxResponse OKX API 通用响应包装。
// code 字段为 string 类型，"0" 表示成功，非 "0" 表示业务错误。
type okxResponse struct {
	// Code 响应码，"0" 表示成功。
	Code string `json:"code"`
	// Msg 响应消息，成功时为空字符串。
	Msg string `json:"msg"`
	// Data 响应数据数组。
	Data []json.RawMessage `json:"data"`
}

// quoteData quote 接口响应中的 data[0] 结构。
type quoteData struct {
	// ChainIndex 链的唯一标识。
	ChainIndex string `json:"chainIndex"`
	// FromTokenAmount 卖出代币数量。
	FromTokenAmount string `json:"fromTokenAmount"`
	// ToTokenAmount 买入代币数量。
	ToTokenAmount string `json:"toTokenAmount"`
	// PriceImpactPercent 价格影响百分比。
	PriceImpactPercent string `json:"priceImpactPercent"`
	// EstimateGasFee 预估 Gas 费用（wei）。
	EstimateGasFee string `json:"estimateGasFee"`
	// Router 兑换路径描述。
	Router string `json:"router"`
	// FromToken 卖出代币信息。
	FromToken tokenData `json:"fromToken"`
	// ToToken 买入代币信息。
	ToToken tokenData `json:"toToken"`
	// SwapMode 交易模式（exactIn 或 exactOut）。
	SwapMode string `json:"swapMode"`
}

// swapData swap 接口响应中的 data[0] 结构，包含报价和交易数据。
type swapData struct {
	// RouterResult 报价路径数据对象。
	RouterResult routerResult `json:"routerResult"`
	// Tx 交易数据。
	Tx txData `json:"tx"`
}

// routerResult swap 响应中的报价结果。
type routerResult struct {
	// ChainIndex 链的唯一标识。
	ChainIndex string `json:"chainIndex"`
	// FromTokenAmount 卖出代币数量。
	FromTokenAmount string `json:"fromTokenAmount"`
	// ToTokenAmount 买入代币数量。
	ToTokenAmount string `json:"toTokenAmount"`
	// PriceImpactPercent 价格影响百分比。
	PriceImpactPercent string `json:"priceImpactPercent"`
	// EstimateGasFee 预估 Gas 费用（wei）。
	EstimateGasFee string `json:"estimateGasFee"`
	// Router 兑换路径描述。
	Router string `json:"router"`
	// FromToken 卖出代币信息。
	FromToken tokenData `json:"fromToken"`
	// ToToken 买入代币信息。
	ToToken tokenData `json:"toToken"`
	// SwapMode 交易模式。
	SwapMode string `json:"swapMode"`
}

// txData swap 响应中的交易数据。
type txData struct {
	// From 用户钱包地址。
	From string `json:"from"`
	// To OKX DEX Router 合约地址。
	To string `json:"to"`
	// Data 交易调用数据（hex）。
	Data string `json:"data"`
	// Value 交易发送的主链币数量（wei）。
	Value string `json:"value"`
	// Gas Gas 限制。
	Gas string `json:"gas"`
	// GasPrice Gas 价格（wei）。
	GasPrice string `json:"gasPrice"`
}

// tokenData OKX 响应中的代币信息。
type tokenData struct {
	// TokenContractAddress 代币合约地址。
	TokenContractAddress string `json:"tokenContractAddress"`
	// TokenSymbol 代币符号。
	TokenSymbol string `json:"tokenSymbol"`
	// Decimal 代币精度。
	Decimal string `json:"decimal"`
}
