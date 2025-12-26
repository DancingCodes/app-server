package file

import (
	"app-server/internal/pkg/app"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		app.Error(c, http.StatusBadRequest, app.CodeServerErr, "未接收到文件")
		return
	}

	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join("./uploads", fileName)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		app.Error(c, http.StatusInternalServerError, app.CodeServerErr, "文件保存失败")
		return
	}

	// 返回图片的可访问路径
	url := "/uploads/" + fileName
	app.Success(c, gin.H{"url": url})
}
