// Package user 实现用户认证模块的业务逻辑和数据访问。
package user

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/mousecake-go/mousecake-go/config"
	"github.com/mousecake-go/mousecake-go/internal/shared/middleware"

	"github.com/mousecake-go/mousecake-go/internal/shared/auth"
	"github.com/mousecake-go/mousecake-go/internal/shared/errs"
	"github.com/mousecake-go/mousecake-go/internal/shared/response"
)

// walletNonceRequest 是 wallet nonce 请求体。
type walletNonceRequest struct {
	Address string `json:"address" binding:"required" example:"0x1234567890abcdef1234567890abcdef12345678"`
}

// walletVerifyRequest 是 wallet verify 请求体。
type walletVerifyRequest struct {
	Message   string `json:"message" binding:"required" example:"mousecake-go wants you to sign in with your Ethereum account..."`
	Signature string `json:"signature" binding:"required" example:"0xabcdef..."`
}

// adminLoginRequest 是管理员登录请求体。
type adminLoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"********"`
}

// tokenResponse 是 token 响应数据。
type tokenResponse struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIs..."`
	TokenType   string `json:"token_type" example:"Bearer"`
	ExpiresIn   int    `json:"expires_in" example:"900"`
}

// nonceResponse 是 nonce 响应数据。
type nonceResponse struct {
	Message string `json:"message" example:"mousecake-go wants you to sign in..."`
	Nonce   string `json:"nonce" example:"abc123def456"`
}

// userInfoResponse 是用户信息响应数据。
type userInfoResponse struct {
	ID       int64  `json:"id" example:"1"`
	Type     string `json:"type" example:"wallet"`
	Address  string `json:"address,omitempty" example:"0x1234567890abcdef1234567890abcdef12345678"`
	Username string `json:"username,omitempty" example:"admin"`
	Name     string `json:"name,omitempty" example:"张三"`
	Nickname string `json:"nickname,omitempty" example:"mouse"`
	IsAdmin  bool   `json:"is_admin" example:"false"`
}

// Handler 处理认证相关的 HTTP 请求。
type Handler struct {
	svc    *Service
	jwtSvc *auth.JWTService
	log    *slog.Logger
}

// NewHandler 创建 Handler 实例。
func NewHandler(svc *Service, jwtSvc *auth.JWTService) *Handler {
	return &Handler{
		svc:    svc,
		jwtSvc: jwtSvc,
		log:    slog.Default().With("module", "user", "layer", "handler"),
	}
}

// RegisterRoutes 将 handler 绑定到 Gin RouterGroup，应用限流和认证中间件。
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup, rlCfg config.RateLimitConfig) {
	addressRL := middleware.NewAddressRateLimit(rlCfg.Address.Rate, rlCfg.Address.Burst)

	auth := rg.Group("")
	{
		auth.POST("/wallet/nonce", addressRL, h.walletNonce)
		auth.POST("/wallet/verify", addressRL, h.walletVerify)
		auth.POST("/admin/login", h.adminLogin)
	}

	accountRL := middleware.NewAccountRateLimit(rlCfg.Account.Rate, rlCfg.Account.Burst)
	me := rg.Group("")
	me.Use(authMiddleware(h.jwtSvc), accountRL)
	{
		me.GET("/me", h.getMe)
	}
}

// walletNonce 获取钱包登录 nonce。
//
//	@Summary      获取钱包登录 nonce
//	@Description  根据钱包地址生成 SIWE 登录消息和 nonce
//	@Tags         认证
//	@Accept       json
//	@Produce      json
//	@Param        body  body      walletNonceRequest  true  "钱包地址"
//	@Success      200   {object}  response.Response{data=nonceResponse}
//	@Failure      400   {object}  response.Response
//	@Router       /auth/wallet/nonce [post]
func (h *Handler) walletNonce(c *gin.Context) {
	var req walletNonceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "缺少必要参数")
		return
	}

	if !ethAddressRegex.MatchString(req.Address) {
		response.Error(c, http.StatusBadRequest, 40001, "地址格式无效")
		return
	}

	msg, nonce, err := h.svc.GenerateSIWENonce(c.Request.Context(), req.Address)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 40001, err.Error())
		return
	}

	response.Success(c, nonceResponse{Message: msg, Nonce: nonce})
}

