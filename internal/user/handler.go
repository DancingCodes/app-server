package user

import (
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc Service
}

func NewUserHandler(s Service) *UserHandler {
	return &UserHandler{svc: s}
}

func (h *UserHandler) Register(c *gin.Context) {
	c.JSON(200, gin.H{"message": "success"})
}
