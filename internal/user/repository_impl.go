package user

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type mysqlRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &mysqlRepository{db: db}
}

func (r *mysqlRepository) Create(ctx context.Context, u *User) error {
	// 使用 WithContext 传递链路上下文，方便超时控制和日志追踪
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *mysqlRepository) GetByAccount(ctx context.Context, account string) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).Where("account = ?", account).First(&u).Error
	if err != nil {
		// 如果没找到用户，First 会报错，我们要处理这个错误
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // 返回 nil 表示用户不存在，而不是抛出数据库错误
		}
		return nil, err
	}
	return &u, nil
}

func (r *mysqlRepository) GetByID(ctx context.Context, id uint) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).First(&u, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *mysqlRepository) Update(ctx context.Context, u *User) error {
	// Save 会更新所有字段，通常用于全量更新
	return r.db.WithContext(ctx).Save(u).Error
}
