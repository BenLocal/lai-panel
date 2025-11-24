package deploypipe

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/benlocal/lai-panel/pkg/node"
)

type CopyInstallerPipeline struct {
}

func (p *CopyInstallerPipeline) Process(ctx context.Context, c *DeployCtx) (*DeployCtx, error) {
	staticPath := c.App.StaticPath
	if staticPath == nil || *staticPath == "" {
		return c, nil
	}

	installerPath := c.GetServicePath()
	exec, err := c.NodeState.GetExec()
	if err != nil {
		return c, err
	}

	cmd := fmt.Sprintf("mkdir -p %s", installerPath)
	err = exec.ExecuteCommand(cmd, c.env, func(s string) {
		c.Send("info", s)
	}, func(s string) {
		c.Send("error", s)
	})
	if err != nil {
		return c, fmt.Errorf("failed to create installer path: %w", err)
	}

	// 下载文件
	var fileData []byte
	var fileName string

	if strings.HasPrefix(*staticPath, "http://") || strings.HasPrefix(*staticPath, "https://") {
		// 从 URL 下载
		c.Send("info", fmt.Sprintf("downloading file from %s", *staticPath))
		resp, err := http.Get(*staticPath)
		if err != nil {
			return c, fmt.Errorf("failed to download file: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return c, fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
		}

		fileData, err = io.ReadAll(resp.Body)
		if err != nil {
			return c, fmt.Errorf("failed to read downloaded file: %w", err)
		}

		// 从 URL 中提取文件名
		fileName = filepath.Base(*staticPath)
		if fileName == "" || fileName == "/" {
			fileName = "downloaded_file"
		}
	} else {
		// 本地文件路径
		c.Send("info", fmt.Sprintf("reading local file from %s", *staticPath))
		fileData, err = exec.ReadFile(*staticPath)
		if err != nil {
			return c, fmt.Errorf("failed to read local file: %w", err)
		}
		fileName = filepath.Base(*staticPath)
	}

	// 检查是否是 tar.gz 文件
	if strings.HasSuffix(strings.ToLower(fileName), ".tar.gz") || strings.HasSuffix(strings.ToLower(fileName), ".tgz") {
		c.Send("info", fmt.Sprintf("extracting tar.gz file to %s", installerPath))
		err = p.extractTarGz(fileData, installerPath, exec, c)
		if err != nil {
			return c, fmt.Errorf("failed to extract tar.gz file: %w", err)
		}
		c.Send("info", "tar.gz file extracted successfully")
	} else {
		// 不是 tar.gz 文件，直接保存到 installerPath
		filePath := path.Join(installerPath, fileName)
		err = exec.WriteFile(filePath, fileData)
		if err != nil {
			return c, fmt.Errorf("failed to write file: %w", err)
		}
		c.Send("info", fmt.Sprintf("file saved to %s", filePath))
	}

	return c, nil
}

func (p *CopyInstallerPipeline) extractTarGz(data []byte, destDir string, exec node.NodeExec, c *DeployCtx) error {
	// 创建 gzip reader
	gzReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	// 创建 tar reader
	tarReader := tar.NewReader(gzReader)

	// 遍历 tar 文件中的所有文件
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}

		// 构建目标文件路径
		targetPath := path.Join(destDir, header.Name)

		// 检查路径安全性（防止路径遍历攻击）
		if !strings.HasPrefix(filepath.Clean(targetPath), filepath.Clean(destDir)) {
			c.Send("warning", fmt.Sprintf("skipping unsafe path: %s", header.Name))
			continue
		}

		switch header.Typeflag {
		case tar.TypeDir:
			// 创建目录
			cmd := fmt.Sprintf("mkdir -p %s", targetPath)
			err = exec.ExecuteCommand(cmd, map[string]string{}, func(s string) {
				// 静默处理
			}, func(s string) {
				// 静默处理
			})
			if err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}

		case tar.TypeReg:
			// 读取文件内容
			fileData, err := io.ReadAll(tarReader)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", header.Name, err)
			}

			// 写入文件
			err = exec.WriteFile(targetPath, fileData)
			if err != nil {
				return fmt.Errorf("failed to write file %s: %w", targetPath, err)
			}

			// 设置文件权限
			if header.Mode > 0 {
				cmd := fmt.Sprintf("chmod %o %s", header.Mode, targetPath)
				_ = exec.ExecuteCommand(cmd, map[string]string{}, func(s string) {
					// 静默处理输出
				}, func(s string) {
					// 静默处理错误
				})
			}

			c.Send("info", fmt.Sprintf("extracted: %s", header.Name))
		}
	}

	return nil
}

func (p *CopyInstallerPipeline) Cancel(c *DeployCtx, err error) {
	// do nothing
}
