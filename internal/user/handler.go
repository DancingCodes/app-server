package user

import (
	"app-server/internal/pkg/app"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(s Service) *Handler {
	return &Handler{svc: s}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// 哪怕是参数不对，我们也统一返回 CodeServerErr (500)
		app.Error(c, http.StatusBadRequest, app.CodeServerErr, "参数格式不正确")
		return
	}

	err := h.svc.Register(c.Request.Context(), &req)
	if err != nil {
		app.Error(c, http.StatusInternalServerError, app.CodeServerErr, err.Error())
		return
	}

	app.Success(c, nil)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.Error(c, http.StatusBadRequest, app.CodeServerErr, "参数格式不正确")
		return
	}

	user, err := h.svc.Login(c.Request.Context(), req.Account, req.Password)
	if err != nil {
		// 这里的业务码也改用常量 CodeServerErr
		app.Error(c, http.StatusUnauthorized, app.CodeServerErr, "账号或密码错误")
		return
	}

	app.Success(c, user)
}
