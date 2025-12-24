package user

import (
	"context"

	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, req *RegisterRequest) error
	Login(ctx context.Context, account string, password string) (*User, error)
}

type ServiceImpl struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &ServiceImpl{repo: r}
}

func (s *ServiceImpl) Register(ctx context.Context, req *RegisterRequest) error {
	// 1. 检查用户是否已存在（这是一个好的习惯）
	existUser, _ := s.repo.GetByAccount(ctx, req.Account)
	if existUser != nil && existUser.ID != 0 {
		return errors.New("用户已存在")
	}

	// 2. 将明文密码加密
	// GenerateFromPassword 的第二个参数是 Cost，通常用 DefaultCost (10)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 3. 将加密后的字符串存入数据库
	user := &User{
		Account:  req.Account,
		Password: string(hashedPassword),
	}

	return s.repo.Create(ctx, user)
}

func (s *ServiceImpl) Login(ctx context.Context, account, password string) (*User, error) {
	// 1. 先根据账号找到这个用户
	user, err := s.repo.GetByAccount(ctx, account)
	if err != nil {
		return nil, errors.New("账号或密码错误") // 统一报错，不告诉前端是账号没找到还是密码错了，更安全
	}

	// 2. 使用 Bcrypt 工具比对【数据库里的密文】和【用户输入的明文】
	// 成功返回 nil，失败返回 error
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("账号或密码错误")
	}

	return user, nil
}
