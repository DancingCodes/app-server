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
	// 1. 声明 req 变量（这一步不能漏！）
	var req LoginRequest

	// 2. 绑定参数：将请求体中的 JSON 映射到 req 结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		app.Error(c, http.StatusBadRequest, app.CodeServerErr, "参数格式不正确")
		return
	}

	// 3. 现在 req.Account 和 req.Password 就可以安全使用了
	user, err := h.svc.Login(c.Request.Context(), req.Account, req.Password)
	if err != nil {
		app.Error(c, http.StatusUnauthorized, app.CodeAuthErr, "账号或密码错误")
		return
	}

	// 4. 生成 Token 并返回（如之前所述）
	token, err := app.GenerateToken(user.ID)
	if err != nil {
		app.Error(c, http.StatusInternalServerError, app.CodeServerErr, "生成令牌失败")
		return
	}

	app.Success(c, gin.H{
		"token": token,
		"user":  user,
	})
}

// GetProfile 获取当前登录用户的个人资料
func (h *Handler) GetProfile(c *gin.Context) {
	// 1. 从 Context 中取出中间件设置的 userID
	// 注意：middleware.Auth 里的 c.Set("userID", claims.UserID) 存的是什么类型，这里就要转成什么类型
	uid, exists := c.Get("userID")
	if !exists {
		app.Error(c, http.StatusUnauthorized, app.CodeAuthErr, "未登录或登录已过期")
		return
	}

	// 2. 类型断言 (将 interface{} 转回 uint)
	userID, ok := uid.(uint)
	if !ok {
		app.Error(c, http.StatusInternalServerError, app.CodeServerErr, "内部系统错误")
		return
	}

	// 3. 模拟获取用户信息 (实际开发中你会调用 h.svc.GetUserByID)
	// 这里我们先直接返回 userID 证明 JWT 校验成功
	app.Success(c, gin.H{
		"user_id": userID,
		"remark":  "恭喜，你能看到这个信息说明你的 JWT 校验成功了！",
	})
}
