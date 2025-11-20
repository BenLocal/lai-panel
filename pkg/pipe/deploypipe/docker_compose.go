package deploypipe

import (
	"context"
	"errors"
	"fmt"
	"path"
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

	name := c.Service.Name

	// write docker compose file to disk
	pa := path.Join(c.options.ServicePath(), name, "docker_compose.yml")
	err := c.NodeState.Exec.WriteFile(pa, []byte(*c.dockerComposeFile))
	if err != nil {
		return c, err
	}

	c.Send("info", "docker compose file written to disk, path: "+pa)
	c.Send("info", "  --> deploying to node: "+c.NodeState.GetNodeInfo())

	// execute docker compose up
	cmd := fmt.Sprintf("docker compose -f %s up -d --build", pa)
	err = c.NodeState.Exec.ExecuteCommand(cmd, c.env, func(s string) {
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
	err := c.NodeState.Exec.ExecuteCommand(fmt.Sprintf("docker compose -f %s down", pa), env, func(s string) {
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
