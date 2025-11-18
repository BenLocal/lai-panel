package handler

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/http1/resp"
)

func (h *BaseHandler) StaticDataPath() string {
	return h.options.StaticDataPath()
}

func (h *BaseHandler) Tar(ctx context.Context, c *app.RequestContext) {
	p := c.Param("path")
	filePath, err := h.validateTarPath(p)
	if err != nil {
		c.Error(err)
		return
	}

	// Get the directory name for the tar filename
	dirName := filepath.Base(filePath)
	if dirName == "" || dirName == "." {
		dirName = "archive"
	}

	// Set response headers for file download
	h.setTarHeaders(c, dirName)

	// Hijack the writer to use ChunkedBodyWriter for chunked transfer
	c.Response.HijackWriter(resp.NewChunkedBodyWriter(&c.Response, c.GetWriter()))

	// Create tar archive
	if err := h.createTarArchive(c, filePath); err != nil {
		c.Error(fmt.Errorf("failed to create tar archive: %w", err))
		return
	}
}

// validateTarPath validates the path parameter and returns the full file path
func (h *BaseHandler) validateTarPath(p string) (string, error) {
	if p == "" {
		return "", errors.New("path is required")
	}

	filePath := path.Join(h.StaticDataPath(), p)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", errors.New("file not found")
	}

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", err
	}

	if !fileInfo.IsDir() {
		return "", errors.New("file is not a directory")
	}

	return filePath, nil
}

// setTarHeaders sets the HTTP headers for tar file download
func (h *BaseHandler) setTarHeaders(c *app.RequestContext, dirName string) {
	c.Header("Content-Type", "application/gzip")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.tar.gz", dirName))
	c.Header("X-Content-Type-Options", "nosniff")
}

// createTarArchive creates a tar.gz archive from the given directory using chunked transfer
func (h *BaseHandler) createTarArchive(c *app.RequestContext, filePath string) error {
	// Create gzip writer with compression level
	gzWriter := gzip.NewWriter(c.Response.BodyWriter())
	defer gzWriter.Close()

	// Create tar writer
	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	// Walk through the directory and add files to tar
	err := filepath.Walk(filePath, func(file string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if err := h.addToTar(tarWriter, filePath, file, info); err != nil {
			return err
		}

		// Flush tar writer to send chunk
		if err := tarWriter.Flush(); err != nil {
			return err
		}

		// Flush gzip writer to send chunk
		if err := gzWriter.Flush(); err != nil {
			return err
		}

		// Flush response writer to send chunk to client using ChunkedBodyWriter
		if err := c.Flush(); err != nil {
			return err
		}

		return nil
	})

	return err
}

// addToTar adds a file or directory to the tar archive
func (h *BaseHandler) addToTar(tarWriter *tar.Writer, basePath, file string, info os.FileInfo) error {
	// Get relative path from the base directory
	relPath, err := filepath.Rel(basePath, file)
	if err != nil {
		return err
	}

	// Skip the root directory itself
	if relPath == "." {
		return nil
	}

	// Create tar header
	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}

	// Set the name in the tar archive (use forward slashes for compatibility)
	header.Name = strings.ReplaceAll(relPath, "\\", "/")
	if info.IsDir() {
		header.Name += "/"
	}

	// Write header
	if err := tarWriter.WriteHeader(header); err != nil {
		return err
	}

	// If it's a file, copy its content
	if !info.IsDir() {
		return h.addFileToTar(tarWriter, file)
	}

	return nil
}

// addFileToTar adds a file's content to the tar archive using chunked transfer
func (h *BaseHandler) addFileToTar(tarWriter *tar.Writer, filePath string) error {
	fileHandle, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer fileHandle.Close()

	// Use buffer for chunked transfer
	buf := make([]byte, 32*1024) // 32KB buffer
	_, err = io.CopyBuffer(tarWriter, fileHandle, buf)
	return err
}
