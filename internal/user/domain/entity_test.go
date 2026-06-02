package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// 4.1 User 实体创建
func TestNewWalletUser(t *testing.T) {
	user, err := NewWalletUser("0x1234567890abcdef1234567890abcdef12345678")
	require.NoError(t, err)

	assert.Equal(t, "0x1234567890abcdef1234567890abcdef12345678", user.Address)
	assert.Empty(t, user.Username)
	assert.Empty(t, user.PasswordHash)
	assert.False(t, user.IsAdmin)
	assert.Equal(t, UserStatusActive, user.Status)
	assert.WithinDuration(t, time.Now(), user.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), user.UpdatedAt, time.Second)
}

func TestNewAdminUser(t *testing.T) {
	user, err := NewAdminUser("admin", "hashed_password")
	require.NoError(t, err)

	assert.Equal(t, "admin", user.Username)
	assert.Equal(t, "hashed_password", user.PasswordHash)
	assert.Empty(t, user.Address)
	assert.True(t, user.IsAdmin)
	assert.Equal(t, UserStatusActive, user.Status)
}

// 4.3 UserStatus 状态机
func TestUser_Disable(t *testing.T) {
	user, _ := NewWalletUser("0x1234567890abcdef1234567890abcdef12345678")
	require.Equal(t, UserStatusActive, user.Status)

	err := user.Disable()
	require.NoError(t, err)
	assert.Equal(t, UserStatusDisabled, user.Status)
}

func TestUser_Enable(t *testing.T) {
	user, _ := NewWalletUser("0x1234567890abcdef1234567890abcdef12345678")
	user.Disable()

	err := user.Enable()
	require.NoError(t, err)
	assert.Equal(t, UserStatusActive, user.Status)
}

func TestUser_DisableAlreadyDisabled(t *testing.T) {
	user, _ := NewWalletUser("0x1234567890abcdef1234567890abcdef12345678")
	user.Disable()

	err := user.Disable()
	assert.Error(t, err)
	assert.Equal(t, UserStatusDisabled, user.Status)
}

// 4.5 钱包地址绑定规则
func TestUser_BindAddress(t *testing.T) {
	user, _ := NewAdminUser("admin", "hash")

	err := user.BindAddress("0x1234567890abcdef1234567890abcdef12345678")
	require.NoError(t, err)
	assert.Equal(t, "0x1234567890abcdef1234567890abcdef12345678", user.Address)
}

func TestUser_BindAddress_InvalidFormat(t *testing.T) {
	user, _ := NewAdminUser("admin", "hash")

	err := user.BindAddress("invalid")
	assert.Error(t, err)
}

func TestUser_BindAddress_AlreadyBound(t *testing.T) {
	user, _ := NewAdminUser("admin", "hash")
	user.BindAddress("0x1234567890abcdef1234567890abcdef12345678")

	err := user.BindAddress("0xabcdef1234567890abcdef1234567890abcdef12")
	assert.Error(t, err)
}

// 4.7 软删除
func TestUser_SoftDelete(t *testing.T) {
	user, _ := NewWalletUser("0x1234567890abcdef1234567890abcdef12345678")
	assert.Nil(t, user.DeletedAt)

	err := user.SoftDelete()
	require.NoError(t, err)
	assert.NotNil(t, user.DeletedAt)
}

func TestUser_SoftDelete_AlreadyDeleted(t *testing.T) {
	user, _ := NewWalletUser("0x1234567890abcdef1234567890abcdef12345678")
	user.SoftDelete()

	err := user.SoftDelete()
	assert.Error(t, err)
}

// 4.9 LoginNonce 值对象
func TestNewLoginNonce(t *testing.T) {
	nonce, err := NewLoginNonce("0xabc...", 5*time.Minute)
	require.NoError(t, err)

	assert.Equal(t, "0xabc...", nonce.Address)
	assert.Len(t, nonce.Nonce, 32)
	assert.WithinDuration(t, time.Now().Add(5*time.Minute), nonce.ExpiresAt, time.Second)
	assert.False(t, nonce.IsExpired())
}

func TestLoginNonce_IsExpired(t *testing.T) {
	nonce := &LoginNonce{
		Address:   "0xabc...",
		Nonce:     "test-nonce",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	assert.True(t, nonce.IsExpired())
}
