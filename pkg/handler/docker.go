package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/benlocal/lai-panel/pkg/node"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/sse"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/volume"
)

func (h *BaseHandler) DockerInfo(ctx context.Context, c *app.RequestContext) {
	nodeState, err := h.getNodeState(ctx, c)
	if err != nil {
		c.Error(err)
		return
	}

	info, err := nodeState.DockerClient.Info(ctx)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(info))
}

func (h *BaseHandler) DockerContainers(ctx context.Context, c *app.RequestContext) {
	nodeState, err := h.getNodeState(ctx, c)
	if err != nil {
		c.Error(err)
		return
	}
	containers, err := nodeState.DockerClient.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(containers))
}

func (h *BaseHandler) DockerImages(ctx context.Context, c *app.RequestContext) {
	nodeState, err := h.getNodeState(ctx, c)
	if err != nil {
		c.Error(err)
		return
	}
	images, err := nodeState.DockerClient.ImageList(ctx, image.ListOptions{})
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(images))
}

func (h *BaseHandler) DockerVolumes(ctx context.Context, c *app.RequestContext) {
	nodeState, err := h.getNodeState(ctx, c)
	if err != nil {
		c.Error(err)
		return
	}
	volumes, err := nodeState.DockerClient.VolumeList(ctx, volume.ListOptions{})
	if err != nil {
		c.Error(err)
		return
	}

	if len(volumes.Volumes) == 0 {
		c.JSON(http.StatusOK, SuccessResponse([]*volume.Volume{}))
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(volumes.Volumes))
}

func (h *BaseHandler) DockerNetworks(ctx context.Context, c *app.RequestContext) {
	nodeState, err := h.getNodeState(ctx, c)
	if err != nil {
		c.Error(err)
		return
	}
	networks, err := nodeState.DockerClient.NetworkList(ctx, network.ListOptions{})
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, networks)
}

func (h *BaseHandler) getNodeState(_ context.Context, c *app.RequestContext) (*node.NodeState, error) {
	nodeId, err := h.getNodeIDFromRequest(c)
	if err != nil {
		return nil, err
	}
	node, err := h.NodeRepository().GetByID(nodeId)
	if err != nil {
		return nil, err
	}

	nodeState, err := h.NodeManager().AddOrGetNode(node)
	if err != nil {
		return nil, err
	}
	return nodeState, nil
}

func (h *BaseHandler) getNodeIDFromRequest(c *app.RequestContext) (int64, error) {
	nodeIdStr := string(c.Request.Header.Peek("X-Node-ID"))
	nodeId, err := strconv.ParseInt(nodeIdStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return nodeId, nil
}

func (h *BaseHandler) DockerList(ctx context.Context, c *app.RequestContext) {
	nodeId, err := h.getNodeIDFromRequest(c)
	if err != nil {
		c.Error(err)
		return
	}
	node, err := h.NodeRepository().GetByID(nodeId)
	if err != nil {
		c.Error(err)
		return
	}
	nodeState, err := h.NodeManager().AddOrGetNode(node)
	if err != nil {
		c.Error(err)
		return
	}
	containers, err := nodeState.DockerClient.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, containers)
}

func (h *BaseHandler) DockerImagePullAuto(ctx context.Context, c *app.RequestContext) {
	type dockerImagePullAutoRequest struct {
		Image string `json:"image"`
	}
	var req dockerImagePullAutoRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.Error(err)
		return
	}
	imageString := req.Image
	if imageString == "" {
		c.Error(errors.New("image is required"))
		return
	}
	nodeId, err := h.getNodeIDFromRequest(c)
	if err != nil {
		c.Error(err)
		return
	}
	dst, err := h.NodeRepository().GetByID(nodeId)
	if err != nil {
		c.Error(err)
		return
	}
	ds, err := h.NodeManager().AddOrGetNode(dst)
	if err != nil {
		c.Error(err)
		return
	}

	// select other nodes have the same image
	nodes, err := h.NodeRepository().List()
	if err != nil {
		c.Error(err)
		return
	}

	var ss *node.NodeState
	for _, srcNode := range nodes {
		if srcNode.ID == nodeId {
			continue
		}
		srcNodeState, err := h.NodeManager().AddOrGetNode(&srcNode)
		if err != nil {
			continue
		}
		_, err = srcNodeState.DockerClient.ImageList(ctx, image.ListOptions{})
		if err != nil {
			continue
		}
		ss = srcNodeState
		break
	}

	if ss == nil {
		c.Error(errors.New("no node has the image"))
		return
	}

	writer := sse.NewWriter(c)
	defer writer.Close()
	err = node.CopyImageBetweenNodes(ctx, ss, ds, imageString, func(ctx context.Context, reader io.ReadCloser) error {
		_, err := io.Copy(&CopyWriter{writer}, reader)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		writer.WriteEvent("", "error", []byte(err.Error()))
		return
	}
	writer.WriteEvent("", "done", []byte("done"))

}

type CopyWriter struct {
	*sse.Writer
}

func (c *CopyWriter) Write(p []byte) (n int, err error) {
	err = c.Writer.WriteEvent("", "data", p)
	if err != nil {
		return 0, err
	}
	return len(p), err
}
