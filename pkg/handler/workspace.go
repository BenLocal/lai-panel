package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

type workspaceEntry struct {
	Name    string    `json:"name"`
	Path    string    `json:"path"`
	IsDir   bool      `json:"is_dir"`
	Size    int64     `json:"size"`
	ModTime time.Time `json:"mod_time"`
}

type workspaceListResponse struct {
	CurrentPath string           `json:"currentPath"`
	Entries     []workspaceEntry `json:"entries"`
}

type workspaceListRequest struct {
	AppName string `json:"app_name" binding:"required"`
	Path    string `json:"path"`
}

type workspaceFileRequest struct {
	AppName string `json:"app_name" binding:"required"`
	Path    string `json:"path" binding:"required"`
}

type workspaceSaveRequest struct {
	AppName string `json:"app_name" binding:"required"`
	Path    string `json:"path" binding:"required"`
	Content string `json:"content"`
}

type workspaceDeleteRequest struct {
	AppName string `json:"app_name" binding:"required"`
	Path    string `json:"path" binding:"required"`
}

type workspaceCreateDirRequest struct {
	AppName string `json:"app_name" binding:"required"`
	Path    string `json:"path" binding:"required"`
}

func (h *BaseHandler) WorkspaceListHandler(ctx context.Context, c *app.RequestContext) {
	var req workspaceListRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.Error(err)
		return
	}

	_, targetPath, relPath, err := h.resolveWorkspacePath(req.AppName, req.Path)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	if _, err := os.Stat(targetPath); os.IsNotExist(err) {
		if err := os.MkdirAll(targetPath, 0o755); err != nil {
			c.Error(err)
			return
		}
	}

	dirEntries, err := os.ReadDir(targetPath)
	if err != nil {
		c.Error(err)
		return
	}

	entries := make([]workspaceEntry, 0, len(dirEntries))
	for _, entry := range dirEntries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		entryPath := entry.Name()
		if relPath != "" {
			entryPath = filepath.Join(relPath, entryPath)
		}

		entries = append(entries, workspaceEntry{
			Name:    entry.Name(),
			Path:    toWorkspacePath(entryPath),
			IsDir:   entry.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime(),
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir == entries[j].IsDir {
			return entries[i].Name < entries[j].Name
		}

		return entries[i].IsDir
	})

	c.JSON(http.StatusOK, SuccessResponse(workspaceListResponse{
		CurrentPath: toWorkspacePath(relPath),
		Entries:     entries,
	}))
}

func (h *BaseHandler) WorkspaceReadHandler(ctx context.Context, c *app.RequestContext) {
	var req workspaceFileRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.Error(err)
		return
	}

	_, filePath, _, err := h.resolveWorkspacePath(req.AppName, req.Path)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	info, err := os.Stat(filePath)
	if err != nil {
		c.Error(err)
		return
	}

	if info.IsDir() {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "path is a directory"))
		return
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(map[string]string{
		"content": string(data),
	}))
}

func (h *BaseHandler) WorkspaceSaveHandler(ctx context.Context, c *app.RequestContext) {
	var req workspaceSaveRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.Error(err)
		return
	}

	_, filePath, _, err := h.resolveWorkspacePath(req.AppName, req.Path)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		c.Error(err)
		return
	}

	if err := os.WriteFile(filePath, []byte(req.Content), 0o644); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, EmptyResponse())
}

func (h *BaseHandler) WorkspaceDeleteHandler(ctx context.Context, c *app.RequestContext) {
	var req workspaceDeleteRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.Error(err)
		return
	}

	_, targetPath, relPath, err := h.resolveWorkspacePath(req.AppName, req.Path)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	if relPath == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "cannot delete workspace root"))
		return
	}

	if err := os.RemoveAll(targetPath); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, EmptyResponse())
}

func (h *BaseHandler) WorkspaceMkdirHandler(ctx context.Context, c *app.RequestContext) {
	var req workspaceCreateDirRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.Error(err)
		return
	}

	_, dirPath, relPath, err := h.resolveWorkspacePath(req.AppName, req.Path)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	if relPath == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "directory name is required"))
		return
	}

	if err := os.MkdirAll(dirPath, 0o755); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, EmptyResponse())
}

func (h *BaseHandler) HandleWorkspaceUpload(ctx context.Context, c *app.RequestContext) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "failed to read upload file"))
		return
	}

	appName := strings.TrimSpace(string(c.FormValue("app_name")))
	if appName == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "app_name is required"))
		return
	}

	pathValue := string(c.FormValue("path"))
	root, targetDir, relDir, err := h.resolveWorkspacePath(appName, pathValue)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, err.Error()))
		return
	}

	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		c.Error(err)
		return
	}

	src, err := file.Open()
	if err != nil {
		c.Error(err)
		return
	}
	defer src.Close()

	safeName := filepath.Base(file.Filename)
	destPath := filepath.Join(targetDir, safeName)
	if !isPathWithinBase(root, destPath) {
		c.JSON(http.StatusBadRequest, ErrorResponse(http.StatusBadRequest, "invalid destination path"))
		return
	}

	dst, err := os.Create(destPath)
	if err != nil {
		c.Error(err)
		return
	}
	defer dst.Close()

	written, err := io.Copy(dst, src)
	if err != nil {
		os.Remove(destPath)
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(map[string]interface{}{
		"filename": safeName,
		"path":     toWorkspacePath(filepath.Join(relDir, safeName)),
		"size":     written,
	}))
}

func (h *BaseHandler) resolveWorkspacePath(appName, rel string) (string, string, string, error) {
	root, err := h.ensureWorkspaceRoot(appName)
	if err != nil {
		return "", "", "", err
	}

	cleanRel := cleanWorkspaceRelativePath(rel)

	targetPath := root
	if cleanRel != "" {
		targetPath = filepath.Join(root, cleanRel)
	}

	if !isPathWithinBase(root, targetPath) {
		return "", "", "", errors.New("invalid workspace path")
	}

	return root, targetPath, cleanRel, nil
}

func (h *BaseHandler) ensureWorkspaceRoot(appName string) (string, error) {
	appName = strings.TrimSpace(appName)
	if appName == "" {
		return "", errors.New("app name is required")
	}

	if strings.Contains(appName, "..") || strings.ContainsAny(appName, `/\`) {
		return "", errors.New("invalid app name")
	}

	root := filepath.Join(h.WorkSpaceDataPath(), appName)
	if err := os.MkdirAll(root, 0o755); err != nil {
		return "", err
	}

	return root, nil
}

func cleanWorkspaceRelativePath(p string) string {
	p = strings.TrimSpace(p)
	if p == "" {
		return ""
	}

	cleanPath := filepath.Clean(p)
	if cleanPath == "." {
		return ""
	}

	cleanPath = strings.TrimPrefix(cleanPath, string(os.PathSeparator))
	cleanPath = strings.TrimPrefix(cleanPath, "./")
	return cleanPath
}

func isPathWithinBase(base, target string) bool {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return false
	}

	return rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))
}

func toWorkspacePath(p string) string {
	if p == "" {
		return ""
	}

	return strings.ReplaceAll(p, string(os.PathSeparator), "/")
}
