package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 标准 JSON 返回结构
type Response struct {
	Code int         `json:"code"` // 业务自定义状态码（不是 HTTP 状态码）
	Msg  string      `json:"msg"`  // 提示信息
	Data interface{} `json:"data"` // 具体的数据内容，没有则返回 null
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess, // 直接用刚才定义的常量 200
		Msg:  "success",
		Data: data,
	})
}

// Error 函数保持不变，但在调用时传入常量
func Error(c *gin.Context, httpCode int, businessCode int, msg string) {
	c.JSON(httpCode, Response{
		Code: businessCode,
		Msg:  msg,
		Data: nil,
	})
}
