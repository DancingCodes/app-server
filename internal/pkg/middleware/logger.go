package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 记录请求的基本信息：方法、路径、状态码、耗时
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 请求前：记录开始时间
		startTime := time.Now()

		// 2. 执行后续逻辑（也就是你的 Handler）
		c.Next()

		// 3. 请求后：计算耗时并打印日志
		endTime := time.Now()
		latencyTime := endTime.Sub(startTime) // 耗时
		reqMethod := c.Request.Method         // GET/POST
		reqUri := c.Request.RequestURI        // 路径
		statusCode := c.Writer.Status()       // HTTP 状态码
		clientIP := c.ClientIP()              // 客户端 IP

		log.Printf("| %3d | %13v | %15s | %s | %s |",
			statusCode,
			latencyTime,
			clientIP,
			reqMethod,
			reqUri,
		)
	}
}
