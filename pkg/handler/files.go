package handler

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
)

func (h *BaseHandler) HandleFileUpload(ctx context.Context, c *app.RequestContext) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, ErrorResponse(400, fmt.Sprintf("获取上传文件失败: %v", err)))
		return
	}

	targetPath := string(c.FormValue("path"))
	if targetPath == "" {
		targetPath = file.Filename
	}

	targetPath = filepath.Clean(targetPath)
	if filepath.IsAbs(targetPath) || targetPath == ".." || strings.HasPrefix(targetPath, "..") {
		c.JSON(400, ErrorResponse(400, "无效的文件路径"))
		return
	}

	savePath := filepath.Join(h.StaticDataPath(), targetPath)

	saveDir := filepath.Dir(savePath)
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		c.JSON(500, ErrorResponse(500, fmt.Sprintf("创建目录失败: %v", err)))
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(500, ErrorResponse(500, fmt.Sprintf("打开上传文件失败: %v", err)))
		return
	}
	defer src.Close()

	dst, err := os.Create(savePath)
	if err != nil {
		c.JSON(500, ErrorResponse(500, fmt.Sprintf("创建目标文件失败: %v", err)))
		return
	}
	defer dst.Close()

	buf := make([]byte, 32*1024)
	written, err := io.CopyBuffer(dst, src, buf)
	if err != nil {
		os.Remove(savePath)
		c.JSON(500, ErrorResponse(500, fmt.Sprintf("保存文件失败: %v", err)))
		return
	}

	// 返回成功响应
	c.JSON(200, SuccessResponse(map[string]interface{}{
		"filename": file.Filename,
		"path":     targetPath,
		"size":     written,
	}))
}
