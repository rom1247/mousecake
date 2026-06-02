package quote

import (
	"errors"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/mousecake-go/mousecake-go/internal/quote/domain"
	"github.com/mousecake-go/mousecake-go/internal/shared/response"
)

var (
	evmAddressRegex = regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)
	amountRegex     = regexp.MustCompile(`^[0-9]+$`)
	txHashRegex     = regexp.MustCompile(`^0x[0-9a-fA-F]{64}$`)
)

// Handler 注册 quote 模块的 HTTP 路由。
type Handler struct {
	svc *QuoteService
	log *slog.Logger
}

// NewHandler 创建 Handler。
func NewHandler(svc *QuoteService) *Handler {
	return &Handler{
		svc: svc,
		log: slog.Default().With("module", "quote", "layer", "handler"),
	}
}

// RegisterRoutes 注册所有路由到 Gin RouterGroup。
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/providers", h.getProviders)
	rg.GET("/quote", h.getQuote)
	rg.POST("/swap", h.getSwap)
	rg.POST("/swap/:id/submit", h.submitSwap)
}

// getProviders 获取可用供应商列表。
//
//	@Summary      获取供应商列表
//	@Description  获取当前可用的报价供应商列表
//	@Tags         报价
//	@Produce      json
//	@Success      200  {object}  response.Response{data=[]string}
//	@Failure      500  {object}  response.Response
//	@Router       /providers [get]
func (h *Handler) getProviders(c *gin.Context) {
	providers := h.svc.GetProviders()
	response.Success(c, providers)
}

// getQuote 获取报价。
//
//	@Summary      获取报价
//	@Description  根据指定参数从供应商获取代币兑换报价
//	@Tags         报价
//	@Accept       json
//	@Produce      json
//	@Param        provider    query     string  true  "供应商名称"       Enums(okx, zerox)
//	@Param        chain_id    query     int     true  "链 ID"           example(1)
//	@Param        from_token  query     string  true  "源代币地址"       example(0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE)
//	@Param        to_token    query     string  true  "目标代币地址"     example(0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48)
//	@Param        amount      query     string  true  "兑换数量"         example(1000000000000000000)
//	@Param        swap_mode   query     string  true  "兑换模式"         Enums(exactIn, exactOut)
//	@Success      200  {object}  response.Response
//	@Failure      400  {object}  response.Response
//	@Failure      500  {object}  response.Response
//	@Router       /quote [get]
func (h *Handler) getQuote(c *gin.Context) {
	provider := c.Query("provider")
	chainIDStr := c.Query("chain_id")
	fromToken := c.Query("from_token")
	toToken := c.Query("to_token")
	amount := c.Query("amount")
	swapMode := c.Query("swap_mode")

	// 参数校验
	if provider == "" || chainIDStr == "" || fromToken == "" || toToken == "" || amount == "" || swapMode == "" {
		response.Error(c, http.StatusBadRequest, ErrCodeInvalidParam, "缺少必填参数")
		return
	}

	chainID, err := strconv.Atoi(chainIDStr)
	if err != nil || chainID <= 0 {
		response.Error(c, http.StatusBadRequest, ErrCodeInvalidParam, "无效的 chain_id")
		return
	}

	if !isValidEVMAddress(fromToken) {
		response.Error(c, http.StatusBadRequest, ErrCodeInvalidParam, "无效的 from_token 地址格式")
		return
	}

	if !isValidEVMAddress(toToken) {
		response.Error(c, http.StatusBadRequest, ErrCodeInvalidParam, "无效的 to_token 地址格式")
		return
	}

	if !isValidAmount(amount) {
		response.Error(c, http.StatusBadRequest, ErrCodeInvalidParam, "无效的 amount 格式")
		return
	}

	if swapMode != string(domain.SwapModeExactIn) && swapMode != string(domain.SwapModeExactOut) {
		response.Error(c, http.StatusBadRequest, ErrCodeInvalidParam, "无效的 swap_mode")
		return
	}

	params := domain.QuoteParams{
		ChainID:   chainID,
		FromToken: fromToken,
		ToToken:   toToken,
		Amount:    amount,
		SwapMode:  domain.SwapMode(swapMode),
	}

	result, err := h.svc.GetQuote(c.Request.Context(), provider, params)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	response.Success(c, result)
}

// getSwapRequest 获取交易数据请求。
type getSwapRequest struct {
	Provider        string  `json:"provider" binding:"required" example:"okx"`
	ChainID         int     `json:"chain_id" binding:"required" example:"1"`
	FromToken       string  `json:"from_token" binding:"required" example:"0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"`
	ToToken         string  `json:"to_token" binding:"required" example:"0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"`
	Amount          string  `json:"amount" binding:"required" example:"1000000000000000000"`
	SwapMode        string  `json:"swap_mode" binding:"required" example:"exactIn"`
	SlippagePercent float64 `json:"slippage_percent" example:"0.5"`
	WalletAddress   string  `json:"user_wallet_address" binding:"required" example:"0x1234567890abcdef1234567890abcdef12345678"`
}

