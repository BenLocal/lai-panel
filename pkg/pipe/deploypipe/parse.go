package deploypipe

import (
	"context"
	"errors"

	"github.com/benlocal/lai-panel/pkg/tmpl"
)

type DockerComposeFileParsePipeline struct {
}

func (p *DockerComposeFileParsePipeline) Process(ctx context.Context, c *DeployCtx) (*DeployCtx, error) {
	tpl := c.App.DockerCompose
	if tpl == nil {
		return c, errors.New("docker compose file is not found")
	}

	v, err := tmpl.ParseDockerCompose("docker compose", *tpl, c.env)
	if err != nil {
		return c, err
	}

	c.Send("info", "docker compose file parsed")
	c.Send("info", v)

	c.dockerComposeFile = &v
	return c, nil
}

func (p *DockerComposeFileParsePipeline) Cancel(c *DeployCtx, err error) {
	// do nothing
}
