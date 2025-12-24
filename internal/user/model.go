package user

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Account   string    `gorm:"unique;not null;comment:账号" json:"account"`
	Password  string    `gorm:"not null;comment:密码" json:"-"`
	Nickname  string    `gorm:"size:50;comment:昵称" json:"nickname"`
	Avatar    string    `gorm:"comment:头像地址" json:"avatar"`
	Signature string    `gorm:"size:255;comment:个人签名" json:"signature"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}

// RegisterRequest 专门用于注册接口，接收前端传来的 JSON
type RegisterRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname"` // 昵称可以选填
}

// LoginRequest 专门用于登录接口
type LoginRequest struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}
