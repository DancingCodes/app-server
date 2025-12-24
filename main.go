package main

import (
	"app-server/internal/pkg/config"
	"app-server/internal/pkg/database"
	"app-server/internal/pkg/middleware"
	"app-server/internal/user"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 加载配置
	config.LoadConfig("config.yaml")

	// 2. 初始化数据库
	db := database.InitDB(config.GlobalConfig.Database.DSN)

	// 3. 初始化 Gin 引擎 (使用 New 而不是 Default)
	r := gin.New()

	// 4. 【核心步骤】注册中间件
	// 注册你刚写的自定义日志中间件
	r.Use(middleware.Logger())
	// 必须加上 Recovery 中间件，防止程序因为某个 Handler 报错而直接挂掉（闪退）
	r.Use(gin.Recovery())

	// 5. 挂载业务路由
	user.RegisterHandlers(r, db)

	// 6. 启动服务
	err := r.Run(config.GlobalConfig.Server.Port)
	if err != nil {
		return
	}
}
