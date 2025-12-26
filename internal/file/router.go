package file

import (
	"app-server/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	h := NewHandler()

	g := r.Group("/common", middleware.Auth())
	{
		g.POST("/upload", h.Upload)
	}
}
