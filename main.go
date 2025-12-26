package main

import (
	"app-server/internal/file"
	"app-server/internal/pkg/config"
	"app-server/internal/pkg/database"
	"app-server/internal/pkg/middleware"
	"app-server/internal/user"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig("config.yaml")

	db := database.InitDB(config.GlobalConfig.Database.DSN)

	r := gin.New()

	r.Use(middleware.Logger())
	r.Use(gin.Recovery())

	r.Static("/uploads", "./uploads")

	user.Router(r, db)
	file.Router(r)

	err := r.Run(config.GlobalConfig.Server.Port)
	if err != nil {
		return
	}
}
