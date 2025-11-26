package deploypipe

import (
	"context"
	"fmt"
	"os"

	"github.com/benlocal/lai-panel/pkg/node"
)

type CleanupWorkspacePipeline struct {
}

func (p *CleanupWorkspacePipeline) Process(ctx context.Context, c *DeployCtx) (*DeployCtx, error) {
	installerPath, err := c.GetServicePath()
	if err != nil {
		return c, err
	}

	exec, err := c.NodeState.GetExec()
	if err != nil {
		return c, err
	}

	err = p.rm(exec, installerPath, c)
	if err != nil {
		return c, err
	}

	c.Send("info", fmt.Sprintf("workspace cleaned up: %s", installerPath))
	err = p.mkdir(exec, installerPath, c)
	if err != nil {
		return c, err
	}
	c.Send("info", fmt.Sprintf("workspace created: %s", installerPath))

	return c, nil
}

func (p *CleanupWorkspacePipeline) Cancel(c *DeployCtx, err error) {
	// do nothing
}

func (p *CleanupWorkspacePipeline) mkdir(exec node.NodeExec, path string, c *DeployCtx) error {
	cmd := fmt.Sprintf("mkdir -p %s", path)
	opt := node.NewNodeExecuteCommandOptions()
	opt.SetEnv(c.env)
	return exec.ExecuteCommand(cmd, opt, func(s string) {
		c.Send("info", s)
	}, func(s string) {
		c.Send("error", s)
	})
}

func (p *CleanupWorkspacePipeline) rm(exec node.NodeExec, path string, c *DeployCtx) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	if path == "/" {
		return fmt.Errorf("cannot remove root directory")
	}

	cmd := fmt.Sprintf("rm -rf %s", path)
	opt := node.NewNodeExecuteCommandOptions()
	opt.SetEnv(c.env)
	return exec.ExecuteCommand(cmd, opt, func(s string) {
		c.Send("info", s)
	}, func(s string) {
		c.Send("error", s)
	})
}
