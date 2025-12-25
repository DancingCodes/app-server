package middleware

import (
	"bytes"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// responseBodyWriter 包装了 gin.ResponseWriter，用于拦截返回内容
type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b) // 将返回内容备份到 buffer
	return w.ResponseWriter.Write(b)
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 颜色代码定义
		const (
			green  = "\033[97;42m"
			yellow = "\033[90;43m"
			red    = "\033[97;41m"
			reset  = "\033[0m"
		)

		startTime := time.Now()

		// 1. 处理请求参数
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 2. 包装 ResponseWriter
		writer := &responseBodyWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = writer

		c.Next()

		// 3. 计算结果
		latencyTime := time.Since(startTime)
		statusCode := c.Writer.Status()

		// --- 动态选择颜色 ---
		var statusColor string
		switch {
		case statusCode >= 200 && statusCode < 300:
			statusColor = green
		case statusCode >= 300 && statusCode < 400:
			statusColor = yellow
		default:
			statusColor = red
		}

		// 4. 打印带颜色的日志
		log.Printf("\n%s [Request] %s %s %s | IP: %s | Status: %d | Latency: %v\n"+
			"[Params] %s\n"+
			"[Response] %s\n"+
			"----------------------------------------------------------------",
			statusColor, c.Request.Method, c.Request.RequestURI, reset,
			c.ClientIP(),
			statusCode,
			latencyTime,
			string(requestBody),
			writer.body.String(),
		)
	}
}
