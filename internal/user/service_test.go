package user

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mousecake-go/mousecake-go/internal/shared/auth"
	"github.com/mousecake-go/mousecake-go/internal/shared/errs"
	"github.com/mousecake-go/mousecake-go/internal/user/domain"
)

// mockUserRepo 是 domain.UserRepository 的 mock 实现。
type mockUserRepo struct {
	users          map[string]*domain.User       // address -> user
	usersByName    map[string]*domain.User       // username -> user
	usersByID      map[int64]*domain.User        // id -> user
	nonces         map[string]*domain.LoginNonce // address -> nonce
	nextID         int64
	deleteNonceErr error
	updateErr      error
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users:       make(map[string]*domain.User),
		usersByName: make(map[string]*domain.User),
		usersByID:   make(map[int64]*domain.User),
		nonces:      make(map[string]*domain.LoginNonce),
		nextID:      1,
	}
}

func (m *mockUserRepo) FindByAddress(_ context.Context, address string) (*domain.User, error) {
	u, ok := m.users[address]
	if !ok || u.DeletedAt != nil {
		return nil, fmt.Errorf("not found")
	}
	return u, nil
}

func (m *mockUserRepo) FindByUsername(_ context.Context, username string) (*domain.User, error) {
	u, ok := m.usersByName[username]
	if !ok || u.DeletedAt != nil {
		return nil, fmt.Errorf("not found")
	}
	return u, nil
}

func (m *mockUserRepo) FindByID(_ context.Context, id int64) (*domain.User, error) {
	u, ok := m.usersByID[id]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return u, nil
}

func (m *mockUserRepo) Create(_ context.Context, user *domain.User) error {
	user.ID = m.nextID
	m.nextID++
	if user.Address != "" {
		m.users[user.Address] = user
	}
	if user.Username != "" {
		m.usersByName[user.Username] = user
	}
	m.usersByID[user.ID] = user
	return nil
}

func (m *mockUserRepo) Update(_ context.Context, user *domain.User) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.usersByID[user.ID] = user
	if user.Address != "" {
		m.users[user.Address] = user
	}
	if user.Username != "" {
		m.usersByName[user.Username] = user
	}
	return nil
}

func (m *mockUserRepo) UpsertByAddress(ctx context.Context, user *domain.User) error {
	existing, ok := m.users[user.Address]
	if ok {
		now := time.Now()
		existing.LastLoginAt = &now
		return nil
	}
	return m.Create(ctx, user)
}

func (m *mockUserRepo) SaveNonce(_ context.Context, nonce *domain.LoginNonce) error {
	m.nonces[nonce.Address] = nonce
	return nil
}

func (m *mockUserRepo) FindNonceByAddress(_ context.Context, address string) (*domain.LoginNonce, error) {
	n, ok := m.nonces[address]
	if !ok {
		return nil, fmt.Errorf("not found")
	}
	return n, nil
}

func (m *mockUserRepo) DeleteNonce(_ context.Context, address string) error {
	if m.deleteNonceErr != nil {
		return m.deleteNonceErr
	}
	delete(m.nonces, address)
	return nil
}

func (m *mockUserRepo) CleanExpiredNonces(_ context.Context) error {
	for addr, n := range m.nonces {
		if n.IsExpired() {
			delete(m.nonces, addr)
		}
	}
	return nil
}

func newTestService(repo *mockUserRepo) *Service {
	jwtSvc := auth.NewJWTService("test-secret", 4*time.Hour, 8*time.Hour)
	return NewService(repo, jwtSvc, []int{1, 5, 11155111}, "admin", "admin123456", "mousecake-go")
}

// 8.1 GenerateSIWENonce
func TestService_GenerateSIWENonce_First(t *testing.T) {
	repo := newMockUserRepo()
	svc := newTestService(repo)

	msg, nonce, err := svc.GenerateSIWENonce(context.Background(), "0x1234567890abcdef1234567890abcdef12345678")
	require.NoError(t, err)
	assert.NotEmpty(t, msg)
	assert.NotEmpty(t, nonce)
}

