package deploypipe

import (
	"context"
	"errors"
	"fmt"
	"path"

	"github.com/benlocal/lai-panel/pkg/node"
)

const (
	DockerComposeFilePath = "docker_compose_file_path"
)

type DockerComposeUpPipeline struct {
}

func (p *DockerComposeUpPipeline) Process(ctx context.Context, c *DeployCtx) (*DeployCtx, error) {
	if c.dockerComposeFile == nil {
		return c, errors.New("docker compose file is not found")
	}

	// write docker compose file to disk
	pa := path.Join(c.GetServicePath(), "docker_compose.yml")
	exec, err := c.NodeState.GetExec()
	if err != nil {
		return c, err
	}
	composeCmd, err := findDockerComposeCommand(exec)
	if err != nil {
		return c, err
	}
	err = exec.WriteFile(pa, []byte(*c.dockerComposeFile))
	if err != nil {
		return c, err
	}

	c.Send("info", "docker compose file written to disk, path: "+pa)
	c.Send("info", "  --> deploying to node: "+c.NodeState.GetNodeInfo())

	// execute docker compose up
	cmd := fmt.Sprintf("%s -f %s up -d --build", composeCmd, pa)
	c.Send("info", "executing command: "+cmd)
	err = exec.ExecuteCommand(cmd, c.env, func(s string) {
		c.Send("info", s)
	}, func(s string) {
		c.Send("error", s)
	})
	if err != nil {
		return c, err
	}

	c.Send("info", "docker compose up executed")
	c.deployInfo[DockerComposeFilePath] = pa

	return c, nil
}

func (p *DockerComposeUpPipeline) Cancel(c *DeployCtx, err error) {
	// do nothing
}

type DockerComposeDownPipeline struct {
}

func (p *DockerComposeDownPipeline) Process(ctx context.Context, c *DownCtx) (*DownCtx, error) {
	pa := c.deployInfo[DockerComposeFilePath]
	env := map[string]string{}
	exec, err := c.NodeState.GetExec()
	if err != nil {
		return c, err
	}
	composeCmd, err := findDockerComposeCommand(exec)
	if err != nil {
		return c, err
	}
	err = exec.ExecuteCommand(fmt.Sprintf("%s -f %s down", composeCmd, pa), env, func(s string) {
		// do nothing
	}, func(s string) {
		// do nothing
	})
	if err != nil {
		return c, err
	}
	return c, nil
}

func (p *DockerComposeDownPipeline) Cancel(c *DownCtx, err error) {
	// do nothing
}

func findDockerComposeCommand(exec node.NodeExec) (string, error) {
	if _, _, err := exec.ExecuteOutput("docker compose version", map[string]string{}); err == nil {
		return "docker compose", nil
	}
	if _, _, err := exec.ExecuteOutput("docker-compose version", map[string]string{}); err == nil {
		return "docker-compose", nil
	}
	return "", errors.New("docker compose command not found (tried 'docker compose' and 'docker-compose')")
}
