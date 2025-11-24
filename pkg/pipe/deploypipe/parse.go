package deploypipe

import (
	"bytes"
	"context"
	"errors"

	"github.com/benlocal/lai-panel/pkg/tmpl"
	"gopkg.in/yaml.v3"
)

const (
	ManagedByLabel = "com.lai-panel.managed-by"
	OwnerLabel     = "com.lai-panel.owner"
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

	v, err = p.editFile(v, map[string]string{
		ManagedByLabel: "lai-panel",
		OwnerLabel:     "lai-panel",
	})
	if err != nil {
		return c, err
	}

	c.dockerComposeFile = &v

	c.Send("info", *c.dockerComposeFile)
	return c, nil
}

func (p *DockerComposeFileParsePipeline) Cancel(c *DeployCtx, err error) {
	// do nothing
}

func (p *DockerComposeFileParsePipeline) editFile(
	file string,
	labels map[string]string) (string, error) {
	var y yaml.Node
	if err := yaml.Unmarshal([]byte(file), &y); err != nil {
		return file, err
	}

	// Get the root document node (first child of DocumentNode)
	root := &y
	if y.Kind == yaml.DocumentNode && len(y.Content) > 0 {
		root = y.Content[0]
	}

	services := lookup(root, "services")
	if services == nil || services.Kind != yaml.MappingNode {
		return file, errors.New("no 'services' map found")
	}
	for i := 0; i < len(services.Content); i += 2 {
		svc := services.Content[i+1]
		labelsNode := lookup(svc, "labels")
		if labelsNode == nil {
			labelsNode = &yaml.Node{Kind: yaml.MappingNode}
			svc.Content = append(svc.Content,
				&yaml.Node{Kind: yaml.ScalarNode, Value: "labels"},
				labelsNode,
			)
		}
		for k, v := range labels {
			setLabel(labelsNode, k, v)
		}
	}

	// Marshal back to string
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)
	defer enc.Close()

	if err := enc.Encode(&y); err != nil {
		return file, err
	}

	return buf.String(), nil
}

func lookup(node *yaml.Node, key string) *yaml.Node {
	if node.Kind != yaml.MappingNode {
		return nil
	}
	for i := 0; i < len(node.Content); i += 2 {
		if i+1 < len(node.Content) && node.Content[i].Value == key {
			return node.Content[i+1]
		}
	}
	return nil
}

func setLabel(labelsNode *yaml.Node, k, v string) {
	for i := 0; i < len(labelsNode.Content); i += 2 {
		if labelsNode.Content[i].Value == k {
			labelsNode.Content[i+1].Value = v
			return
		}
	}
	labelsNode.Content = append(labelsNode.Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: k},
		&yaml.Node{Kind: yaml.ScalarNode, Value: v},
	)
}
