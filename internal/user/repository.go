package user

// Repository 定义了用户模块对数据库的操作规范
type Repository interface {
	Create(u *User) error
	GetByID(id int64) (*User, error)
}

// 具体的实现类（比如 MySQL 实现）
type mysqlRepo struct {
	// db *sql.DB 或者 gorm.DB
}

func NewUserRepository() Repository {
	return &mysqlRepo{}
}

func (m *mysqlRepo) Create(u *User) error {
	// 这里写具体的 SQL 操作
	return nil
}

func (m *mysqlRepo) GetByID(id int64) (*User, error) {
	return &User{ID: id, Username: "test"}, nil
}
