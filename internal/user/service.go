package user

type Service interface {
	Register(username, password string) error
}

type userService struct {
	repo Repository // 依赖接口而非具体实现
}

func NewUserService(r Repository) Service {
	return &userService{repo: r}
}

func (s *userService) Register(username, password string) error {
	// 这里写复杂的业务：比如密码加密、敏感词过滤等
	u := &User{Username: username, Password: password}
	return s.repo.Create(u)
}
