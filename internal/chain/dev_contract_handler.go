package chain

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mousecake-go/mousecake-go/internal/shared/response"
)

// DevContractHandler 开发环境通用合约查询 HTTP handler。
// 仅在非 release 模式下注册路由。
type DevContractHandler struct {
	svc *DevContractService
	log *slog.Logger
}

// NewDevContractHandler 创建通用合约查询 handler。
func NewDevContractHandler(svc *DevContractService) *DevContractHandler {
	return &DevContractHandler{
		svc: svc,
		log: slog.Default().With("module", "chain", "layer", "handler"),
	}
}

// RegisterRoutes 注册开发环境路由，仅在非 release 模式下生效。
func (h *DevContractHandler) RegisterRoutes(rg *gin.RouterGroup) {
	if gin.Mode() == gin.ReleaseMode {
		return
	}

	dev := rg.Group("/dev/contract")
	{
		dev.POST("/query", h.query)
		dev.GET("/contracts", h.listContracts)
		dev.GET("/contracts/:name/methods", h.listMethods)
	}
}

// query 处理通用合约查询/执行请求。
// 根据方法的 StateMutability 自动选择只读查询（eth_call）或写入交易（sendTransaction）。
//
//	@Summary      通用合约调用
//	@Description  根据合约名和方法名动态调用合约，view/pure 方法执行只读查询，其他方法签名广播交易
//	@Tags         开发-合约查询
//	@Accept       json
//	@Produce      json
//	@Param        body  body      ContractQueryRequest  true  "合约调用请求"
//	@Success      200   {object}  response.Response{data=ContractQueryResult}
//	@Failure      400   {object}  response.Response
//	@Failure      500   {object}  response.Response
//	@Router       /dev/contract/query [post]
func (h *DevContractHandler) query(c *gin.Context) {
	var req ContractQueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 400, "请求参数错误: "+err.Error())
		return
	}

	result, err := h.svc.Query(c.Request.Context(), req)
	if err != nil {
		h.log.ErrorContext(c.Request.Context(), "合约查询失败",
			"contract", req.Contract,
			"method", req.Method,
			"error", err,
		)
		response.Error(c, http.StatusInternalServerError, 500, err.Error())
		return
	}

	response.Success(c, result)
}

// listContracts 返回所有已注册合约名称。
//
//	@Summary      列出已注册合约
//	@Description  返回所有已注册到 ABI 注册中心的合约名称列表
//	@Tags         开发-合约查询
//	@Produce      json
//	@Success      200  {object}  response.Response{data=[]string}
//	@Router       /dev/contract/contracts [get]
func (h *DevContractHandler) listContracts(c *gin.Context) {
	names := h.svc.registry.ListContracts()
	response.Success(c, names)
}

// listMethods 返回指定合约的所有方法签名。
//
//	@Summary      列出合约方法
//	@Description  根据合约名返回所有方法的签名、状态可变性和参数类型信息
//	@Tags         开发-合约查询
//	@Produce      json
//	@Param        name  path      string  true  "合约名称"
//	@Success      200   {object}  response.Response{data=[]MethodInfo}
//	@Failure      404   {object}  response.Response
//	@Router       /dev/contract/contracts/{name}/methods [get]
func (h *DevContractHandler) listMethods(c *gin.Context) {
	name := c.Param("name")
	methods, err := h.svc.registry.ListMethods(name)
	if err != nil {
		response.Error(c, http.StatusNotFound, 404, err.Error())
		return
	}
	response.Success(c, methods)
}
