package handler

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/benlocal/lai-panel/pkg/node"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
)

func (h *BaseHandler) DockerInfo(ctx context.Context, c *app.RequestContext) {
	nodeId, err := h.getNodeIDFromRequest(c)
	if err != nil {
		c.Error(err)
		return
	}
	node, err := h.nodeRepository.GetByID(nodeId)
	if err != nil {
		c.Error(err)
		return
	}
	nodeState, err := h.nodeManager.AddOrGetNode(node)
	if err != nil {
		c.Error(err)
		return
	}

	info, err := nodeState.DockerClient.Info(ctx)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, info)
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
	node, err := h.nodeRepository.GetByID(nodeId)
	if err != nil {
		c.Error(err)
		return
	}
	nodeState, err := h.nodeManager.AddOrGetNode(node)
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
	dst, err := h.nodeRepository.GetByID(nodeId)
	if err != nil {
		c.Error(err)
		return
	}
	_, err = h.nodeManager.AddOrGetNode(dst)
	if err != nil {
		c.Error(err)
		return
	}

	// select other nodes have the same image
	nodes, err := h.nodeRepository.List()
	if err != nil {
		c.Error(err)
		return
	}

	var ss *node.NodeState
	for _, srcNode := range nodes {
		if srcNode.ID == nodeId {
			continue
		}
		srcNodeState, err := h.nodeManager.AddOrGetNode(&srcNode)
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

	c.SetStatusCode(http.StatusOK)
	c.Response.Header.Set("Content-Type", "text/event-stream")
	c.Response.Header.Set("Cache-Control", "no-cache")
	c.Response.Header.Set("Connection", "keep-alive")

	// c.SetBodyStream(io.Reader(func(writer *bufio.Writer) {
	// 	err := node.CopyImageBetweenNodes(ctx, ss, ds, imageString, func(ctx context.Context, reader io.ReadCloser) error {
	// 		_, err := io.Copy(writer, reader)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		writer.Flush()
	// 		return nil
	// 	})
	// 	if err != nil {
	// 		writer.WriteString("error\n")
	// 		writer.Flush()
	// 		return
	// 	}
	// 	writer.WriteString("done\n")
	// 	writer.Flush()
	// }, 0), 0)

}
