package deploypipe

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
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

	installerPath, err := c.GetServicePath()
	if err != nil {
		return c, err
	}
	exec, err := c.NodeState.GetExec()
	if err != nil {
		return c, err
	}

	err = p.mkdir(exec, installerPath, c)
	if err != nil {
		return c, fmt.Errorf("failed to create installer path: %w", err)
	}

	filename, reader, err := p.downloadFile(exec, *staticPath, c)
	if err != nil {
		return c, fmt.Errorf("failed to download file: %w", err)
	}
	defer reader.Close()

	if strings.HasSuffix(strings.ToLower(filename), ".tar.gz") || strings.HasSuffix(strings.ToLower(filename), ".tgz") {
		c.Send("info", fmt.Sprintf("extracting tar.gz file to %s", installerPath))
		err = p.extractTarGz(reader, installerPath, exec, c)
		if err != nil {
			return c, fmt.Errorf("failed to extract tar.gz file: %w", err)
		}
		c.Send("info", "tar.gz file extracted successfully")
	} else {
		filePath := path.Join(installerPath, filename)
		err = exec.WriteFileStream(filePath, reader)
		if err != nil {
			return c, fmt.Errorf("failed to write file: %w", err)
		}
		c.Send("info", fmt.Sprintf("file saved to %s", filePath))
	}

	return c, nil
}

func (p *CopyInstallerPipeline) downloadFile(exec node.NodeExec, path string, c *DeployCtx) (string, io.ReadCloser, error) {
	var fileName string
	var reader io.ReadCloser

	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		c.Send("info", fmt.Sprintf("downloading file from %s", path))
		resp, err := http.Get(path)
		if err != nil {
			return "", nil, fmt.Errorf("failed to download file: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return "", nil, fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
		}

		fileName = filepath.Base(path)
		if fileName == "" || fileName == "/" {
			fileName = "downloaded_file"
		}
		reader = resp.Body
	} else {
		c.Send("info", fmt.Sprintf("reading local file from %s", path))
		fileName = filepath.Base(path)
		r, err := os.Open(path)
		if err != nil {
			return "", nil, fmt.Errorf("failed to open local file: %w", err)
		}
		reader = r
	}

	return fileName, reader, nil
}

func (p *CopyInstallerPipeline) mkdir(exec node.NodeExec, path string, c *DeployCtx) error {
	cmd := fmt.Sprintf("mkdir -p %s", path)
	return exec.ExecuteCommand(cmd, c.env, func(s string) {
		c.Send("info", s)
	}, func(s string) {
		c.Send("error", s)
	})
}

func (p *CopyInstallerPipeline) extractTarGz(reader io.ReadCloser, destDir string, exec node.NodeExec, c *DeployCtx) error {
	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()
	tarReader := tar.NewReader(gzReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar header: %w", err)
		}
		targetPath := path.Join(destDir, header.Name)
		if !strings.HasPrefix(filepath.Clean(targetPath), filepath.Clean(destDir)) {
			c.Send("warning", fmt.Sprintf("skipping unsafe path: %s", header.Name))
			continue
		}

		switch header.Typeflag {
		case tar.TypeDir:
			cmd := fmt.Sprintf("mkdir -p %s", targetPath)
			err = exec.ExecuteCommand(cmd, map[string]string{}, func(s string) {
				c.Send("info", s)
			}, func(s string) {
				c.Send("error", s)
			})
			if err != nil {
				return fmt.Errorf("failed to create directory %s: %w", targetPath, err)
			}

		case tar.TypeReg:
			fileData, err := io.ReadAll(tarReader)
			if err != nil {
				return fmt.Errorf("failed to read file %s: %w", header.Name, err)
			}

			err = exec.WriteFile(targetPath, fileData)
			if err != nil {
				return fmt.Errorf("failed to write file %s: %w", targetPath, err)
			}

			if header.Mode > 0 {
				cmd := fmt.Sprintf("chmod %o %s", header.Mode, targetPath)
				_ = exec.ExecuteCommand(cmd, map[string]string{}, func(s string) {
					c.Send("info", s)
				}, func(s string) {
					c.Send("error", s)
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
