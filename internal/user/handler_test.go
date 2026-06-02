package user

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mousecake-go/mousecake-go/internal/shared/auth"
	"github.com/mousecake-go/mousecake-go/internal/user/domain"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupHandler() (*Handler, *mockUserRepo, *gin.Engine) {
	repo := newMockUserRepo()
	jwtSvc := auth.NewJWTService("test-secret", 4*time.Hour, 8*time.Hour)
	svc := NewService(repo, jwtSvc, []int{1, 5, 11155111}, "admin", "admin123456", "mousecake-go")
	handler := NewHandler(svc, jwtSvc)

	r := gin.New()
	handler.RegisterRoutes(r.Group("/api/v1/auth"))
	return handler, repo, r
}

// 9.1 POST /api/v1/auth/wallet/nonce
func TestHandler_WalletNonce_Success(t *testing.T) {
	_, _, r := setupHandler()

	body, _ := json.Marshal(map[string]string{"address": "0x1234567890abcdef1234567890abcdef12345678"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/wallet/nonce", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(0), resp["code"])
	data := resp["data"].(map[string]interface{})
	assert.NotEmpty(t, data["message"])
	assert.NotEmpty(t, data["nonce"])
}

func TestHandler_WalletNonce_InvalidAddress(t *testing.T) {
	_, _, r := setupHandler()

	body, _ := json.Marshal(map[string]string{"address": "invalid"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/wallet/nonce", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestHandler_WalletNonce_MissingAddress(t *testing.T) {
	_, _, r := setupHandler()

	body, _ := json.Marshal(map[string]string{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/wallet/nonce", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

// 9.5 POST /api/v1/auth/admin/login
func TestHandler_AdminLogin_Success(t *testing.T) {
	_, repo, r := setupHandler()

	admin, _ := domain.NewAdminUser("admin", mustHashPassword("password123"))
	admin.ID = 1
	repo.usersByName["admin"] = admin
	repo.usersByID[1] = admin

	body, _ := json.Marshal(map[string]string{"username": "admin", "password": "password123"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/admin/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	assert.NotEmpty(t, data["access_token"])
	assert.Equal(t, "Bearer", data["token_type"])
}

func TestHandler_AdminLogin_BadCredentials(t *testing.T) {
	_, _, r := setupHandler()

	body, _ := json.Marshal(map[string]string{"username": "admin", "password": "wrongpassword"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/admin/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

// 9.7 GET /api/v1/auth/me
func TestHandler_GetMe_WalletUser(t *testing.T) {
	_, repo, r := setupHandler()

	walletUser, _ := domain.NewWalletUser("0x1234567890abcdef1234567890abcdef12345678")
	walletUser.ID = 1
	repo.users["0x1234567890abcdef1234567890abcdef12345678"] = walletUser
	repo.usersByID[1] = walletUser

	jwtSvc := auth.NewJWTService("test-secret", 4*time.Hour, 8*time.Hour)
	token, _ := jwtSvc.IssueToken(auth.TokenTypeWallet, 1, false)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "wallet", data["type"])
}

func TestHandler_GetMe_NoToken(t *testing.T) {
	_, _, r := setupHandler()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/auth/me", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

// W3: Handler 层补充测试

func TestHandler_WalletVerify_MissingFields(t *testing.T) {
	_, _, r := setupHandler()

	body, _ := json.Marshal(map[string]string{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/wallet/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(40001), resp["code"])
}

func TestHandler_AdminLogin_Disabled(t *testing.T) {
	_, repo, r := setupHandler()

	admin, _ := domain.NewAdminUser("disabledadmin", mustHashPassword("password123"))
	admin.ID = 1
	admin.Disable()
	repo.usersByName["disabledadmin"] = admin
	repo.usersByID[1] = admin

	body, _ := json.Marshal(map[string]string{"username": "disabledadmin", "password": "password123"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/admin/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(40104), resp["code"])
}

func TestHandler_AdminLogin_ShortPassword(t *testing.T) {
	_, _, r := setupHandler()

	body, _ := json.Marshal(map[string]string{"username": "admin", "password": "short"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/admin/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
}

func TestHandler_GetMe_Admin(t *testing.T) {
	_, repo, r := setupHandler()

	admin, _ := domain.NewAdminUser("admin", "hash")
	admin.ID = 2
	repo.usersByName["admin"] = admin
	repo.usersByID[2] = admin

	jwtSvc := auth.NewJWTService("test-secret", 4*time.Hour, 8*time.Hour)
	token, _ := jwtSvc.IssueToken(auth.TokenTypeAdmin, 2, true)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, "admin", data["type"])
	assert.Equal(t, true, data["is_admin"])
}

func TestHandler_GetMe_ExpiredToken(t *testing.T) {
	_, _, r := setupHandler()

	jwtSvc := auth.NewJWTService("test-secret", 1*time.Nanosecond, 1*time.Nanosecond)
	token, _ := jwtSvc.IssueToken(auth.TokenTypeWallet, 1, false)

	time.Sleep(10 * time.Millisecond)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(40105), resp["code"])
}

// === Issue 4: ExpiresIn 应匹配 JWT 配置 ===

func TestHandler_WalletVerify_ExpiresInMatchesConfig(t *testing.T) {
	walletExpire := 2 * time.Hour // 7200 秒
	repo := newMockUserRepo()
	jwtSvc := auth.NewJWTService("test-secret", walletExpire, 8*time.Hour)
	svc := NewService(repo, jwtSvc, []int{1, 5, 11155111}, "admin", "admin123456", "mousecake-go")
	handler := NewHandler(svc, jwtSvc)

	privateKey, address := generateTestKey()
	msg, _, err := svc.GenerateSIWENonce(context.Background(), address)
	require.NoError(t, err)

	signature := signTestMessage(msg, privateKey)

	r := gin.New()
	handler.RegisterRoutes(r.Group("/api/v1/auth"))

	body, _ := json.Marshal(map[string]string{"message": msg, "signature": signature})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/wallet/verify", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(7200), data["expires_in"], "wallet ExpiresIn 应等于 wallet_expire 配置（秒）")
}

func TestHandler_AdminLogin_ExpiresInMatchesConfig(t *testing.T) {
	adminExpire := 4 * time.Hour // 14400 秒
	repo := newMockUserRepo()
	jwtSvc := auth.NewJWTService("test-secret", 4*time.Hour, adminExpire)
	svc := NewService(repo, jwtSvc, []int{1, 5, 11155111}, "admin", "admin123456", "mousecake-go")
	handler := NewHandler(svc, jwtSvc)

	admin, _ := domain.NewAdminUser("admin", mustHashPassword("password123"))
	admin.ID = 1
	repo.usersByName["admin"] = admin
	repo.usersByID[1] = admin

	r := gin.New()
	handler.RegisterRoutes(r.Group("/api/v1/auth"))

	body, _ := json.Marshal(map[string]string{"username": "admin", "password": "password123"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/auth/admin/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	assert.Equal(t, float64(14400), data["expires_in"], "admin ExpiresIn 应等于 admin_expire 配置（秒）")
}
