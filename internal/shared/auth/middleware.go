package auth

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/mousecake-go/mousecake-go/internal/shared/errs"
	"github.com/mousecake-go/mousecake-go/internal/shared/response"
)

// NewAuthMiddleware 创建 JWT 认证中间件，从 Authorization 头提取并验证 Bearer Token。
func NewAuthMiddleware(svc *JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, errs.CodeTokenInvalid, errs.GetErrorMessage(errs.CodeTokenInvalid))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Error(c, http.StatusUnauthorized, errs.CodeTokenInvalid, errs.GetErrorMessage(errs.CodeTokenInvalid))
			c.Abort()
			return
		}

		claims, err := svc.ValidateToken(parts[1])
		if err != nil {
			response.Error(c, http.StatusUnauthorized, errs.CodeTokenInvalid, errs.GetErrorMessage(errs.CodeTokenInvalid))
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

// NewAdminMiddleware 创建管理员权限检查中间件。
func NewAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("is_admin")
		if !exists || !isAdmin.(bool) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Next()
	}
}