func TestService_GenerateSIWENonce_Override(t *testing.T) {
	repo := newMockUserRepo()
	svc := newTestService(repo)

	_, _, _ = svc.GenerateSIWENonce(context.Background(), "0x1234567890abcdef1234567890abcdef12345678")
	_, newNonce, err := svc.GenerateSIWENonce(context.Background(), "0x1234567890abcdef1234567890abcdef12345678")
	require.NoError(t, err)
	assert.NotEmpty(t, newNonce)
}

func TestService_GenerateSIWENonce_InvalidAddress(t *testing.T) {
	repo := newMockUserRepo()
	svc := newTestService(repo)

	_, _, err := svc.GenerateSIWENonce(context.Background(), "invalid")
	assert.Error(t, err)
}

// 8.5 AdminLogin
func TestService_AdminLogin_Success(t *testing.T) {
	repo := newMockUserRepo()
	svc := newTestService(repo)

	// 先 seed 一个管理员
	admin, _ := domain.NewAdminUser("admin", mustHashPassword("password123"))
	repo.usersByName["admin"] = admin
	repo.usersByID[1] = admin
	admin.ID = 1

	token, err := svc.AdminLogin(context.Background(), "admin", "password123")
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestService_AdminLogin_UserNotFound(t *testing.T) {
	repo := newMockUserRepo()
	svc := newTestService(repo)

	_, err := svc.AdminLogin(context.Background(), "nonexistent", "password123")
	assert.Error(t, err)
}

func TestService_AdminLogin_WrongPassword(t *testing.T) {
	repo := newMockUserRepo()
	svc := newTestService(repo)

	admin, _ := domain.NewAdminUser("admin", mustHashPassword("password123"))
	repo.usersByName["admin"] = admin

	_, err := svc.AdminLogin(context.Background(), "admin", "wrong-password")
	assert.Error(t, err)
}

func TestService_AdminLogin_EmptyPassword(t *testing.T) {
	repo := newMockUserRepo()
	svc := newTestService(repo)

	_, err := svc.AdminLogin(context.Background(), "admin", "")
	assert.Error(t, err)
}

// 8.7 SeedAdmin
func TestService_SeedAdmin_First(t *testing.T) {
	repo := newMockUserRepo()
	svc := newTestService(repo)

	err := svc.SeedAdmin(context.Background())
	require.NoError(t, err)

	admin, err := repo.FindByUsername(context.Background(), "admin")
	require.NoError(t, err)
	assert.True(t, admin.IsAdmin)
}

func TestService_SeedAdmin_AlreadyExists(t *testing.T) {
	repo := newMockUserRepo()
	svc := newTestService(repo)

	existing, _ := domain.NewAdminUser("admin", "existing_hash")
	existing.ID = 1
	repo.usersByName["admin"] = existing
	repo.usersByID[1] = existing

	err := svc.SeedAdmin(context.Background())
	require.NoError(t, err)
}

// 8.9 GetCurrentUser
func TestService_GetCurrentUser_WalletUser(t *testing.T) {
	repo := newMockUserRepo()
	svc := newTestService(repo)

	walletUser, _ := domain.NewWalletUser("0x1234567890abcdef1234567890abcdef12345678")
	walletUser.ID = 1
	repo.users["0x1234567890abcdef1234567890abcdef12345678"] = walletUser
	repo.usersByID[1] = walletUser

	user, err := svc.GetCurrentUser(context.Background(), 1)
	require.NoError(t, err)
	assert.Equal(t, int64(1), user.ID)
}

func TestService_GetCurrentUser_Admin(t *testing.T) {
	repo := newMockUserRepo()
	svc := newTestService(repo)

	admin, _ := domain.NewAdminUser("admin", "hash")
	admin.ID = 2
	repo.usersByName["admin"] = admin
	repo.usersByID[2] = admin

	user, err := svc.GetCurrentUser(context.Background(), 2)
	require.NoError(t, err)
	assert.True(t, user.IsAdmin)
}

func TestService_CleanExpiredNonces(t *testing.T) {
	repo := newMockUserRepo()
	svc := newTestService(repo)

	err := svc.CleanExpiredNonces(context.Background())
	require.NoError(t, err)
}

func mustHashPassword(password string) string {
	hash, err := hashPassword(password)
	if err != nil {
		panic(err)
	}
	return hash
}

func generateTestKey() (*ecdsa.PrivateKey, string) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	address := strings.ToLower(crypto.PubkeyToAddress(privateKey.PublicKey).Hex())
	return privateKey, address
}

