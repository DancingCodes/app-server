package user

import (
	"context"
)

type Repository interface {
	// Create 插入一条新用户记录
	Create(ctx context.Context, user *User) error

	// GetByAccount 根据账号查找用户
	GetByAccount(ctx context.Context, account string) (*User, error)

	// GetByID 根据 ID 查找用户
	GetByID(ctx context.Context, id uint) (*User, error)

	// Update 更新用户信息
	Update(ctx context.Context, user *User) error
}