// walletVerify 验证钱包签名。
//
//	@Summary      验证钱包签名
//	@Description  验证 SIWE 消息签名，签发 JWT token
//	@Tags         认证
//	@Accept       json
//	@Produce      json
//	@Param        body  body      walletVerifyRequest  true  "签名验证请求"
//	@Success      200   {object}  response.Response{data=tokenResponse}
//	@Failure      400   {object}  response.Response
//	@Failure      401   {object}  response.Response
//	@Router       /auth/wallet/verify [post]
func (h *Handler) walletVerify(c *gin.Context) {
	var req walletVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "缺少必要字段")
		return
	}

	token, err := h.svc.VerifySIWESignature(c.Request.Context(), req.Message, req.Signature)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, tokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int(h.jwtSvc.ExpireDuration(auth.TokenTypeWallet).Seconds()),
	})
}

// adminLogin 管理员登录。
//
//	@Summary      管理员登录
//	@Description  使用用户名和密码登录管理员账号，签发 JWT token
//	@Tags         认证
//	@Accept       json
//	@Produce      json
//	@Param        body  body      adminLoginRequest  true  "管理员登录请求"
//	@Success      200   {object}  response.Response{data=tokenResponse}
//	@Failure      400   {object}  response.Response
//	@Failure      401   {object}  response.Response
//	@Router       /auth/admin/login [post]
func (h *Handler) adminLogin(c *gin.Context) {
	var req adminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, 40001, "缺少必要参数")
		return
	}

	if len(req.Password) < 8 {
		response.Error(c, http.StatusBadRequest, 40001, "密码长度不足")
		return
	}

	token, err := h.svc.AdminLogin(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.Success(c, tokenResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int(h.jwtSvc.ExpireDuration(auth.TokenTypeAdmin).Seconds()),
	})
}

// getMe 获取当前用户信息。
//
//	@Summary      获取当前用户信息
//	@Description  根据 Bearer Token 获取当前登录用户信息
//	@Tags         认证
//	@Accept       json
//	@Produce      json
//	@Param        Authorization  header    string  true  "Bearer Token"  default(Bearer )
//	@Success      200  {object}  response.Response{data=userInfoResponse}
//	@Failure      401  {object}  response.Response
//	@Router       /auth/me [get]
//	@Security     BearerAuth
func (h *Handler) getMe(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, 40105, "Token 无效或已过期")
		return
	}
	userID, ok := userIDVal.(int64)
	if !ok {
		response.Error(c, http.StatusUnauthorized, 40105, "Token 无效或已过期")
		return
	}

	user, err := h.svc.GetCurrentUser(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, 40105, "用户不存在")
		return
	}

	tokenTypeVal, _ := c.Get("type")
	tokenType, ok := tokenTypeVal.(string)
	if !ok {
		response.Error(c, http.StatusUnauthorized, 40105, "Token 无效或已过期")
		return
	}

	resp := userInfoResponse{
		ID:       user.ID,
		Type:     tokenType,
		Address:  user.Address,
		Username: user.Username,
		Name:     user.Name,
		Nickname: user.Nickname,
		IsAdmin:  user.IsAdmin,
	}

	response.Success(c, resp)
}

func handleServiceError(c *gin.Context, err error) {
	var svcErr *serviceError
	if errors.As(err, &svcErr) {
		response.Error(c, http.StatusUnauthorized, svcErr.Code(), svcErr.Error())
		return
	}

	response.Error(c, http.StatusInternalServerError, errs.CodeInternal, errs.GetErrorMessage(errs.CodeInternal))
}

func authMiddleware(jwtSvc *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, 40105, "Token 无效或已过期")
			c.Abort()
			return
		}

		token, ok := strings.CutPrefix(authHeader, "Bearer ")
		if !ok {
			response.Error(c, http.StatusUnauthorized, 40105, "Token 无效或已过期")
			c.Abort()
			return
		}

		claims, err := jwtSvc.ValidateToken(token)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, 40105, "Token 无效或已过期")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("type", string(claims.Type))
		c.Set("is_admin", claims.IsAdmin)
		c.Set("account_key", "account:"+strconv.FormatInt(claims.UserID, 10))
		c.Next()
	}
}
