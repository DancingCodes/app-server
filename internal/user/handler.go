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
		app.Error(c, http.StatusInternalServerError, app.CodeServerErr, "参数格式不正确")
		return
	}

	// 调用 Service
	user, err := h.svc.Login(c.Request.Context(), req.Account, req.Password)

	// 这里的 user == nil 判断非常重要
	// 即使 Service 没写好，Handler 这里也能兜底防止 app.GenerateToken(user.ID) 崩溃
	if err != nil || user == nil {
		app.Error(c, http.StatusInternalServerError, app.CodeServerErr, "账号或密码错误")
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
	// 1. 获取 userID
	uid, exists := c.Get("userID")
	if !exists {
		app.Error(c, http.StatusInternalServerError, app.CodeAuthErr, "身份验证失败")
		return
	}

	// 2. 类型断言
	userID := uid.(uint)

	// 3. 获取用户信息
	user, err := h.svc.GetByID(c.Request.Context(), userID)

	// 🔥 关键修改：在这里同时判断 err 和 user 是否为空
	// 如果数据库没查到，即使 err 是 nil，但 user 为空，我们也认为验证/查询失败
	if err != nil || user == nil {
		// 如果你想返回 401（未授权/身份失效）
		app.Error(c, http.StatusInternalServerError, app.CodeAuthErr, "用户不存在或登录已失效")
		return
	}

	// 4. 返回完整的用户信息
	app.Success(c, gin.H{
		"user": user,
	})
}

// UpdateProfile 修改当前登录用户的资料
func (h *Handler) UpdateProfile(c *gin.Context) {
	// 1. 从 Context 获取 userID (Auth 中间件存入的)
	uid, exists := c.Get("userID")
	if !exists {
		app.Error(c, http.StatusUnauthorized, app.CodeAuthErr, "身份验证失败")
		return
	}
	userID := uid.(uint)

	// 2. 绑定请求参数
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		app.Error(c, http.StatusBadRequest, app.CodeServerErr, "参数格式不正确")
		return
	}

	// 3. 调用 Service 执行更新
	if err := h.svc.UpdateProfile(c.Request.Context(), userID, &req); err != nil {
		app.Error(c, http.StatusInternalServerError, app.CodeServerErr, err.Error())
		return
	}

	// 4. 返回成功
	app.Success(c, "个人资料更新成功")
}
