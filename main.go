package main // 声明这是主程序入口包

import (
	"app-server/internal/user" // 你自己写的业务包：包含用户相关的逻辑
	"log"                      // Go标准库：用于记录日志和处理致命错误

	"github.com/gin-gonic/gin" // 第三方库：用于处理 HTTP 请求和路由
	"gorm.io/driver/mysql"     // 第三方库：GORM 的 MySQL 驱动，负责跟 MySQL 通信
	"gorm.io/gorm"             // 第三方库：GORM 数据库操作核心框架
)

func main() {
	// 1. 定义 DSN (Data Source Name)：
	// 这里的格式是 用户名:密码@tcp(服务器IP:端口)/数据库名?配置参数
	dsn := "app_db:DancingCodes1227@tcp(82.156.9.114:3306)/app_db?charset=utf8mb4&parseTime=True&loc=Local"

	// 2. 使用 GORM 打开数据库连接：
	// mysql.Open(dsn) 告之使用 MySQL 驱动
	// &gorm.Config{} 是数据库的高级配置，现在传空代表使用默认配置
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// 检查连接是否失败，如果 err 不为空，打印错误并强制停止程序
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 3. 自动迁移 (AutoMigrate)：
	// 这是 GORM 的核心功能，它会检查数据库。如果表不存在，它会按照 user.User 结构体自动建表。
	// 如果你修改了 User 结构体增加了字段，它也会自动帮你修改表结构。
	err = db.AutoMigrate(&user.User{})
	if err != nil {
		return // 如果建表失败，直接退出
	}

	// 4. 依赖注入 (Dependency Injection)：

	// (A) 把数据库连接 (db) 交给“仓库层”负责存取
	repo := user.NewRepository(db)

	// (B) 把仓库层 (repo) 交给“服务层”负责处理业务逻辑
	svc := user.NewService(repo)

	// (C) 把服务层 (svc) 交给“接口层”负责接收网络请求
	handler := user.NewHandler(svc)

	// 5. 初始化 Gin 引擎：
	// Default() 会自带 Logger (打印日志) 和 Recovery (崩溃恢复) 中间件
	r := gin.Default()

	// 6. 定义路由规则：
	// 当有人用 POST 方法访问 /register 时，交给 handler 里的 Register 函数处理
	r.POST("/register", handler.Register)
	// 当有人用 POST 方法访问 /login 时，交给 handler 里的 Login 函数处理
	r.POST("/login", handler.Login)

	// 7. 启动 HTTP 服务器：
	// 监听 8080 端口。这是一个阻塞操作，程序会停在这里一直等待外部请求。
	err = r.Run(":8080")
	if err != nil {
		return // 如果端口被占用导致启动失败，直接退出
	}
}
