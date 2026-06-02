// Package domain 定义用户模块的领域模型，包含实体、值对象和仓库接口。
// 本包不能导入任何外层包（编译期隔离）。
package domain

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// UserStatus 表示用户状态。
type UserStatus string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusDisabled UserStatus = "disabled"
)

// ethAddressRegex 以太坊地址格式校验正则。
var ethAddressRegex = regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)

// User 是用户聚合根实体。
type User struct {
	// ID 数据库主键。
	ID int64
	// Address 用户钱包地址，格式为 0x 开头的 42 位十六进制小写字符串。
	Address string
	// Username 管理员用户名，钱包用户为空。
	Username string
	// PasswordHash 管理员密码哈希（bcrypt），钱包用户为空。
	PasswordHash string
	// Name 用户真实姓名。
	Name string
	// Nickname 用户昵称。
	Nickname string
	// IsAdmin 是否为管理员。
	IsAdmin bool
	// Status 用户状态（active 或 disabled）。
	Status UserStatus
	// LastLoginAt 最近一次登录时间。
	LastLoginAt *time.Time
	// CreatedAt 记录创建时间。
	CreatedAt time.Time
	// UpdatedAt 记录更新时间。
	UpdatedAt time.Time
	// DeletedAt 软删除时间，nil 表示未删除。
	DeletedAt *time.Time
}

// NewWalletUser 创建钱包用户实体，地址转为小写。
func NewWalletUser(address string) (*User, error) {
	if !ethAddressRegex.MatchString(address) {
		return nil, fmt.Errorf("invalid address format: %s", address)
	}

	now := time.Now()
	return &User{
		Address:   strings.ToLower(address),
		Status:    UserStatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

// NewAdminUser 创建管理员用户实体。
func NewAdminUser(username, passwordHash string) (*User, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}
	if passwordHash == "" {
		return nil, errors.New("password hash is required")
	}

	now := time.Now()
	return &User{
		Username:     username,
		PasswordHash: passwordHash,
		IsAdmin:      true,
		Status:       UserStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// Disable 将用户状态设为 Disabled。
func (u *User) Disable() error {
	if u.Status == UserStatusDisabled {
		return errors.New("user already disabled")
	}
	u.Status = UserStatusDisabled
	u.UpdatedAt = time.Now()
	return nil
}

// Enable 将用户状态设为 Active。
func (u *User) Enable() error {
	if u.Status == UserStatusActive {
		return errors.New("user already active")
	}
	u.Status = UserStatusActive
	u.UpdatedAt = time.Now()
	return nil
}

// BindAddress 为管理员用户绑定钱包地址。
func (u *User) BindAddress(address string) error {
	if u.Address != "" {
		return errors.New("address already bound")
	}
	if !ethAddressRegex.MatchString(address) {
		return fmt.Errorf("invalid address format: %s", address)
	}

	u.Address = strings.ToLower(address)
	u.UpdatedAt = time.Now()
	return nil
}

// SoftDelete 软删除用户。
func (u *User) SoftDelete() error {
	if u.DeletedAt != nil {
		return errors.New("user already deleted")
	}
	now := time.Now()
	u.DeletedAt = &now
	u.UpdatedAt = now
	return nil
}

// LoginNonce 是 SIWE 登录的 nonce 值对象。
type LoginNonce struct {
	// Address 钱包地址，与 SIWE 消息中的 address 对应。
	Address string
	// Nonce 服务端生成的随机数，用于防重放攻击。
	Nonce string
	// ExpiresAt nonce 过期时间。
	ExpiresAt time.Time
	// CreatedAt 记录创建时间。
	CreatedAt time.Time
}

// NewLoginNonce 创建新的 LoginNonce 值对象。
func NewLoginNonce(address string, ttl time.Duration) (*LoginNonce, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return nil, fmt.Errorf("generate nonce: %w", err)
	}

	now := time.Now()
	return &LoginNonce{
		Address:   address,
		Nonce:     hex.EncodeToString(bytes),
		ExpiresAt: now.Add(ttl),
		CreatedAt: now,
	}, nil
}

// IsExpired 判断 nonce 是否已过期。
func (n *LoginNonce) IsExpired() bool {
	return time.Now().After(n.ExpiresAt)
}