func signTestMessage(message string, privateKey *ecdsa.PrivateKey) string {
	hash := crypto.Keccak256([]byte(message))
	sig, err := crypto.Sign(hash, privateKey)
	if err != nil {
		panic(err)
	}
	sig[64] += 27
	return "0x" + hex.EncodeToString(sig)
}

// W1: SIWE 签名验证场景测试

func TestService_VerifySIWESignature_NewUser(t *testing.T) {
	repo := newMockUserRepo()
	svc := newTestService(repo)

	privateKey, address := generateTestKey()
	msg, _, err := svc.GenerateSIWENonce(context.Background(), address)
	require.NoError(t, err)

	signature := signTestMessage(msg, privateKey)
	token, err := svc.VerifySIWESignature(context.Background(), msg, signature)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	user, err := repo.FindByAddress(context.Background(), address)
	require.NoError(t, err)
	assert.Equal(t, address, user.Address)
}

func TestService_VerifySIWESignature_ExistingUser(t *testing.T) {
	repo := newMockUserRepo()
	svc := newTestService(repo)

	privateKey, address := generateTestKey()

	existing, _ := domain.NewWalletUser(address)
	existing.ID = 1
	repo.users[address] = existing
	repo.usersByID[1] = existing

	msg, _, err := svc.GenerateSIWENonce(context.Background(), address)
	require.NoError(t, err)

	signature := signTestMessage(msg, privateKey)
	token, err := svc.VerifySIWESignature(context.Background(), msg, signature)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	updated, _ := repo.FindByAddress(context.Background(), address)
	assert.NotNil(t, updated.LastLoginAt)
}

func TestService_VerifySIWESignature_NonceNotFound(t *testing.T) {
	repo := newMockUserRepo()
	svc := newTestService(repo)

	privateKey, address := generateTestKey()

	msg := fmt.Sprintf(
		"mousecake-go wants you to sign in with your Ethereum account:\n%s\n\nI accept the terms of service\n\nURI: https://mousecake.io\nNonce: nonexist\nChain ID: 1\nIssued At: %s",
		address, time.Now().UTC().Format(time.RFC3339),
	)
	signature := signTestMessage(msg, privateKey)

	_, err := svc.VerifySIWESignature(context.Background(), msg, signature)
	assert.Error(t, err)
	assertErrorCode(t, err, errs.CodeNonceExpired)
}

func TestService_VerifySIWESignature_SignatureMismatch(t *testing.T) {
	repo := newMockUserRepo()
	svc := newTestService(repo)

	_, addressA := generateTestKey()
	msg, _, err := svc.GenerateSIWENonce(context.Background(), addressA)
	require.NoError(t, err)

	privateKeyB, _ := generateTestKey()
	signature := signTestMessage(msg, privateKeyB)

	_, err = svc.VerifySIWESignature(context.Background(), msg, signature)
	assert.Error(t, err)
	assertErrorCode(t, err, errs.CodeSignInvalid)
}

func TestService_VerifySIWESignature_InvalidMessage(t *testing.T) {
	repo := newMockUserRepo()
	svc := newTestService(repo)

	_, err := svc.VerifySIWESignature(context.Background(), "invalid message", "0x00")
	assert.Error(t, err)
	assertErrorCode(t, err, errs.CodeSIWEFormat)
}

func TestService_VerifySIWESignature_UserDisabled(t *testing.T) {
	repo := newMockUserRepo()
	svc := newTestService(repo)

	privateKey, address := generateTestKey()

	disabledUser, _ := domain.NewWalletUser(address)
	disabledUser.ID = 1
	disabledUser.Disable()
	repo.users[address] = disabledUser
	repo.usersByID[1] = disabledUser

	msg, _, err := svc.GenerateSIWENonce(context.Background(), address)
	require.NoError(t, err)

	signature := signTestMessage(msg, privateKey)
	_, err = svc.VerifySIWESignature(context.Background(), msg, signature)
	assert.Error(t, err)
	assertErrorCode(t, err, errs.CodeAccountDisabled)
}

