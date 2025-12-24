package user

import (
	"context"

	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, req *RegisterRequest) error
	Login(ctx context.Context, account, password string) (*User, error)
}

type ServiceImpl struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &ServiceImpl{repo: r}
}

func (s *ServiceImpl) Register(ctx context.Context, req *RegisterRequest) error {
	existingUser, _ := s.repo.GetByAccount(ctx, req.Account)
	if existingUser != nil {
		return errors.New("账号已存在")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newUser := &User{
		Account:  req.Account,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
	}

	return s.repo.Create(ctx, newUser)
}

func (s *ServiceImpl) Login(ctx context.Context, account, password string) (*User, error) {
	u, err := s.repo.GetByAccount(ctx, account)
	if err != nil {
		return nil, errors.New("账号不存在")
	}
	if u.Password != password {
		return nil, errors.New("密码错误")
	}
	return u, nil
}
