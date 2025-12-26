package user

import (
	"app-server/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Router(r *gin.Engine, db *gorm.DB) {
	// 初始化依赖
	repo := NewRepository(db)
	svc := NewService(repo)
	h := NewHandler(svc)

	// --- 公开路由：不需要登录 ---
	r.POST("/register", h.Register)
	r.POST("/login", h.Login)

	// --- 受保护路由：需要带上 Token 才能访问 ---
	// 只有经过 middleware.Auth() 检查通过的请求，才会到达里面的 Handler
	authGroup := r.Group("/u")
	authGroup.Use(middleware.Auth())
	{
		// 这里的路径会自动变为 /u/profile
		authGroup.GET("/profile", h.GetProfile)
	}
}