func TestService_VerifySIWESignature_NonceOneTimeUse(t *testing.T) {
	repo := newMockUserRepo()
	svc := newTestService(repo)

	privateKey, address := generateTestKey()
	msg, _, err := svc.GenerateSIWENonce(context.Background(), address)
	require.NoError(t, err)

	signature := signTestMessage(msg, privateKey)

	_, err = svc.VerifySIWESignature(context.Background(), msg, signature)
	require.NoError(t, err)

	_, err = svc.VerifySIWESignature(context.Background(), msg, signature)
	assert.Error(t, err)
	assertErrorCode(t, err, errs.CodeNonceExpired)
}

// W4: chain_id 校验

func TestService_VerifySIWESignature_UnsupportedChainID(t *testing.T) {
	repo := newMockUserRepo()
	svc := newTestService(repo)

	privateKey, address := generateTestKey()

	msg := fmt.Sprintf(
		"mousecake-go wants you to sign in with your Ethereum account:\n%s\n\nI accept the terms of service\n\nURI: https://mousecake.io\nNonce: testnonce\nChain ID: 999\nIssued At: %s",
		address, time.Now().UTC().Format(time.RFC3339),
	)
	signature := signTestMessage(msg, privateKey)

	_, err := svc.VerifySIWESignature(context.Background(), msg, signature)
	assert.Error(t, err)
	assertErrorCode(t, err, errs.CodeChainUnsupported)
}

// W2: 管理员登录边界场景

func TestService_AdminLogin_NonAdmin(t *testing.T) {
	repo := newMockUserRepo()
	svc := newTestService(repo)

	walletUser, _ := domain.NewWalletUser("0x1234567890abcdef1234567890abcdef12345678")
	walletUser.ID = 1
	repo.usersByName["walletuser"] = walletUser
	repo.usersByID[1] = walletUser

	_, err := svc.AdminLogin(context.Background(), "walletuser", "password123")
	assert.Error(t, err)
	assertErrorCode(t, err, errs.CodeCredentialBad)
}

func TestService_AdminLogin_Disabled(t *testing.T) {
	repo := newMockUserRepo()
	svc := newTestService(repo)

	admin, _ := domain.NewAdminUser("disabledadmin", mustHashPassword("password123"))
	admin.ID = 1
	admin.Disable()
	repo.usersByName["disabledadmin"] = admin
	repo.usersByID[1] = admin

	_, err := svc.AdminLogin(context.Background(), "disabledadmin", "password123")
	assert.Error(t, err)
	assertErrorCode(t, err, errs.CodeAccountDisabled)
}

// assertErrorCode 验证错误是否为 serviceError 且码值匹配。
func assertErrorCode(t *testing.T, err error, expectedCode int) {
	t.Helper()
	var svcErr *serviceError
	if !assert.True(t, errors.As(err, &svcErr), "expected *serviceError, got %T: %v", err, err) {
		return
	}
	assert.Equal(t, expectedCode, svcErr.Code())
}

// === Issue 1: 忽略错误返回值 ===

func TestService_VerifySIWESignature_DeleteNonceError(t *testing.T) {
	repo := newMockUserRepo()
	repo.deleteNonceErr = fmt.Errorf("db connection lost")
	svc := newTestService(repo)

	privateKey, address := generateTestKey()
	msg, _, err := svc.GenerateSIWENonce(context.Background(), address)
	require.NoError(t, err)

	signature := signTestMessage(msg, privateKey)
	_, err = svc.VerifySIWESignature(context.Background(), msg, signature)
	assert.Error(t, err, "DeleteNonce 失败应返回错误")
	assert.Contains(t, err.Error(), "delete nonce")
}

func TestService_VerifySIWESignature_UpdateError(t *testing.T) {
	repo := newMockUserRepo()
	repo.updateErr = fmt.Errorf("db connection lost")
	svc := newTestService(repo)

	privateKey, address := generateTestKey()
	msg, _, err := svc.GenerateSIWENonce(context.Background(), address)
	require.NoError(t, err)

	signature := signTestMessage(msg, privateKey)
	_, err = svc.VerifySIWESignature(context.Background(), msg, signature)
	assert.Error(t, err, "Update LastLoginAt 失败应返回错误")
	assert.Contains(t, err.Error(), "update last login")
}

