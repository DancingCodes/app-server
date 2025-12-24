package user

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterHandlers(r *gin.Engine, db *gorm.DB) {
	repo := NewRepository(db)
	svc := NewService(repo)
	handler := NewHandler(svc)

	r.POST("/register", handler.Register)
	r.POST("/login", handler.Login)
}
