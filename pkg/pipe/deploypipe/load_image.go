package deploypipe

import (
	"context"
	"errors"
	"io"

	"github.com/benlocal/lai-panel/pkg/node"
	"github.com/cloudwego/hertz/pkg/protocol/sse"
	"gopkg.in/yaml.v3"
)

type LoadImagePipeline struct {
}

func (p *LoadImagePipeline) Process(ctx context.Context, c *DeployCtx) (*DeployCtx, error) {
	if c.dockerComposeFile == nil {
		return c, nil
	}

	images, err := p.getImages(*c.dockerComposeFile)
	if err != nil {
		return c, nil
	}

	if len(images) == 0 {
		return c, nil
	}

	for _, image := range images {
		err := p.loadImage(ctx, c, image)
		if err != nil {
			c.Send("warning", "load image "+image+" failed: "+err.Error())
			continue
		}
	}

	return c, nil
}

func (p *LoadImagePipeline) Cancel(c *DeployCtx, err error) {
	// do nothing
}

func (p *LoadImagePipeline) getImages(dockerComposeFile string) ([]string, error) {
	var y yaml.Node
	if err := yaml.Unmarshal([]byte(dockerComposeFile), &y); err != nil {
		return nil, err
	}

	root := &y
	if y.Kind == yaml.DocumentNode && len(y.Content) > 0 {
		root = y.Content[0]
	}

	services := lookup(root, "services")
	if services == nil || services.Kind != yaml.MappingNode {
		return nil, nil
	}

	// 遍历所有 services，找出不需要 build 的 image
	images := []string{}
	for i := 0; i < len(services.Content); i += 2 {
		if i+1 >= len(services.Content) {
			continue
		}
		svc := services.Content[i+1]

		buildNode := lookup(svc, "build")
		if buildNode != nil {
			continue
		}

		imageNode := lookup(svc, "image")
		if imageNode != nil && imageNode.Kind == yaml.ScalarNode && imageNode.Value != "" {
			images = append(images, imageNode.Value)
		}
	}

	return images, nil
}

func (p *LoadImagePipeline) loadImage(ctx context.Context,
	c *DeployCtx,
	image string) error {
	// load image to local registry
	currentState := c.NodeState
	dc, err := currentState.GetDockerClient()
	if err != nil {
		return err
	}

	// check if image exists
	_, err = dc.ImageInspect(ctx, image)
	if err == nil {
		// image already exists
		c.Send("info", "image "+image+" already exists")
		return nil
	}

	// load image to local registry
	nodes, err := c.appCtx.NodeRepository().List()
	if err != nil {
		return err
	}

	var ss *node.NodeState
	for _, node := range nodes {
		if node.ID == currentState.GetNodeID() {
			continue
		}
		ssState, err := c.appCtx.NodeManager().GetNodeState(node.ID)
		if err != nil {
			continue
		}
		dc, err := ssState.GetDockerClient()
		if err != nil {
			continue
		}
		_, err = dc.ImageInspect(ctx, image)
		if err != nil {
			continue
		}

		ss = ssState
		break
	}

	if ss == nil {
		return errors.New("no node has the image")
	}

	return node.CopyImageBetweenNodes(ctx, ss, currentState, image, func(ctx context.Context, reader io.ReadCloser) error {
		_, err := io.Copy(&CopyWriter{c.writer}, reader)
		if err != nil {
			return err
		}
		c.Send("info", "load image "+image+" form "+ss.GetNodeInfo()+" to local node success")
		return nil
	})
}

type CopyWriter struct {
	*sse.Writer
}

func (c *CopyWriter) Write(p []byte) (n int, err error) {
	err = c.Writer.WriteEvent("", "info", p)
	if err != nil {
		return 0, err
	}
	return len(p), err
}
