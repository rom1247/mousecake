// Package zerox 实现 0x Swap API 报价供应商适配器，
// 通过 0x Swap API v2 (AllowanceHolder) 获取代币兑换报价和交易数据。
package zerox

// priceResponse /swap/allowance-holder/price 接口响应结构（指示性报价）。
type priceResponse struct {
	// BuyAmount 买入代币数量（wei）。
	BuyAmount string `json:"buyAmount"`
	// BuyToken 买入代币地址。
	BuyToken string `json:"buyToken"`
	// SellAmount 卖出代币数量（wei）。
	SellAmount string `json:"sellAmount"`
	// SellToken 卖出代币地址。
	SellToken string `json:"sellToken"`
	// Gas 预估 Gas 限制。
	Gas string `json:"gas"`
	// GasPrice Gas 价格（wei）。
	GasPrice string `json:"gasPrice"`
	// MinBuyAmount 最小买入数量（受滑点影响）。
	MinBuyAmount string `json:"minBuyAmount"`
	// Route 路由信息。
	Route routeInfo `json:"route"`
}

// quoteResponse /swap/allowance-holder/quote 接口响应结构（确定性报价，含交易数据）。
type quoteResponse struct {
	// BuyAmount 买入代币数量（wei）。
	BuyAmount string `json:"buyAmount"`
	// BuyToken 买入代币地址。
	BuyToken string `json:"buyToken"`
	// SellAmount 卖出代币数量（wei）。
	SellAmount string `json:"sellAmount"`
	// SellToken 卖出代币地址。
	SellToken string `json:"sellToken"`
	// Gas 预估 Gas 限制。
	Gas string `json:"gas"`
	// GasPrice Gas 价格（wei）。
	GasPrice string `json:"gasPrice"`
	// MinBuyAmount 最小买入数量（受滑点影响）。
	MinBuyAmount string `json:"minBuyAmount"`
	// Route 路由信息。
	Route routeInfo `json:"route"`
	// Transaction 交易数据，仅在 /quote 接口返回。
	Transaction *transactionData `json:"transaction"`
}

// routeInfo 0x 路由信息。
type routeInfo struct {
	// Fills 路由填充详情列表。
	Fills []routeFill `json:"fills"`
	// Tokens 路由中涉及的代币列表。
	Tokens []routeToken `json:"tokens"`
}

// routeFill 路由中的一跳填充信息。
type routeFill struct {
	// From 源代币地址。
	From string `json:"from"`
	// To 目标代币地址。
	To string `json:"to"`
	// Source 流动性来源（如 "Uniswap_V3"）。
	Source string `json:"source"`
	// ProportionBps 占比（基点，10000 = 100%）。
	ProportionBps string `json:"proportionBps"`
}

// routeToken 路由中的代币信息。
type routeToken struct {
	// Address 代币地址。
	Address string `json:"address"`
	// Symbol 代币符号。
	Symbol string `json:"symbol"`
}

// transactionData 0x 返回的交易数据。
type transactionData struct {
	// To 交易目标地址。
	To string `json:"to"`
	// Data 交易调用数据（hex）。
	Data string `json:"data"`
	// Value 交易发送的 ETH 数量（wei）。
	Value string `json:"value"`
	// Gas Gas 限制。
	Gas string `json:"gas"`
	// GasPrice Gas 价格（wei）。
	GasPrice string `json:"gasPrice"`
}

// zeroxError 0x API 错误响应。
type zeroxError struct {
	// Code HTTP 状态码。
	Code int `json:"code"`
	// Reason 错误原因描述。
	Reason string `json:"reason"`
	// ValidationErrors 校验错误详情。
	ValidationErrors []validationError `json:"validationErrors,omitempty"`
}

// validationError 校验错误详情。
type validationError struct {
	// Field 错误字段。
	Field string `json:"field"`
	// Reason 错误原因。
	Reason string `json:"reason"`
}
