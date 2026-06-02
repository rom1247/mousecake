package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJWTService_IssueAndValidate_WalletToken(t *testing.T) {
	svc := NewJWTService("test-secret", 4*time.Hour, 8*time.Hour)

	token, err := svc.IssueToken(TokenTypeWallet, 1, false)
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := svc.ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, TokenTypeWallet, claims.Type)
	assert.Equal(t, int64(1), claims.UserID)
	assert.False(t, claims.IsAdmin)
}

func TestJWTService_IssueAndValidate_AdminToken(t *testing.T) {
	svc := NewJWTService("test-secret", 4*time.Hour, 8*time.Hour)

	token, err := svc.IssueToken(TokenTypeAdmin, 2, true)
	require.NoError(t, err)

	claims, err := svc.ValidateToken(token)
	require.NoError(t, err)
	assert.Equal(t, TokenTypeAdmin, claims.Type)
	assert.Equal(t, int64(2), claims.UserID)
	assert.True(t, claims.IsAdmin)
}

func TestJWTService_ValidateToken_Expired(t *testing.T) {
	svc := NewJWTService("test-secret", -1*time.Hour, -1*time.Hour)

	token, err := svc.IssueToken(TokenTypeWallet, 1, false)
	require.NoError(t, err)

	_, err = svc.ValidateToken(token)
	assert.Error(t, err)
}

func TestJWTService_ValidateToken_InvalidSignature(t *testing.T) {
	svc1 := NewJWTService("secret-1", 4*time.Hour, 8*time.Hour)
	svc2 := NewJWTService("secret-2", 4*time.Hour, 8*time.Hour)

	token, err := svc1.IssueToken(TokenTypeWallet, 1, false)
	require.NoError(t, err)

	_, err = svc2.ValidateToken(token)
	assert.Error(t, err)
}

func TestJWTService_ValidateToken_InvalidFormat(t *testing.T) {
	svc := NewJWTService("test-secret", 4*time.Hour, 8*time.Hour)

	_, err := svc.ValidateToken("not-a-jwt")
	assert.Error(t, err)
}
