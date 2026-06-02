// Package user 实现用户认证模块的业务逻辑和数据访问。
package user

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/mousecake-go/mousecake-go/internal/shared/auth"
	"github.com/mousecake-go/mousecake-go/internal/shared/errs"
	"github.com/mousecake-go/mousecake-go/internal/user/domain"
)

// ethAddressRegex 以太坊地址格式校验。
var ethAddressRegex = regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)

// siweIssuedAtMaxAge 是 SIWE 消息 issued-at 字段允许的最大年龄。
const siweIssuedAtMaxAge = 10 * time.Minute

// serviceError 是携带业务错误码的服务层错误类型。
type serviceError struct {
	code int
	msg  string
}

// Error 实现 error 接口。
func (e *serviceError) Error() string { return e.msg }

// Code 返回业务错误码。
func (e *serviceError) Code() int { return e.code }

// newServiceError 创建业务错误。
func newServiceError(code int, msg string) *serviceError {
	return &serviceError{code: code, msg: msg}
}

// Service 编排用户认证相关的用例。
type Service struct {
	repo            domain.UserRepository
	jwtSvc          *auth.JWTService
	allowedChainIDs []int
	adminUsername   string
	adminPassword   string
	siweDomain      string
	log             *slog.Logger
}

// NewService 创建用户服务实例。
func NewService(repo domain.UserRepository, jwtSvc *auth.JWTService, allowedChainIDs []int, adminUsername, adminPassword, siweDomain string) *Service {
	return &Service{
		repo:            repo,
		jwtSvc:          jwtSvc,
		allowedChainIDs: allowedChainIDs,
		adminUsername:   adminUsername,
		adminPassword:   adminPassword,
		siweDomain:      siweDomain,
		log:             slog.Default().With("module", "user", "layer", "service"),
	}
}

// GenerateSIWENonce 生成 SIWE nonce 并构造 EIP-4361 消息。
func (s *Service) GenerateSIWENonce(ctx context.Context, address string) (string, string, error) {
	if !ethAddressRegex.MatchString(address) {
		return "", "", fmt.Errorf("invalid address format: %s", address)
	}

	address = strings.ToLower(address)
	nonce, err := domain.NewLoginNonce(address, 5*time.Minute)
	if err != nil {
		return "", "", fmt.Errorf("generate nonce: %w", err)
	}

	if err := s.repo.SaveNonce(ctx, nonce); err != nil {
		return "", "", fmt.Errorf("save nonce: %w", err)
	}

	if len(s.allowedChainIDs) == 0 {
		return "", "", fmt.Errorf("allowed chain IDs is not configured")
	}
	chainID := s.allowedChainIDs[0]

	msg := buildSIWEMessage(address, nonce.Nonce, s.siweDomain, chainID)
	return msg, nonce.Nonce, nil
}

