package main

import (
	"app-server/internal/user"

	"github.com/gin-gonic/gin"
)

func main() {
	userRepo := user.NewUserRepository()
	userSvc := user.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userSvc)

	// 4. 注册路由
	r := gin.Default()
	r.POST("/register", userHandler.Register)
	r.Run(":8080")
}
