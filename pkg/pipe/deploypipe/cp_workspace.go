package deploypipe

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/benlocal/lai-panel/pkg/node"
	"github.com/benlocal/lai-panel/pkg/options"
	"github.com/benlocal/lai-panel/pkg/tmpl"
)

type CopyWorkspacePipeline struct {
}

func (p *CopyWorkspacePipeline) Process(ctx context.Context, c *DeployCtx) (*DeployCtx, error) {
	workspace := path.Join(c.options.DataPath(), options.WORK_SPACE_BASE_PATH)
	appws := path.Join(workspace, c.App.Name)
	_, err := os.Stat(appws)
	if err != nil {
		return c, err
	}
	if os.IsNotExist(err) {
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

	// walk directory and copy to appws
	err = filepath.Walk(appws, func(filePath string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		// Calculate relative path from appws
		relPath, err := filepath.Rel(appws, filePath)
		if err != nil {
			return err
		}

		// Skip the root directory itself
		if relPath == "." {
			return nil
		}

		// Target path in installerPath
		targetPath := filepath.Join(installerPath, relPath)

		// If it's a directory, create it in the target location using exec
		if info.IsDir() {
			cmd := fmt.Sprintf("mkdir -p %s", targetPath)
			opt := node.NewNodeExecuteCommandOptions()
			opt.SetEnv(c.env)
			return exec.ExecuteCommand(cmd, opt, func(s string) {
				c.Send("info", s)
			}, func(s string) {
				c.Send("error", s)
			})
		}

		// If it's a file, read, process with template, and write
		// Read file content from local filesystem
		content, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		// Process content with template
		processedContent, err := tmpl.ParseWithEnv(relPath, string(content), c.env)
		if err != nil {
			return err
		}

		// Ensure target directory exists using exec
		targetDir := filepath.Dir(targetPath)
		cmd := fmt.Sprintf("mkdir -p %s", targetDir)
		opt := node.NewNodeExecuteCommandOptions()
		opt.SetEnv(c.env)
		if err := exec.ExecuteCommand(cmd, opt, func(s string) {
			c.Send("info", s)
		}, func(s string) {
			c.Send("error", s)
		}); err != nil {
			return err
		}

		// Write processed content to target file using exec
		return exec.WriteFile(targetPath, []byte(processedContent))
	})

	if err != nil {
		return c, err
	}

	return c, nil
}

func (p *CopyWorkspacePipeline) Cancel(c *DeployCtx, err error) {
	// No cleanup needed for copy operation
}
