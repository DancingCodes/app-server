package middleware

import (
	"app-server/internal/pkg/app"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从 Header 拿 Token (标准格式：Authorization: Bearer <token>)
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			app.Error(c, http.StatusUnauthorized, app.CodeAuthErr, "未登录")
			c.Abort() // 拦截，不许往后走
			return
		}

		// 2. 截取 Token 字符串
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			app.Error(c, http.StatusUnauthorized, app.CodeAuthErr, "Token格式错误")
			c.Abort()
			return
		}

		// 3. 解析并校验 Token
		claims, err := app.ParseToken(parts[1])
		if err != nil {
			app.Error(c, http.StatusUnauthorized, app.CodeAuthErr, "无效的令牌")
			c.Abort()
			return
		}

		// 4. 将解析出来的用户信息塞进上下文，方便后续接口直接拿
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