// getSwap 获取交易数据。
//
//	@Summary      获取交易数据
//	@Description  根据报价参数获取链上交易所需的交易数据
//	@Tags         报价
//	@Accept       json
//	@Produce      json
//	@Param        body  body      getSwapRequest  true  "获取交易数据请求"
//	@Success      200   {object}  response.Response
//	@Failure      400   {object}  response.Response
//	@Failure      500   {object}  response.Response
//	@Router       /swap [post]
func (h *Handler) getSwap(c *gin.Context) {
	var req getSwapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, ErrCodeInvalidParam, "缺少必填参数")
		return
	}

	if req.ChainID <= 0 {
		response.Error(c, http.StatusBadRequest, ErrCodeInvalidParam, "无效的 chain_id")
		return
	}

	if !isValidEVMAddress(req.FromToken) {
		response.Error(c, http.StatusBadRequest, ErrCodeInvalidParam, "无效的 from_token 地址格式")
		return
	}

	if !isValidEVMAddress(req.ToToken) {
		response.Error(c, http.StatusBadRequest, ErrCodeInvalidParam, "无效的 to_token 地址格式")
		return
	}

	if !isValidAmount(req.Amount) {
		response.Error(c, http.StatusBadRequest, ErrCodeInvalidParam, "无效的 amount 格式")
		return
	}

	if req.SwapMode != string(domain.SwapModeExactIn) && req.SwapMode != string(domain.SwapModeExactOut) {
		response.Error(c, http.StatusBadRequest, ErrCodeInvalidParam, "无效的 swap_mode")
		return
	}

	if !isValidEVMAddress(req.WalletAddress) {
		response.Error(c, http.StatusBadRequest, ErrCodeInvalidParam, "无效的钱包地址格式")
		return
	}

	slippagePercent := req.SlippagePercent
	if slippagePercent == 0 {
		slippagePercent = 0.5 // 默认滑点
	}
	if slippagePercent < 0 || slippagePercent > 100 {
		response.Error(c, http.StatusBadRequest, ErrCodeInvalidParam, "slippage_percent 超出范围 (0-100)")
		return
	}

	params := domain.SwapParams{
		QuoteParams: domain.QuoteParams{
			ChainID:   req.ChainID,
			FromToken: req.FromToken,
			ToToken:   req.ToToken,
			Amount:    req.Amount,
			SwapMode:  domain.SwapMode(req.SwapMode),
		},
		SlippagePercent:   slippagePercent,
		UserWalletAddress: req.WalletAddress,
	}

	result, err := h.svc.GetSwap(c.Request.Context(), req.Provider, params)
	if err != nil {
		h.handleServiceError(c, err)
		return
	}

	response.Success(c, result)
}

// submitSwapRequest 提交交易哈希的请求。
type submitSwapRequest struct {
	TxHash string `json:"tx_hash" binding:"required" example:"0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"`
}

// submitSwap 提交交易哈希。
//
//	@Summary      提交交易哈希
//	@Description  提交已广播的链上交易哈希，完成 swap 记录
//	@Tags         报价
//	@Accept       json
//	@Produce      json
//	@Param        id     path      int                true  "Swap 记录 ID"
//	@Param        body   body      submitSwapRequest  true  "交易哈希"
//	@Success      200    {object}  response.Response
//	@Failure      400    {object}  response.Response
//	@Failure      404    {object}  response.Response
//	@Failure      409    {object}  response.Response
//	@Failure      500    {object}  response.Response
//	@Router       /swap/{id}/submit [post]
func (h *Handler) submitSwap(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, ErrCodeInvalidParam, "无效的 swap ID")
		return
	}

	var req submitSwapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, ErrCodeInvalidParam, "缺少 tx_hash 参数")
		return
	}

	if !isValidTxHash(req.TxHash) {
		response.Error(c, http.StatusBadRequest, ErrCodeInvalidParam, "无效的 tx_hash 格式")
		return
	}

	err = h.svc.SubmitSwap(c.Request.Context(), id, req.TxHash)
	if err != nil {
		if errors.Is(err, domain.ErrSwapRecordNotFound) {
			response.Error(c, http.StatusNotFound, ErrCodeSwapNotFound, "swap 记录不存在")
			return
		}
		if errors.Is(err, domain.ErrAlreadySubmitted) {
			response.Error(c, http.StatusConflict, ErrCodeAlreadySubmitted, "swap 已提交")
			return
		}
		h.handleServiceError(c, err)
		return
	}

	response.Success(c, nil)
}

// handleServiceError 处理 Service 层错误。
func (h *Handler) handleServiceError(c *gin.Context, err error) {
	if errors.Is(err, domain.ErrProviderNotFound) {
		response.Error(c, http.StatusBadRequest, ErrCodeProviderNotFound, err.Error())
		return
	}
	h.log.ErrorContext(c.Request.Context(), "供应商请求失败", "error", err)
	response.Error(c, http.StatusInternalServerError, ErrCodeProviderError, "供应商请求失败")
}

// isValidEVMAddress 校验 EVM 地址格式。
func isValidEVMAddress(addr string) bool {
	return evmAddressRegex.MatchString(addr)
}

// isValidAmount 校验 amount 格式（正整数字符串）。
func isValidAmount(amount string) bool {
	if !amountRegex.MatchString(amount) {
		return false
	}
	return amount != ""
}

// isValidTxHash 校验交易哈希格式。
func isValidTxHash(hash string) bool {
	return txHashRegex.MatchString(hash)
}
