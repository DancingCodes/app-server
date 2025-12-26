package user

import (
	"context"

	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Register(ctx context.Context, req *RegisterRequest) (*User, error)
	Login(ctx context.Context, account string, password string) (*User, error)
	GetByID(ctx context.Context, id uint) (*User, error)
	UpdateProfile(ctx context.Context, userID uint, req *UpdateProfileRequest) error
}

type ServiceImpl struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &ServiceImpl{repo: r}
}

func (s *ServiceImpl) Register(ctx context.Context, req *RegisterRequest) (*User, error) {
	// 1. 检查用户是否已存在
	existUser, _ := s.repo.GetByAccount(ctx, req.Account)
	if existUser != nil && existUser.ID != 0 {
		return nil, errors.New("用户已存在")
	}

	// 2. 将明文密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	defaultAvatar := "https://filestore.moonc.love/uploadFiles/1766716767733-395871637.png"

	// 3. 构造用户对象
	user := &User{
		Nickname: req.Nickname,
		Account:  req.Account,
		Password: string(hashedPassword),
		Avatar:   defaultAvatar,
	}

	// 4. 将用户存入数据库
	// 此时 s.repo.Create 会填充 user 的 ID
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	// 5. 返回 user 指针和 nil
	return user, nil
}

// ... Login 实现保持不变

func (s *ServiceImpl) Login(ctx context.Context, account, password string) (*User, error) {
	// 1. 查询用户
	user, err := s.repo.GetByAccount(ctx, account)

	// --- 核心修复点 ---
	// 手动删除数据后，GetByAccount 可能会返回 (nil, nil) 或 (nil, err)
	// 我们必须在这里拦截，不能让程序往下走到 bcrypt 那一行
	if err != nil || user == nil {
		return nil, errors.New("账号不存在")
	}

	// 2. 只有 user != nil，访问 user.Password 才是安全的
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("密码错误")
	}

	return user, nil
}

func (s *ServiceImpl) GetByID(ctx context.Context, id uint) (*User, error) {
	// 调用 repository 层的查询方法
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *ServiceImpl) UpdateProfile(ctx context.Context, userID uint, req *UpdateProfileRequest) error {
	// 1. 先从数据库查出老数据
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("用户不存在")
	}

	// 2. 处理账号修改 (Account) - 需要查重
	if req.Account != "" && req.Account != user.Account {
		existUser, _ := s.repo.GetByAccount(ctx, req.Account)
		if existUser != nil && existUser.ID != 0 {
			return errors.New("该账号已被他人占用")
		}
		user.Account = req.Account
	}

	// 3. 处理密码修改 (Password) - 需要重新加密
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}

	// 4. 处理其他基础字段 (只有前端传了值才覆盖)
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Signature != "" {
		user.Signature = req.Signature
	}

	// 5. 调用 repository 的 Update 方法保存回数据库
	return s.repo.Update(ctx, user)
}
