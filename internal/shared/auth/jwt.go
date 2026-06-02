// Package auth 提供 JWT 签发、验证和 Gin 认证中间件功能。
package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenType 表示 JWT token 的身份类型。
type TokenType string

const (
	TokenTypeWallet TokenType = "wallet"
	TokenTypeAdmin  TokenType = "admin"
)

// CustomClaims 是自定义的 JWT claims 结构体。
type CustomClaims struct {
	// Token 身份类型（wallet 或 admin）。
	Type TokenType `json:"type"`
	// UserID 用户 ID。
	UserID int64 `json:"user_id"`
	// IsAdmin 是否为管理员。
	IsAdmin bool `json:"is_admin"`
	jwt.RegisteredClaims
}

// JWTService 提供 JWT 签发和验证功能。
type JWTService struct {
	secret       []byte
	walletExpire time.Duration
	adminExpire  time.Duration
}

// NewJWTService 创建 JWT 服务实例。
func NewJWTService(secret string, walletExpire, adminExpire time.Duration) *JWTService {
	return &JWTService{
		secret:       []byte(secret),
		walletExpire: walletExpire,
		adminExpire:  adminExpire,
	}
}

// IssueToken 签发 JWT access token。
func (s *JWTService) IssueToken(tokenType TokenType, userID int64, isAdmin bool) (string, error) {
	expire := s.walletExpire
	if tokenType == TokenTypeAdmin {
		expire = s.adminExpire
	}

	now := time.Now()
	claims := CustomClaims{
		Type:    tokenType,
		UserID:  userID,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expire)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    "mousecake",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(s.secret)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signed, nil
}

// ValidateToken 验证 JWT 并返回解析后的 claims。
func (s *JWTService) ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}

// ExpireDuration 返回指定 token 类型的过期时长。
func (s *JWTService) ExpireDuration(tokenType TokenType) time.Duration {
	if tokenType == TokenTypeAdmin {
		return s.adminExpire
	}
	return s.walletExpire
}
