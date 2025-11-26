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

type DownloadInstallerPipeline struct {
}

func (p *DownloadInstallerPipeline) Process(ctx context.Context, c *DeployCtx) (*DeployCtx, error) {
	staticPath := c.App.StaticPath
	if staticPath == nil || *staticPath == "" {
		return c, nil
	}

	exec, err := c.NodeState.GetExec()
	if err != nil {
		return c, err
	}

	installerPath, err := c.GetServicePath()
	if err != nil {
		return c, err
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

func (p *DownloadInstallerPipeline) downloadFile(exec node.NodeExec, path string, c *DeployCtx) (string, io.ReadCloser, error) {
	var fileName string
	var reader io.ReadCloser

	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		c.Send("info", fmt.Sprintf("downloading file from %s", path))
		resp, err := http.Get(path)
		if err != nil {
			return "", nil, fmt.Errorf("failed to download file: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
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

func (p *DownloadInstallerPipeline) extractTarGz(reader io.ReadCloser, destDir string, exec node.NodeExec, c *DeployCtx) error {
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
			opt := node.NewNodeExecuteCommandOptions()
			err = exec.ExecuteCommand(cmd, opt, func(s string) {
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
				opt := node.NewNodeExecuteCommandOptions()
				_ = exec.ExecuteCommand(cmd, opt, func(s string) {
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

func (p *DownloadInstallerPipeline) Cancel(c *DeployCtx, err error) {
	// do nothing
}
