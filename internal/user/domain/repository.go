package domain

import "context"

// UserRepository 定义用户聚合根的持久化接口。
type UserRepository interface {
	// FindByAddress 根据钱包地址查询用户（排除软删除）。
	FindByAddress(ctx context.Context, address string) (*User, error)
	// FindByUsername 根据用户名查询用户（排除软删除）。
	FindByUsername(ctx context.Context, username string) (*User, error)
	// FindByID 根据 ID 查询用户（排除软删除）。
	FindByID(ctx context.Context, id int64) (*User, error)
	// Create 创建新用户。
	Create(ctx context.Context, user *User) error
	// Update 更新已有用户。
	Update(ctx context.Context, user *User) error
	// UpsertByAddress 根据地址做 UPSERT（存在则更新 LastLoginAt，否则创建）。
	UpsertByAddress(ctx context.Context, user *User) error
	// SaveNonce 保存或更新 nonce（按地址 UPSERT）。
	SaveNonce(ctx context.Context, nonce *LoginNonce) error
	// FindNonceByAddress 根据地址查询 nonce。
	FindNonceByAddress(ctx context.Context, address string) (*LoginNonce, error)
	// DeleteNonce 根据 地址删除 nonce。
	DeleteNonce(ctx context.Context, address string) error
	// CleanExpiredNonces 清理所有过期的 nonce 记录。
	CleanExpiredNonces(ctx context.Context) error
}