// VerifySIWESignature 验证 SIWE 签名，自动注册新用户，签发 JWT。
func (s *Service) VerifySIWESignature(ctx context.Context, message, signature string) (string, error) {
	address, nonceStr, chainID, err := parseSIWEMessage(message, s.siweDomain)
	if err != nil {
		return "", err
	}

	if !s.isChainAllowed(chainID) {
		return "", newServiceError(errs.CodeChainUnsupported, errs.GetErrorMessage(errs.CodeChainUnsupported))
	}

	storedNonce, err := s.repo.FindNonceByAddress(ctx, address)
	if err != nil {
		return "", newServiceError(errs.CodeNonceExpired, errs.GetErrorMessage(errs.CodeNonceExpired))
	}

	if storedNonce.Nonce != nonceStr {
		return "", newServiceError(errs.CodeNonceExpired, errs.GetErrorMessage(errs.CodeNonceExpired))
	}

	if storedNonce.IsExpired() {
		return "", newServiceError(errs.CodeNonceExpired, errs.GetErrorMessage(errs.CodeNonceExpired))
	}

	recoveredAddr, err := recoverAddress(message, signature)
	if err != nil {
		return "", newServiceError(errs.CodeSignInvalid, errs.GetErrorMessage(errs.CodeSignInvalid))
	}

	if !strings.EqualFold(recoveredAddr, address) {
		return "", newServiceError(errs.CodeSignInvalid, errs.GetErrorMessage(errs.CodeSignInvalid))
	}

	if err := s.repo.DeleteNonce(ctx, address); err != nil {
		return "", fmt.Errorf("delete nonce: %w", err)
	}

	user, err := s.repo.FindByAddress(ctx, address)
	if err != nil {
		user, err = domain.NewWalletUser(address)
		if err != nil {
			return "", fmt.Errorf("create wallet user: %w", err)
		}
		if err := s.repo.Create(ctx, user); err != nil {
			return "", fmt.Errorf("create user: %w", err)
		}
	}

	if user.Status == domain.UserStatusDisabled {
		return "", newServiceError(errs.CodeAccountDisabled, errs.GetErrorMessage(errs.CodeAccountDisabled))
	}

	now := time.Now()
	user.LastLoginAt = &now
	if err := s.repo.Update(ctx, user); err != nil {
		return "", fmt.Errorf("update last login: %w", err)
	}

	token, err := s.jwtSvc.IssueToken(auth.TokenTypeWallet, user.ID, user.IsAdmin)
	if err != nil {
		return "", fmt.Errorf("issue token: %w", err)
	}

	return token, nil
}

// AdminLogin 管理员密码登录。
func (s *Service) AdminLogin(ctx context.Context, username, password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password is required")
	}

	user, err := s.repo.FindByUsername(ctx, username)
	if err != nil {
		return "", newServiceError(errs.CodeCredentialBad, errs.GetErrorMessage(errs.CodeCredentialBad))
	}

	if !user.IsAdmin {
		return "", newServiceError(errs.CodeCredentialBad, errs.GetErrorMessage(errs.CodeCredentialBad))
	}

	if user.Status == domain.UserStatusDisabled {
		return "", newServiceError(errs.CodeAccountDisabled, errs.GetErrorMessage(errs.CodeAccountDisabled))
	}

	match, err := argon2id.ComparePasswordAndHash(password, user.PasswordHash)
	if err != nil || !match {
		return "", newServiceError(errs.CodeCredentialBad, errs.GetErrorMessage(errs.CodeCredentialBad))
	}

	now := time.Now()
	user.LastLoginAt = &now
	if err := s.repo.Update(ctx, user); err != nil {
		return "", fmt.Errorf("update last login: %w", err)
	}

	token, err := s.jwtSvc.IssueToken(auth.TokenTypeAdmin, user.ID, true)
	if err != nil {
		return "", fmt.Errorf("issue token: %w", err)
	}

	return token, nil
}

// CreateAdmin 通过 CLI 创建管理员账号。
func (s *Service) CreateAdmin(ctx context.Context, username, password string) error {
	if len(password) < 8 {
		return fmt.Errorf("密码长度不足（最少 8 位）")
	}

	_, err := s.repo.FindByUsername(ctx, username)
	if err == nil {
		return fmt.Errorf("用户名 %s 已存在", username)
	}

	hash, err := hashPassword(password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	admin, err := domain.NewAdminUser(username, hash)
	if err != nil {
		return fmt.Errorf("create admin user: %w", err)
	}

	return s.repo.Create(ctx, admin)
}

// SeedAdmin 初始化默认管理员账号。
func (s *Service) SeedAdmin(ctx context.Context) error {
	existing, err := s.repo.FindByUsername(ctx, s.adminUsername)
	if err == nil && existing != nil {
		return nil
	}

	hash, err := hashPassword(s.adminPassword)
	if err != nil {
		return fmt.Errorf("hash admin password: %w", err)
	}

	admin, err := domain.NewAdminUser(s.adminUsername, hash)
	if err != nil {
		return fmt.Errorf("create admin user: %w", err)
	}

	if err := s.repo.Create(ctx, admin); err != nil {
		return fmt.Errorf("create admin: %w", err)
	}

	return nil
}

// GetCurrentUser 获取当前用户信息。
func (s *Service) GetCurrentUser(ctx context.Context, userID int64) (*domain.User, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find user %d: %w", userID, err)
	}
	return user, nil
}

