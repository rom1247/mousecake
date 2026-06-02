package domain

import (
	"errors"
	"fmt"
)

// 领域模型哨兵错误
var (
	// ErrInvalidChainID 无效的链 ID。
	ErrInvalidChainID = errors.New("quote: 无效的链 ID")
	// ErrAmountEmpty amount 不能为空。
	ErrAmountEmpty = errors.New("quote: amount 不能为空")
	// ErrWalletAddressEmpty 钱包地址不能为空。
	ErrWalletAddressEmpty = errors.New("quote: 钱包地址不能为空")
	// ErrInvalidSlippage 无效的滑点设置。
	ErrInvalidSlippage = errors.New("quote: 无效的滑点设置")
)

// TokenInfo 代币信息。
type TokenInfo struct {
	// Address 代币合约地址。
	Address string
	// Decimals 代币精度。
	Decimals int
	// Symbol 代币符号。
	Symbol string
}

// RouteHop 路由跳信息，表示兑换路径中的一跳。
type RouteHop struct {
	// PoolAddress 流动性池地址。
	PoolAddress string
	// TokenAddresses 本跳涉及的代币地址列表。
	TokenAddresses []string
}

// QuoteParams 统一的报价请求参数。
type QuoteParams struct {
	// ChainID 链 ID。
	ChainID int
	// FromToken 源代币地址。
	FromToken string
	// ToToken 目标代币地址。
	ToToken string
	// Amount 兑换数量（wei）。
	Amount string
	// SwapMode 兑换模式。
	SwapMode SwapMode
}

// Validate 校验 QuoteParams 必填字段。
func (p QuoteParams) Validate() error {
	if p.ChainID <= 0 {
		return ErrInvalidChainID
	}
	if p.FromToken == "" {
		return ErrFromTokenEmpty
	}
	if p.ToToken == "" {
		return ErrToTokenEmpty
	}
	if p.Amount == "" {
		return ErrAmountEmpty
	}
	if !validSwapModes[p.SwapMode] {
		return ErrInvalidSwapMode
	}
	return nil
}

// QuoteResult 统一的报价结果。
type QuoteResult struct {
	// Provider 供应商名称。
	Provider string
	// FromToken 源代币信息。
	FromToken TokenInfo
	// ToToken 目标代币信息。
	ToToken TokenInfo
	// FromAmount 源代币数量。
	FromAmount string
	// ToAmount 目标代币数量。
	ToAmount string
	// PriceImpact 价格影响。
	PriceImpact string
	// EstimatedGas 预估 Gas。
	EstimatedGas string
	// Route 兑换路由。
	Route []RouteHop
}

// SwapParams 统一的交易数据请求参数，继承 QuoteParams 并增加滑点和钱包地址。
type SwapParams struct {
	QuoteParams
	// SlippagePercent 滑点百分比（0-100）。
	SlippagePercent float64
	// UserWalletAddress 用户钱包地址。
	UserWalletAddress string
}

// Validate 校验 SwapParams 所有必填字段。
func (p SwapParams) Validate() error {
	if err := p.QuoteParams.Validate(); err != nil {
		return err
	}
	if p.UserWalletAddress == "" {
		return ErrWalletAddressEmpty
	}
	if p.SlippagePercent < 0 || p.SlippagePercent > 100 {
		return ErrInvalidSlippage
	}
	return nil
}

// TxData 交易数据。
type TxData struct {
	// To 交易目标地址。
	To string
	// Data 交易调用数据（hex）。
	Data string
	// Value 交易发送的 ETH 数量（wei）。
	Value string
	// Gas Gas 限制。
	Gas string
	// GasPrice Gas 价格（wei）。
	GasPrice string
}

// SwapResult 统一的交易数据结果，继承 QuoteResult 并增加 TxData。
type SwapResult struct {
	QuoteResult
	// ID SwapRecord 的 ID，前端提交交易哈希时使用。
	ID int64
	// TxData 交易数据。
	TxData TxData
}

// CacheKey 生成缓存 Key，格式为 quote:{provider}:{chain_id}:{from_token}:{to_token}:{amount}:{swap_mode}。
func CacheKey(provider string, params QuoteParams) string {
	return fmt.Sprintf("quote:%s:%d:%s:%s:%s:%s", provider, params.ChainID, params.FromToken, params.ToToken, params.Amount, params.SwapMode)
}