func TestService_AdminLogin_UpdateError(t *testing.T) {
	repo := newMockUserRepo()
	repo.updateErr = fmt.Errorf("db connection lost")
	svc := newTestService(repo)

	admin, _ := domain.NewAdminUser("admin", mustHashPassword("password123"))
	admin.ID = 1
	repo.usersByName["admin"] = admin
	repo.usersByID[1] = admin

	_, err := svc.AdminLogin(context.Background(), "admin", "password123")
	assert.Error(t, err, "Update LastLoginAt 失败应返回错误")
	assert.Contains(t, err.Error(), "update last login")
}

// === Issue 3: typed error 提取 ===

func TestService_TypedErrorCodeExtraction(t *testing.T) {
	repo := newMockUserRepo()
	svc := newTestService(repo)

	_, err := svc.VerifySIWESignature(context.Background(), "invalid message", "0x00")
	require.Error(t, err)

	var svcErr *serviceError
	assert.True(t, errors.As(err, &svcErr), "应能通过 errors.As 提取 serviceError")
	if svcErr != nil {
		assert.Equal(t, errs.CodeSIWEFormat, svcErr.Code())
	}
}

// === Issue 5: SIWE domain 和 issued-at 验证 ===

func TestService_VerifySIWESignature_DomainMismatch(t *testing.T) {
	repo := newMockUserRepo()
	svc := newTestService(repo)
	privateKey, address := generateTestKey()

	msg := fmt.Sprintf(
		"evil.com wants you to sign in with your Ethereum account:\n%s\n\nI accept the terms of service\n\nURI: https://evil.com\nNonce: testnonce\nChain ID: 1\nIssued At: %s",
		address, time.Now().UTC().Format(time.RFC3339),
	)

	nonce, _ := domain.NewLoginNonce(address, 5*time.Minute)
	nonce.Nonce = "testnonce"
	repo.nonces[address] = nonce

	signature := signTestMessage(msg, privateKey)
	_, err := svc.VerifySIWESignature(context.Background(), msg, signature)
	assert.Error(t, err, "域名为不匹配时应返回错误")
	assertErrorCode(t, err, errs.CodeSIWEFormat)
}

func TestService_VerifySIWESignature_StaleIssuedAt(t *testing.T) {
	repo := newMockUserRepo()
	svc := newTestService(repo)
	privateKey, address := generateTestKey()

	staleTime := time.Now().Add(-30 * time.Minute).UTC().Format(time.RFC3339)
	msg := fmt.Sprintf(
		"mousecake-go wants you to sign in with your Ethereum account:\n%s\n\nI accept the terms of service\n\nURI: https://mousecake.io\nNonce: testnonce\nChain ID: 1\nIssued At: %s",
		address, staleTime,
	)

	nonce, _ := domain.NewLoginNonce(address, 5*time.Minute)
	nonce.Nonce = "testnonce"
	repo.nonces[address] = nonce

	signature := signTestMessage(msg, privateKey)
	_, err := svc.VerifySIWESignature(context.Background(), msg, signature)
	assert.Error(t, err, "issued-at 过旧时应返回错误")
	assertErrorCode(t, err, errs.CodeSIWEFormat)
}

func TestService_VerifySIWESignature_MissingIssuedAt(t *testing.T) {
	repo := newMockUserRepo()
	svc := newTestService(repo)
	privateKey, address := generateTestKey()

	msg := fmt.Sprintf(
		"mousecake-go wants you to sign in with your Ethereum account:\n%s\n\nI accept the terms of service\n\nURI: https://mousecake.io\nNonce: testnonce\nChain ID: 1\n",
		address,
	)

	nonce, _ := domain.NewLoginNonce(address, 5*time.Minute)
	nonce.Nonce = "testnonce"
	repo.nonces[address] = nonce

	signature := signTestMessage(msg, privateKey)
	_, err := svc.VerifySIWESignature(context.Background(), msg, signature)
	assert.Error(t, err, "缺少 Issued At 应返回错误")
	assertErrorCode(t, err, errs.CodeSIWEFormat)
}
