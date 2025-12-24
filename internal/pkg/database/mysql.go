package database

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 删掉 var DB *gorm.DB

func InitDB(dsn string) *gorm.DB {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	return db // 只返回，不存全局
}
