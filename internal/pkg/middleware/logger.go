package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// 用于解析响应中的业务代码
type bizResponse struct {
	Code int `json:"code"`
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseBodyWriter) Write(b []byte) (int, error) {
	if b != nil {
		w.body.Write(b)
	}
	return w.ResponseWriter.Write(b)
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ANSI 颜色转义字符
		const (
			blue   = "\033[97;44m" // 蓝色背景 (200 成功)
			yellow = "\033[90;43m" // 黄色背景 (401 警告)
			red    = "\033[97;41m" // 红色背景 (500 错误)
			reset  = "\033[0m"     // 重置颜色
		)

		startTime := time.Now()

		// 1. 备份请求参数
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 2. 包装 ResponseWriter 以捕获响应内容
		writer := responseBodyWriter{
			body:           bytes.NewBuffer(nil),
			ResponseWriter: c.Writer,
		}
		c.Writer = writer

		c.Next()

		// 3. 计算耗时并转换为毫秒 (ms)
		// 使用 Duration 的 Milliseconds 方法或浮点数计算以保留小数
		latency := float64(time.Since(startTime).Microseconds()) / 1000.0

		// 4. 解析业务 Code 决定颜色
		var statusColor = reset
		respBytes := writer.body.Bytes()
		var biz bizResponse

		if len(respBytes) > 0 {
			if err := json.Unmarshal(respBytes, &biz); err == nil {
				switch biz.Code {
				case 200:
					statusColor = blue
				case 401:
					statusColor = yellow
				case 500:
					statusColor = red
				default:
					statusColor = reset
				}
			}
		}

		// 5. 格式化打印中文日志
		// Latency 改为毫秒显示，标签全部中文化
		log.Printf("\n%s [日志] %s %s %s | 来源IP: %s | 耗时: %.2f毫秒\n"+
			"  ├─ [请求参数]: %s\n"+
			"  └─ [返回结果]: %s\n"+
			"----------------------------------------------------------------",
			statusColor, c.Request.Method, c.Request.RequestURI, reset,
			c.ClientIP(),
			latency,
			string(requestBody),
			string(respBytes),
		)
	}
}
