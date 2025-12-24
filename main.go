package main

import (
	"app-server/internal/config"
	"app-server/internal/database"

	"app-server/internal/user"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. 加载配置（把 config.yaml 里的东西读到内存里）
	config.LoadConfig("config.yaml")

	// 2. 从配置里拿出 DSN 传给数据库初始化函数
	db := database.InitDB(config.GlobalConfig.Database.DSN)

	// 3. 剩下的逻辑保持不变
	r := gin.Default()
	user.RegisterHandlers(r, db)

	// 端口也可以从配置读
	err := r.Run(config.GlobalConfig.Server.Port)
	if err != nil {
		return
	}
}
