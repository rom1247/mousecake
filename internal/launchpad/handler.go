package launchpad

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mousecake-go/mousecake-go/internal/launchpad/domain"
	"github.com/mousecake-go/mousecake-go/internal/shared/response"
)

// PrepareTxResponse Prepare 交易的 HTTP 响应结构。
type PrepareTxResponse struct {
	// ID 数据库主键。
	ID int64 `json:"id"`
	// SaleID 关联的销售 ID。
	SaleID *int64 `json:"sale_id,omitempty"`
	// PoolIndex 关联的池序号。
	PoolIndex *int64 `json:"pool_index,omitempty"`
	// OperationType 操作类型。
	OperationType string `json:"operation_type"`
	// CallerAddress 调用者地址。
	CallerAddress string `json:"caller_address"`
	// Status 交易状态。
	Status string `json:"status"`
	// Tx 交易参数（to, data, value）。
	Tx TxParams `json:"tx"`
	// ExpiresAt 过期时间。
	ExpiresAt time.Time `json:"expires_at"`
	// CreatedAt 创建时间。
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt 更新时间。
	UpdatedAt time.Time `json:"updated_at"`
}

// TxParams 交易参数结构。
type TxParams struct {
	// To 目标合约地址。
	To string `json:"to"`
	// Data 调用数据（hex 编码）。
	Data string `json:"data"`
	// Value 原生代币数量（wei）。
	Value string `json:"value"`
}

// DevExecuteResponse 开发执行接口的响应结构，包含 receipt 和事件解析数据。
type DevExecuteResponse struct {
	// ID 数据库主键。
	ID int64 `json:"id"`
	// SaleID 关联的销售 ID。
	SaleID *int64 `json:"sale_id,omitempty"`
	// PoolIndex 关联的池序号。
	PoolIndex *int64 `json:"pool_index,omitempty"`
	// OperationType 操作类型。
	OperationType string `json:"operation_type"`
	// CallerAddress 调用者地址。
	CallerAddress string `json:"caller_address"`
	// Status 交易状态。
	Status string `json:"status"`
	// Tx 交易参数。
	Tx TxParams `json:"tx"`
	// ExpiresAt 过期时间。
	ExpiresAt time.Time `json:"expires_at"`
	// CreatedAt 创建时间。
	CreatedAt time.Time `json:"created_at"`
	// UpdatedAt 更新时间。
	UpdatedAt time.Time `json:"updated_at"`
	// Receipt 交易 receipt 摘要信息。
	Receipt *ReceiptResult `json:"receipt,omitempty"`
	// Events 解析后的事件列表。
	Events []EventResponse `json:"events,omitempty"`
}

// ReceiptResult 交易 receipt 摘要信息。
type ReceiptResult struct {
	// TxHash 交易哈希。
	TxHash string `json:"tx_hash"`
	// BlockNumber 区块号。
	BlockNumber uint64 `json:"block_number"`
	// GasUsed 消耗的 gas。
	GasUsed uint64 `json:"gas_used"`
	// Status 交易状态（1=成功，0=失败）。
	Status uint64 `json:"status"`
}

// EventResponse 解析后的事件响应结构。
type EventResponse struct {
	// Name 事件名称。
	Name string `json:"name"`
	// Fields 事件字段值。
	Fields map[string]any `json:"fields"`
	// Address 合约地址。
	Address string `json:"address"`
	// LogIndex 日志索引。
	LogIndex uint `json:"log_index"`
}

// toPrepareTxResponse 将 domain.PrepareTx 转换为 HTTP 响应结构。
func toPrepareTxResponse(tx *domain.PrepareTx) PrepareTxResponse {
	return PrepareTxResponse{
		ID:            tx.ID,
		SaleID:        tx.SaleID,
		PoolIndex:     tx.PoolIndex,
		OperationType: string(tx.OperationType),
		CallerAddress: tx.CallerAddress,
		Status:        string(tx.Status),
		Tx: TxParams{
			To:    tx.TargetAddress,
			Data:  tx.Calldata,
			Value: tx.Value,
		},
		ExpiresAt: tx.ExpiresAt,
		CreatedAt: tx.CreatedAt,
		UpdatedAt: tx.UpdatedAt,
	}
}

