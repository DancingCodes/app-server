package user

import (
	"context"

	"gorm.io/gorm"
)

type mysqlRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &mysqlRepository{db: db}
}

func (r *mysqlRepository) Create(ctx context.Context, u *User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *mysqlRepository) GetByAccount(ctx context.Context, account string) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).Where("account = ?", account).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *mysqlRepository) GetByID(ctx context.Context, id uint) (*User, error) {
	var u User
	err := r.db.WithContext(ctx).First(&u, id).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *mysqlRepository) Update(ctx context.Context, u *User) error {
	return r.db.WithContext(ctx).Save(u).Error
}
