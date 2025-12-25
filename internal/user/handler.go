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
		app.Error(c, http.StatusBadRequest, app.CodeServerErr, "参数格式不正确")
		return
	}

	// 1. 调用 Service 注册
	user, err := h.svc.Register(c.Request.Context(), &req)
	if err != nil {
		app.Error(c, http.StatusInternalServerError, app.CodeServerErr, err.Error())
		return
	}

	// 2. 注册成功生成 Token
	token, err := app.GenerateToken(user.ID)
	if err != nil {
		app.Error(c, http.StatusInternalServerError, app.CodeServerErr, "生成令牌失败")
		return
	}

	// 3. 移除 user 对象，仅返回 token
	app.Success(c, gin.H{
		"token": token,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.Error(c, http.StatusBadRequest, app.CodeServerErr, "参数格式不正确")
		return
	}

	// 调用 Service
	user, err := h.svc.Login(c.Request.Context(), req.Account, req.Password)

	// 这里的 user == nil 判断非常重要
	// 即使 Service 没写好，Handler 这里也能兜底防止 app.GenerateToken(user.ID) 崩溃
	if err != nil || user == nil {
		app.Error(c, http.StatusUnauthorized, app.CodeServerErr, "账号或密码错误")
		return
	}

	token, err := app.GenerateToken(user.ID)
	if err != nil {
		app.Error(c, http.StatusInternalServerError, app.CodeServerErr, "生成令牌失败")
		return
	}

	app.Success(c, gin.H{
		"token": token,
	})
}

// GetProfile 获取当前登录用户的个人资料
func (h *Handler) GetProfile(c *gin.Context) {
	// 1. 从 Context 中取出中间件已经设置好的 userID
	// 这个值是 middleware.Auth 校验 Token 成功后存进去的
	uid, exists := c.Get("userID")
	if !exists {
		app.Error(c, http.StatusUnauthorized, app.CodeAuthErr, "身份验证失败")
		return
	}

	// 2. 类型断言
	userID := uid.(uint)

	// 3. 直接用这个 ID 去数据库查，不需要前端传任何 ID 参数
	user, err := h.svc.GetByID(c.Request.Context(), userID)
	if err != nil {
		app.Error(c, http.StatusNotFound, app.CodeServerErr, "用户不存在")
		return
	}

	// 4. 返回完整的用户信息
	app.Success(c, gin.H{
		"user": user,
	})
}