// CleanExpiredNonces 清理过期的 nonce 记录。
func (s *Service) CleanExpiredNonces(ctx context.Context) error {
	return s.repo.CleanExpiredNonces(ctx)
}

// isChainAllowed 检查 chain ID 是否在允许列表中。
func (s *Service) isChainAllowed(chainID int) bool {
	for _, id := range s.allowedChainIDs {
		if id == chainID {
			return true
		}
	}
	return false
}

func buildSIWEMessage(address, nonce, domain string, chainID int) string {
	return fmt.Sprintf(
		"%s wants you to sign in with your Ethereum account:\n%s\n\nI accept the terms of service\n\nURI: https://mousecake.io\nNonce: %s\nChain ID: %d\nIssued At: %s",
		domain, address, nonce, chainID, time.Now().UTC().Format(time.RFC3339),
	)
}

func parseSIWEMessage(message, expectedDomain string) (string, string, int, error) {
	lines := strings.Split(message, "\n")
	var address, nonce, domain, issuedAtStr string
	var chainID int

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "Nonce: "):
			nonce = strings.TrimPrefix(line, "Nonce: ")
		case strings.HasPrefix(line, "Chain ID: "):
			idStr := strings.TrimPrefix(line, "Chain ID: ")
			if v, err := strconv.Atoi(idStr); err == nil {
				chainID = v
			}
		case strings.HasPrefix(line, "Issued At: "):
			issuedAtStr = strings.TrimPrefix(line, "Issued At: ")
		}
	}

	if len(lines) >= 1 {
		parts := strings.SplitN(lines[0], " wants you to sign in", 2)
		if len(parts) == 2 {
			domain = parts[0]
		}
	}

	if len(lines) >= 2 {
		address = strings.TrimSpace(lines[1])
	}

	if address == "" || nonce == "" {
		return "", "", 0, newServiceError(errs.CodeSIWEFormat, errs.GetErrorMessage(errs.CodeSIWEFormat))
	}

	if domain != expectedDomain {
		return "", "", 0, newServiceError(errs.CodeSIWEFormat, errs.GetErrorMessage(errs.CodeSIWEFormat))
	}

	if issuedAtStr == "" {
		return "", "", 0, newServiceError(errs.CodeSIWEFormat, errs.GetErrorMessage(errs.CodeSIWEFormat))
	}

	issuedAt, err := time.Parse(time.RFC3339, issuedAtStr)
	if err != nil {
		return "", "", 0, newServiceError(errs.CodeSIWEFormat, errs.GetErrorMessage(errs.CodeSIWEFormat))
	}

	if time.Since(issuedAt) > siweIssuedAtMaxAge {
		return "", "", 0, newServiceError(errs.CodeSIWEFormat, errs.GetErrorMessage(errs.CodeSIWEFormat))
	}

	return strings.ToLower(address), nonce, chainID, nil
}

func recoverAddress(message, signature string) (string, error) {
	sig, err := decodeSignature(signature)
	if err != nil {
		return "", fmt.Errorf("decode signature: %w", err)
	}

	pubKey, err := crypto.SigToPub(crypto.Keccak256([]byte(message)), sig)
	if err != nil {
		return "", fmt.Errorf("recover public key: %w", err)
	}

	addr := crypto.PubkeyToAddress(*pubKey)
	return strings.ToLower(addr.Hex()), nil
}

func decodeSignature(sig string) ([]byte, error) {
	sig = strings.TrimPrefix(sig, "0x")
	bytes := make([]byte, len(sig)/2)
	for i := range bytes {
		_, err := fmt.Sscanf(sig[i*2:i*2+2], "%02x", &bytes[i])
		if err != nil {
			return nil, err
		}
	}
	if len(bytes) >= 65 {
		bytes[64] -= 27
	}
	return bytes, nil
}

func hashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return hash, nil
}
