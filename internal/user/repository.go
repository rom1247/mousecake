// Package user 实现用户认证模块的业务逻辑和数据访问。
package user

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/mousecake-go/mousecake-go/internal/user/domain"
)

// userPO 是 users 表的持久化对象。
type userPO struct {
	ID           int64      `gorm:"column:id;primaryKey"`
	Address      *string    `gorm:"column:address"`
	Username     *string    `gorm:"column:username"`
	PasswordHash *string    `gorm:"column:password_hash"`
	Name         *string    `gorm:"column:name"`
	Nickname     *string    `gorm:"column:nickname"`
	IsAdmin      bool       `gorm:"column:is_admin"`
	Status       string     `gorm:"column:status"`
	LastLoginAt  *time.Time `gorm:"column:last_login_at"`
	CreatedAt    time.Time  `gorm:"column:created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at"`
	DeletedAt    *time.Time `gorm:"column:deleted_at"`
}

func (userPO) TableName() string { return "users" }

// noncePO 是 login_nonces 表的持久化对象。
type noncePO struct {
	ID        int64     `gorm:"column:id;primaryKey"`
	Address   string    `gorm:"column:address"`
	Nonce     string    `gorm:"column:nonce"`
	ExpiresAt time.Time `gorm:"column:expires_at"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (noncePO) TableName() string { return "login_nonces" }

// UserRepository 是 domain.UserRepository 的 Gorm 实现。
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建仓库实例。
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByAddress 根据钱包地址查询用户（排除软删除）。
func (r *UserRepository) FindByAddress(ctx context.Context, address string) (*domain.User, error) {
	var po userPO
	err := r.db.WithContext(ctx).
		Select("id, address, username, password_hash, name, nickname, is_admin, status, last_login_at, created_at, updated_at, deleted_at").
		Where("address = ? AND deleted_at IS NULL", address).
		First(&po).Error
	if err != nil {
		return nil, fmt.Errorf("find user by address %s: %w", address, err)
	}
	return poToEntity(&po), nil
}

// FindByUsername 根据用户名查询用户（排除软删除）。
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	var po userPO
	err := r.db.WithContext(ctx).
		Select("id, address, username, password_hash, name, nickname, is_admin, status, last_login_at, created_at, updated_at, deleted_at").
		Where("username = ? AND deleted_at IS NULL", username).
		First(&po).Error
	if err != nil {
		return nil, fmt.Errorf("find user by username %s: %w", username, err)
	}
	return poToEntity(&po), nil
}

// FindByID 根据 ID 查询用户（排除软删除）。
func (r *UserRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	var po userPO
	err := r.db.WithContext(ctx).
		Select("id, address, username, password_hash, name, nickname, is_admin, status, last_login_at, created_at, updated_at, deleted_at").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&po).Error
	if err != nil {
		return nil, fmt.Errorf("find user by id %d: %w", id, err)
	}
	return poToEntity(&po), nil
}

// Create 创建新用户。
func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	po := entityToPO(user)
	err := r.db.WithContext(ctx).Create(po).Error
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	user.ID = po.ID
	return nil
}

// Update 更新已有用户。
func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	po := entityToPO(user)
	err := r.db.WithContext(ctx).
		Select("address", "username", "password_hash", "name", "nickname", "is_admin", "status", "last_login_at", "updated_at", "deleted_at").
		Where("id = ?", user.ID).
		Updates(po).Error
	if err != nil {
		return fmt.Errorf("update user %d: %w", user.ID, err)
	}
	return nil
}

// UpsertByAddress 根据地址做 UPSERT（存在则更新 LastLoginAt，否则创建）。
func (r *UserRepository) UpsertByAddress(ctx context.Context, user *domain.User) error {
	po := entityToPO(user)
	result := r.db.WithContext(ctx).
		Where("address = ? AND deleted_at IS NULL", user.Address).
		FirstOrCreate(po)
	if result.Error != nil {
		return fmt.Errorf("upsert user by address %s: %w", user.Address, result.Error)
	}

	if po.ID != user.ID {
		user.ID = po.ID
	}

	now := time.Now()
	user.LastLoginAt = &now
	po.LastLoginAt = &now

	err := r.db.WithContext(ctx).
		Select("last_login_at", "updated_at").
		Where("id = ?", po.ID).
		Updates(po).Error
	if err != nil {
		return fmt.Errorf("update last login at for user %d: %w", po.ID, err)
	}
	return nil
}

// SaveNonce 保存或更新 nonce（按地址 UPSERT）。
func (r *UserRepository) SaveNonce(ctx context.Context, nonce *domain.LoginNonce) error {
	po := nonceToPO(nonce)
	err := r.db.WithContext(ctx).
		Where("address = ?", nonce.Address).
		Assign(po).
		FirstOrCreate(po).Error
	if err != nil {
		return fmt.Errorf("save nonce for address %s: %w", nonce.Address, err)
	}
	nonce.CreatedAt = po.CreatedAt
	return nil
}

// FindNonceByAddress 根据地址查询 nonce。
func (r *UserRepository) FindNonceByAddress(ctx context.Context, address string) (*domain.LoginNonce, error) {
	var po noncePO
	err := r.db.WithContext(ctx).
		Select("id, address, nonce, expires_at, created_at").
		Where("address = ?", address).
		First(&po).Error
	if err != nil {
		return nil, fmt.Errorf("find nonce by address %s: %w", address, err)
	}
	return poToNonce(&po), nil
}

// DeleteNonce 根据地址删除 nonce。
func (r *UserRepository) DeleteNonce(ctx context.Context, address string) error {
	err := r.db.WithContext(ctx).
		Where("address = ?", address).
		Delete(&noncePO{}).Error
	if err != nil {
		return fmt.Errorf("delete nonce for address %s: %w", address, err)
	}
	return nil
}

// CleanExpiredNonces 清理所有过期的 nonce 记录。
func (r *UserRepository) CleanExpiredNonces(ctx context.Context) error {
	err := r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&noncePO{}).Error
	if err != nil {
		return fmt.Errorf("clean expired nonces: %w", err)
	}
	return nil
}

func entityToPO(u *domain.User) *userPO {
	po := &userPO{
		ID:          u.ID,
		IsAdmin:     u.IsAdmin,
		Status:      string(u.Status),
		LastLoginAt: u.LastLoginAt,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		DeletedAt:   u.DeletedAt,
	}
	if u.Address != "" {
		po.Address = &u.Address
	}
	if u.Username != "" {
		po.Username = &u.Username
	}
	if u.PasswordHash != "" {
		po.PasswordHash = &u.PasswordHash
	}
	if u.Name != "" {
		po.Name = &u.Name
	}
	if u.Nickname != "" {
		po.Nickname = &u.Nickname
	}
	return po
}

func poToEntity(po *userPO) *domain.User {
	u := &domain.User{
		ID:          po.ID,
		IsAdmin:     po.IsAdmin,
		Status:      domain.UserStatus(po.Status),
		LastLoginAt: po.LastLoginAt,
		CreatedAt:   po.CreatedAt,
		UpdatedAt:   po.UpdatedAt,
		DeletedAt:   po.DeletedAt,
	}
	if po.Address != nil {
		u.Address = *po.Address
	}
	if po.Username != nil {
		u.Username = *po.Username
	}
	if po.PasswordHash != nil {
		u.PasswordHash = *po.PasswordHash
	}
	if po.Name != nil {
		u.Name = *po.Name
	}
	if po.Nickname != nil {
		u.Nickname = *po.Nickname
	}
	return u
}

func nonceToPO(n *domain.LoginNonce) *noncePO {
	return &noncePO{
		Address:   n.Address,
		Nonce:     n.Nonce,
		ExpiresAt: n.ExpiresAt,
		CreatedAt: n.CreatedAt,
	}
}

func poToNonce(po *noncePO) *domain.LoginNonce {
	return &domain.LoginNonce{
		Address:   po.Address,
		Nonce:     po.Nonce,
		ExpiresAt: po.ExpiresAt,
		CreatedAt: po.CreatedAt,
	}
}