// DevExecutor 开发环境一键执行 PrepareTx 的用例接口。
type DevExecutor interface {
	// Execute 一键执行指定 ID 的 PrepareTx，返回交易结果（含事件）。
	Execute(ctx context.Context, id int64) (*DevExecuteResult, error)
}

// Handler 注册 launchpad 模块的 HTTP 路由。
type Handler struct {
	adminSvc        *AdminService
	userSvc         *UserService
	prepareSvc      *PrepareService
	querySvc        *UserQueryService
	metaSvc         *SaleMetaService
	tokenSvc        *TokenService
	chainRefreshSvc *ChainRefreshService
	devExecuteSvc   DevExecutor
	log             *slog.Logger
}

// NewHandler 创建 Handler。
func NewHandler(
	adminSvc *AdminService,
	userSvc *UserService,
	prepareSvc *PrepareService,
	querySvc *UserQueryService,
	metaSvc *SaleMetaService,
	tokenSvc *TokenService,
	chainRefreshSvc *ChainRefreshService,
	devExecuteSvc DevExecutor,
) *Handler {
	return &Handler{
		adminSvc:        adminSvc,
		userSvc:         userSvc,
		prepareSvc:      prepareSvc,
		querySvc:        querySvc,
		metaSvc:         metaSvc,
		tokenSvc:        tokenSvc,
		chainRefreshSvc: chainRefreshSvc,
		devExecuteSvc:   devExecuteSvc,
		log:             slog.Default().With("module", "launchpad", "layer", "handler"),
	}
}

// RegisterRoutes 注册所有路由到 Gin RouterGroup。
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	admin := rg.Group("/admin")
	{
		admin.POST("/prepare/create-sale", h.adminCreateSale)
		admin.POST("/prepare/set-pool", h.adminSetPool)
		admin.POST("/prepare/set-tier-limits", h.adminSetTierLimits)
		admin.POST("/prepare/add-whitelist", h.adminAddWhitelist)
		admin.POST("/prepare/remove-whitelist", h.adminRemoveWhitelist)
		admin.POST("/prepare/set-start-end-block", h.adminSetStartEndBlock)
		admin.POST("/prepare/revoke", h.adminRevoke)
		admin.POST("/prepare/final-withdraw", h.adminFinalWithdraw)
		admin.POST("/prepare/recover-token", h.adminRecoverToken)
		admin.POST("/sale-meta", h.createSaleMeta)
		admin.PUT("/sale-meta", h.updateSaleMeta)
		admin.POST("/token", h.createToken)
		admin.PUT("/token", h.updateToken)
		admin.POST("/chain-state/refresh", h.adminChainStateRefresh)
	}

	user := rg.Group("/user")
	{
		user.POST("/prepare/deposit", h.userDeposit)
		user.POST("/prepare/harvest", h.userHarvest)
		user.POST("/prepare/release", h.userRelease)
		user.GET("/sales", h.listSales)
		user.GET("/sales/:id", h.getSaleDetail)
		user.GET("/tier", h.getUserTier)
		user.GET("/whitelist-check", h.checkWhitelist)
		user.GET("/deposits", h.getUserDeposits)
		user.GET("/harvest", h.getUserHarvest)
		user.GET("/vesting", h.getUserVesting)
		user.GET("/estimate-allocation", h.estimateAllocation)
	}

	// Prepare 通用操作
	rg.POST("/prepare/submit", h.submitPrepare)
	rg.POST("/prepare/cancel/:id", h.cancelPrepare)
	rg.GET("/prepare/:id", h.getPrepare)

	// 开发环境路由（仅非 release 模式注册）
	if gin.Mode() != gin.ReleaseMode {
		rg.POST("/dev/prepare/execute/:id", h.devExecute)
	}
}

// --- 管理员 Handler ---

// adminCreateSale 创建 IDO 销售合约。
//
//	@Summary      创建 IDO 销售
//	@Description  管理员创建新的 IDO 销售合约，生成链上交易数据
//	@Tags         Launchpad-管理员
//	@Accept       json
//	@Produce      json
//	@Param        body  body      CreateSaleInput  true  "创建销售参数"
//	@Success      200   {object}  response.Response{data=internal_launchpad_domain.PrepareTx}
//	@Failure      400   {object}  response.Response
//	@Failure      500   {object}  response.Response
//	@Router       /launchpad/admin/prepare/create-sale [post]
func (h *Handler) adminCreateSale(c *gin.Context) {
	var input CreateSaleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	tx, err := h.adminSvc.CreateSale(c.Request.Context(), input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	prepareTxResponse := toPrepareTxResponse(tx)
	response.Success(c, prepareTxResponse)
}

// adminSetPool 设置池参数。
//
//	@Summary      设置池参数
//	@Description  管理员设置 IDO 销售池的参数（发售量、募资量、用户限额、vesting 配置等）
//	@Tags         Launchpad-管理员
//	@Accept       json
//	@Produce      json
//	@Param        body  body      SetPoolInput  true  "池参数"
//	@Success      200   {object}  response.Response{data=internal_launchpad_domain.PrepareTx}
//	@Failure      400   {object}  response.Response
//	@Failure      500   {object}  response.Response
//	@Router       /launchpad/admin/prepare/set-pool [post]
func (h *Handler) adminSetPool(c *gin.Context) {
	var input SetPoolInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	tx, err := h.adminSvc.SetPool(c.Request.Context(), input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Success(c, toPrepareTxResponse(tx))
}

// adminSetTierLimits 设置 Tier 额度。
//
//	@Summary      设置 Tier 额度
//	@Description  管理员设置指定 Tier 的认购额度限制
//	@Tags         Launchpad-管理员
//	@Accept       json
//	@Produce      json
//	@Param        body  body      SetTierLimitsInput  true  "Tier 额度参数"
//	@Success      200   {object}  response.Response{data=internal_launchpad_domain.PrepareTx}
//	@Failure      400   {object}  response.Response
//	@Failure      500   {object}  response.Response
//	@Router       /launchpad/admin/prepare/set-tier-limits [post]
func (h *Handler) adminSetTierLimits(c *gin.Context) {
	var input SetTierLimitsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	tx, err := h.adminSvc.SetTierLimits(c.Request.Context(), input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Success(c, toPrepareTxResponse(tx))
}

// adminAddWhitelist 添加白名单。
//
//	@Summary      添加白名单
//	@Description  管理员将用户地址添加到销售白名单
//	@Tags         Launchpad-管理员
//	@Accept       json
//	@Produce      json
//	@Param        body  body      WhitelistInput  true  "白名单参数"
//	@Success      200   {object}  response.Response{data=internal_launchpad_domain.PrepareTx}
//	@Failure      400   {object}  response.Response
//	@Failure      500   {object}  response.Response
//	@Router       /launchpad/admin/prepare/add-whitelist [post]
func (h *Handler) adminAddWhitelist(c *gin.Context) {
	var input WhitelistInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	tx, err := h.adminSvc.AddWhitelist(c.Request.Context(), input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Success(c, toPrepareTxResponse(tx))
}

// adminRemoveWhitelist 移除白名单。
//
//	@Summary      移除白名单
//	@Description  管理员将用户地址从销售白名单中移除
//	@Tags         Launchpad-管理员
//	@Accept       json
//	@Produce      json
//	@Param        body  body      WhitelistInput  true  "白名单参数"
//	@Success      200   {object}  response.Response{data=internal_launchpad_domain.PrepareTx}
//	@Failure      400   {object}  response.Response
//	@Failure      500   {object}  response.Response
//	@Router       /launchpad/admin/prepare/remove-whitelist [post]
func (h *Handler) adminRemoveWhitelist(c *gin.Context) {
	var input WhitelistInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	tx, err := h.adminSvc.RemoveWhitelist(c.Request.Context(), input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Success(c, tx)
}

// adminSetStartEndBlock 设置销售时间窗。
//
//	@Summary      设置销售时间窗
//	@Description  管理员设置销售的开始和结束区块号
//	@Tags         Launchpad-管理员
//	@Accept       json
//	@Produce      json
//	@Param        body  body      SetStartEndBlockInput  true  "时间窗参数"
//	@Success      200   {object}  response.Response{data=internal_launchpad_domain.PrepareTx}
//	@Failure      400   {object}  response.Response
//	@Failure      500   {object}  response.Response
//	@Router       /launchpad/admin/prepare/set-start-end-block [post]
func (h *Handler) adminSetStartEndBlock(c *gin.Context) {
	var input SetStartEndBlockInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	tx, err := h.adminSvc.SetStartEndBlock(c.Request.Context(), input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Success(c, toPrepareTxResponse(tx))
}

// adminRevoke 撤销 vesting。
//
//	@Summary      撤销 vesting
//	@Description  管理员撤销指定池的 vesting 计划
//	@Tags         Launchpad-管理员
//	@Accept       json
//	@Produce      json
//	@Param        body  body      RevokeInput  true  "撤销参数"
//	@Success      200   {object}  response.Response{data=internal_launchpad_domain.PrepareTx}
//	@Failure      400   {object}  response.Response
//	@Failure      500   {object}  response.Response
//	@Router       /launchpad/admin/prepare/revoke [post]
func (h *Handler) adminRevoke(c *gin.Context) {
	var input RevokeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	tx, err := h.adminSvc.Revoke(c.Request.Context(), input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Success(c, toPrepareTxResponse(tx))
}

// adminFinalWithdraw 提取资金。
//
//	@Summary      提取资金
//	@Description  管理员提取销售结束后的募资和剩余发售代币
//	@Tags         Launchpad-管理员
//	@Accept       json
//	@Produce      json
//	@Param        body  body      FinalWithdrawInput  true  "提取参数"
//	@Success      200   {object}  response.Response{data=internal_launchpad_domain.PrepareTx}
//	@Failure      400   {object}  response.Response
//	@Failure      500   {object}  response.Response
//	@Router       /launchpad/admin/prepare/final-withdraw [post]
func (h *Handler) adminFinalWithdraw(c *gin.Context) {
	var input FinalWithdrawInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	tx, err := h.adminSvc.FinalWithdraw(c.Request.Context(), input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Success(c, toPrepareTxResponse(tx))
}

// adminRecoverToken 救援误转代币。
//
//	@Summary      救援误转代币
//	@Description  管理员将误转入合约的代币提取到指定地址
//	@Tags         Launchpad-管理员
//	@Accept       json
//	@Produce      json
//	@Param        body  body      RecoverTokenInput  true  "救援代币参数"
//	@Success      200   {object}  response.Response{data=internal_launchpad_domain.PrepareTx}
//	@Failure      400   {object}  response.Response
//	@Failure      500   {object}  response.Response
//	@Router       /launchpad/admin/prepare/recover-token [post]
func (h *Handler) adminRecoverToken(c *gin.Context) {
	var input RecoverTokenInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	tx, err := h.adminSvc.RecoverToken(c.Request.Context(), input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Success(c, toPrepareTxResponse(tx))
}

// adminChainStateRefresh 刷新链上状态到数据库。
//
//	@Summary      刷新链上状态
//	@Description  管理员按 scope 刷新链上合约状态（Sale/Pool/Tier/UserPool/Vesting）并更新数据库
//	@Tags         Launchpad-管理员
//	@Accept       json
//	@Produce      json
//	@Param        body  body  ChainStateRefreshRequest  true  "刷新参数"
//	@Success      200   {object}  response.Response{data=ChainStateRefreshResponse}
//	@Failure      400   {object}  response.Response
//	@Failure      500   {object}  response.Response
//	@Router       /launchpad/admin/chain-state/refresh [post]
func (h *Handler) adminChainStateRefresh(c *gin.Context) {
	var req ChainStateRefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	result, err := h.chainRefreshSvc.ChainStateRefresh(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "链上状态刷新失败")
		return
	}
	response.Success(c, result)
}

// --- 销售元信息和代币信息管理 Handler ---

// createSaleMeta 创建销售元信息。
//
//	@Summary      创建销售元信息
//	@Description  管理员为销售创建展示元信息（标题、描述、Banner 等）
//	@Tags         Launchpad-管理员
//	@Accept       json
//	@Produce      json
//	@Param        body  body      CreateSaleMetaInput  true  "销售元信息"
//	@Success      200   {object}  response.Response{data=internal_launchpad_domain.SaleMeta}
//	@Failure      400   {object}  response.Response
//	@Failure      500   {object}  response.Response
//	@Router       /launchpad/admin/sale-meta [post]
func (h *Handler) createSaleMeta(c *gin.Context) {
	var input CreateSaleMetaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	meta, err := h.metaSvc.CreateSaleMeta(c.Request.Context(), input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Success(c, meta)
}

// updateSaleMeta 更新销售元信息。
//
//	@Summary      更新销售元信息
//	@Description  管理员更新销售的展示元信息
//	@Tags         Launchpad-管理员
//	@Accept       json
//	@Produce      json
//	@Param        body  body      UpdateSaleMetaInput  true  "更新元信息"
//	@Success      200   {object}  response.Response{data=internal_launchpad_domain.SaleMeta}
//	@Failure      400   {object}  response.Response
//	@Failure      404   {object}  response.Response
//	@Failure      500   {object}  response.Response
//	@Router       /launchpad/admin/sale-meta [put]
func (h *Handler) updateSaleMeta(c *gin.Context) {
	var input UpdateSaleMetaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	meta, err := h.metaSvc.UpdateSaleMeta(c.Request.Context(), input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Success(c, meta)
}

// createToken 创建代币元信息。
//
//	@Summary      创建代币信息
//	@Description  管理员创建代币的元信息（名称、符号、精度等）
//	@Tags         Launchpad-管理员
//	@Accept       json
//	@Produce      json
//	@Param        body  body      CreateTokenInput  true  "代币信息"
//	@Success      200   {object}  response.Response{data=internal_launchpad_domain.Token}
//	@Failure      400   {object}  response.Response
//	@Failure      500   {object}  response.Response
//	@Router       /launchpad/admin/token [post]
func (h *Handler) createToken(c *gin.Context) {
	var input CreateTokenInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	token, err := h.tokenSvc.CreateToken(c.Request.Context(), input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Success(c, token)
}

// updateToken 更新代币元信息。
//
//	@Summary      更新代币信息
//	@Description  管理员更新代币的元信息
//	@Tags         Launchpad-管理员
//	@Accept       json
//	@Produce      json
//	@Param        body  body      UpdateTokenInput  true  "更新代币信息"
//	@Success      200   {object}  response.Response{data=internal_launchpad_domain.Token}
//	@Failure      400   {object}  response.Response
//	@Failure      404   {object}  response.Response
//	@Failure      500   {object}  response.Response
//	@Router       /launchpad/admin/token [put]
func (h *Handler) updateToken(c *gin.Context) {
	var input UpdateTokenInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	token, err := h.tokenSvc.UpdateToken(c.Request.Context(), input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Success(c, token)
}

// --- 用户操作 Handler ---

// userDeposit 用户申购。
//
//	@Summary      用户申购
//	@Description  用户参与 IDO 申购，校验资格后生成链上交易数据
//	@Tags         Launchpad-用户
//	@Accept       json
//	@Produce      json
//	@Param        body  body      DepositInput  true  "申购参数"
//	@Success      200   {object}  response.Response{data=internal_launchpad_domain.PrepareTx}
//	@Failure      400   {object}  response.Response
//	@Failure      500   {object}  response.Response
//	@Router       /launchpad/user/prepare/deposit [post]
func (h *Handler) userDeposit(c *gin.Context) {
	var input DepositInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	tx, err := h.userSvc.Deposit(c.Request.Context(), input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Success(c, toPrepareTxResponse(tx))
}

// userHarvest 用户结算。
//
//	@Summary      用户结算
//	@Description  用户结算 IDO 认购结果，生成链上交易数据
//	@Tags         Launchpad-用户
//	@Accept       json
//	@Produce      json
//	@Param        body  body      HarvestInput  true  "结算参数"
//	@Success      200   {object}  response.Response{data=internal_launchpad_domain.PrepareTx}
//	@Failure      400   {object}  response.Response
//	@Failure      500   {object}  response.Response
//	@Router       /launchpad/user/prepare/harvest [post]
func (h *Handler) userHarvest(c *gin.Context) {
	var input HarvestInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	tx, err := h.userSvc.Harvest(c.Request.Context(), input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Success(c, toPrepareTxResponse(tx))
}

// userRelease 用户释放 vesting。
//
//	@Summary      用户释放 vesting
//	@Description  用户释放 vesting 锁仓代币，生成链上交易数据
//	@Tags         Launchpad-用户
//	@Accept       json
//	@Produce      json
//	@Param        body  body      ReleaseInput  true  "释放参数"
//	@Success      200   {object}  response.Response{data=internal_launchpad_domain.PrepareTx}
//	@Failure      400   {object}  response.Response
//	@Failure      500   {object}  response.Response
//	@Router       /launchpad/user/prepare/release [post]
func (h *Handler) userRelease(c *gin.Context) {
	var input ReleaseInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	tx, err := h.userSvc.Release(c.Request.Context(), input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Success(c, toPrepareTxResponse(tx))
}

// --- 用户查询 Handler ---

// listSales 获取销售列表。
//
//	@Summary      获取销售列表
//	@Description  分页获取公开的 IDO 销售列表
//	@Tags         Launchpad-用户
//	@Produce      json
//	@Param        page       query  int  false  "页码"       default(1)
//	@Param        page_size  query  int  false  "每页数量"   default(10)  maximum(100)
//	@Success      200  {object}  response.Response
//	@Failure      500  {object}  response.Response
//	@Router       /launchpad/user/sales [get]
func (h *Handler) listSales(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	result, err := h.querySvc.ListPublicSales(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Success(c, result)
}

// getSaleDetail 获取销售详情。
//
//	@Summary      获取销售详情
//	@Description  根据 ID 获取 IDO 销售的详细信息
//	@Tags         Launchpad-用户
//	@Produce      json
//	@Param        id  path  int  true  "销售 ID"  example(1)
//	@Success      200  {object}  response.Response
//	@Failure      400  {object}  response.Response
//	@Failure      500  {object}  response.Response
//	@Router       /launchpad/user/sales/{id} [get]
func (h *Handler) getSaleDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "无效 ID")
		return
	}
	detail, err := h.querySvc.GetSaleDetail(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Success(c, detail)
}

// getUserTier 获取用户 Tier 等级。
//
//	@Summary      获取用户 Tier 等级
//	@Description  根据链上 MouseTier 合约查询用户的 Tier 等级
//	@Tags         Launchpad-用户
//	@Produce      json
//	@Param        mouse_tier_address  query  string  true  "MouseTier 合约地址"
//	@Param        user_address        query  string  true  "用户钱包地址"
//	@Success      200  {object}  response.Response
//	@Failure      400  {object}  response.Response
//	@Failure      500  {object}  response.Response
//	@Router       /launchpad/user/tier [get]
func (h *Handler) getUserTier(c *gin.Context) {
	mouseTierAddr := c.Query("mouse_tier_address")
	userAddr := c.Query("user_address")
	if mouseTierAddr == "" || userAddr == "" {
		response.Error(c, http.StatusBadRequest, 400, "缺少 mouse_tier_address 或 user_address")
		return
	}
	tier, err := h.querySvc.GetUserTier(c.Request.Context(), mouseTierAddr, userAddr)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Success(c, gin.H{"tier": tier})
}

// checkWhitelist 检查白名单。
//
//	@Summary      检查白名单
//	@Description  检查用户是否在指定销售的白名单中
//	@Tags         Launchpad-用户
//	@Produce      json
//	@Param        sale_id        query  int     true  "销售 ID"
//	@Param        user_address   query  string  true  "用户钱包地址"
//	@Success      200  {object}  response.Response
//	@Failure      400  {object}  response.Response
//	@Failure      500  {object}  response.Response
//	@Router       /launchpad/user/whitelist-check [get]
func (h *Handler) checkWhitelist(c *gin.Context) {
	saleID, _ := strconv.ParseInt(c.Query("sale_id"), 10, 64)
	userAddr := c.Query("user_address")
	if saleID == 0 || userAddr == "" {
		response.Error(c, http.StatusBadRequest, 400, "缺少 sale_id 或 user_address")
		return
	}
	whitelisted, err := h.querySvc.CheckWhitelist(c.Request.Context(), saleID, userAddr)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Success(c, gin.H{"whitelisted": whitelisted})
}

// getUserDeposits 获取用户申购记录。
//
//	@Summary      获取用户申购记录
//	@Description  查询用户在指定销售中的申购记录
//	@Tags         Launchpad-用户
//	@Produce      json
//	@Param        sale_id        query  int     true  "销售 ID"
//	@Param        user_address   query  string  true  "用户钱包地址"
//	@Success      200  {object}  response.Response
//	@Failure      400  {object}  response.Response
//	@Failure      500  {object}  response.Response
//	@Router       /launchpad/user/deposits [get]
func (h *Handler) getUserDeposits(c *gin.Context) {
	saleID, _ := strconv.ParseInt(c.Query("sale_id"), 10, 64)
	userAddr := c.Query("user_address")
	if saleID == 0 || userAddr == "" {
		response.Error(c, http.StatusBadRequest, 400, "缺少 sale_id 或 user_address")
		return
	}
	deposits, err := h.querySvc.GetUserDeposits(c.Request.Context(), userAddr, saleID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Success(c, deposits)
}

// getUserHarvest 获取用户结算记录。
//
//	@Summary      获取用户结算记录
//	@Description  查询用户在指定销售中的结算记录
//	@Tags         Launchpad-用户
//	@Produce      json
//	@Param        sale_id        query  int     true  "销售 ID"
//	@Param        user_address   query  string  true  "用户钱包地址"
//	@Success      200  {object}  response.Response
//	@Failure      400  {object}  response.Response
//	@Failure      500  {object}  response.Response
//	@Router       /launchpad/user/harvest [get]
func (h *Handler) getUserHarvest(c *gin.Context) {
	saleID, _ := strconv.ParseInt(c.Query("sale_id"), 10, 64)
	userAddr := c.Query("user_address")
	if saleID == 0 || userAddr == "" {
		response.Error(c, http.StatusBadRequest, 400, "缺少 sale_id 或 user_address")
		return
	}
	harvests, err := h.querySvc.GetUserHarvest(c.Request.Context(), userAddr, saleID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Success(c, harvests)
}

// getUserVesting 获取用户 vesting 计划。
//
//	@Summary      获取用户 vesting 计划
//	@Description  查询用户的 vesting 锁仓释放计划
//	@Tags         Launchpad-用户
//	@Produce      json
//	@Param        beneficiary  query  string  true  "受益人钱包地址"
//	@Success      200  {object}  response.Response
//	@Failure      400  {object}  response.Response
//	@Failure      500  {object}  response.Response
//	@Router       /launchpad/user/vesting [get]
func (h *Handler) getUserVesting(c *gin.Context) {
	beneficiary := c.Query("beneficiary")
	if beneficiary == "" {
		response.Error(c, http.StatusBadRequest, 400, "缺少 beneficiary")
		return
	}
	schedules, err := h.querySvc.GetUserVesting(c.Request.Context(), beneficiary)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Success(c, schedules)
}

// estimateAllocation 预估配售。
//
//	@Summary      预估配售
//	@Description  根据用户投入金额预估可获得的发售代币数量
//	@Tags         Launchpad-用户
//	@Accept       json
//	@Produce      json
//	@Param        sale_id      query  int     true  "销售 ID"
//	@Param        pool_index   query  int     true  "池索引"
//	@Param        user_total   query  string  true  "用户投入金额"
//	@Success      200  {object}  response.Response
//	@Failure      400  {object}  response.Response
//	@Failure      500  {object}  response.Response
//	@Router       /launchpad/user/estimate-allocation [get]
func (h *Handler) estimateAllocation(c *gin.Context) {
	var input EstimateAllocationInput
	if err := c.ShouldBindQuery(&input); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	result, err := h.querySvc.EstimateAllocation(c.Request.Context(), input)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Success(c, result)
}

// --- Prepare 通用 Handler ---

// SubmitPrepareRequest 提交 Prepare txHash 的请求。
type SubmitPrepareRequest struct {
	ID     int64  `json:"id" binding:"required"`
	TxHash string `json:"tx_hash" binding:"required"`
}

// submitPrepare 提交 Prepare 交易哈希。
//
//	@Summary      提交 Prepare 交易
//	@Description  提交已广播的链上交易哈希，完成 Prepare 记录
//	@Tags         Launchpad-通用
//	@Accept       json
//	@Produce      json
//	@Param        body  body      SubmitPrepareRequest  true  "提交参数"
//	@Success      200   {object}  response.Response{data=internal_launchpad_domain.PrepareTx}
//	@Failure      400   {object}  response.Response
//	@Failure      500   {object}  response.Response
//	@Router       /launchpad/prepare/submit [post]
func (h *Handler) submitPrepare(c *gin.Context) {
	var req SubmitPrepareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, err.Error())
		return
	}
	tx, err := h.prepareSvc.Submit(c.Request.Context(), req.ID, req.TxHash)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Success(c, toPrepareTxResponse(tx))
}

// cancelPrepare 取消 Prepare 交易。
//
//	@Summary      取消 Prepare 交易
//	@Description  取消指定 ID 的 Prepare 交易
//	@Tags         Launchpad-通用
//	@Produce      json
//	@Param        id  path  int  true  "Prepare ID"  example(1)
//	@Success      200  {object}  response.Response
//	@Failure      400  {object}  response.Response
//	@Failure      500  {object}  response.Response
//	@Router       /launchpad/prepare/cancel/{id} [post]
func (h *Handler) cancelPrepare(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "无效 ID")
		return
	}
	if err := h.prepareSvc.Cancel(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Success(c, nil)
}

// getPrepare 获取 Prepare 交易详情。
//
//	@Summary      获取 Prepare 详情
//	@Description  根据 ID 获取 Prepare 交易的详细信息
//	@Tags         Launchpad-通用
//	@Produce      json
//	@Param        id  path  int  true  "Prepare ID"  example(1)
//	@Success      200  {object}  response.Response{data=internal_launchpad_domain.PrepareTx}
//	@Failure      400  {object}  response.Response
//	@Failure      500  {object}  response.Response
//	@Router       /launchpad/prepare/{id} [get]
func (h *Handler) getPrepare(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "无效 ID")
		return
	}
	// 通过 Service 层查询
	tx, err := h.prepareSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Success(c, toPrepareTxResponse(tx))
}

// devExecute 开发环境一键执行 PrepareTx。
//
//	@Summary      开发执行 PrepareTx
//	@Description  开发环境一键签名广播 PrepareTx（仅 debug/test 模式）
//	@Tags         Launchpad-开发
//	@Produce      json
//	@Param        id  path  int  true  "Prepare ID"  example(1)
//	@Success      200  {object}  response.Response
//	@Failure      400  {object}  response.Response
//	@Failure      500  {object}  response.Response
//	@Router       /launchpad/dev/prepare/execute/{id} [post]
func (h *Handler) devExecute(c *gin.Context) {
	if h.devExecuteSvc == nil {
		response.Error(c, http.StatusServiceUnavailable, 503, "DevExecute 服务未配置，请检查 AdminPrivateKey")
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "无效 ID")
		return
	}
	result, err := h.devExecuteSvc.Execute(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	// 查询最新状态
	tx, err := h.prepareSvc.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}
	response.Success(c, toDevExecuteResponse(tx, result))
}

// toDevExecuteResponse 将 PrepareTx 和 DevExecuteResult 合并为开发执行响应。
func toDevExecuteResponse(tx *domain.PrepareTx, result *DevExecuteResult) DevExecuteResponse {
	resp := DevExecuteResponse{
		ID:            tx.ID,
		SaleID:        tx.SaleID,
		PoolIndex:     tx.PoolIndex,
		OperationType: string(tx.OperationType),
		CallerAddress: tx.CallerAddress,
		Status:        string(tx.Status),
		Tx: TxParams{
			To:    tx.TargetAddress,
			Data:  tx.Calldata,
			Value: tx.Value,
		},
		ExpiresAt: tx.ExpiresAt,
		CreatedAt: tx.CreatedAt,
		UpdatedAt: tx.UpdatedAt,
	}
	if result != nil {
		resp.Receipt = &ReceiptResult{
			TxHash:      result.TxHash,
			BlockNumber: result.BlockNumber,
			GasUsed:     result.GasUsed,
			Status:      result.Status,
		}
		for _, e := range result.Events {
			resp.Events = append(resp.Events, EventResponse{
				Name:     e.Name,
				Fields:   e.Fields,
				Address:  e.Address.Hex(),
				LogIndex: e.LogIndex,
			})
		}
	}
	return resp
}
