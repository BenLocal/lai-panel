package handler

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strconv"

	"github.com/benlocal/lai-panel/pkg/node"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/valyala/fasthttp"
)

func (h *BaseHandler) DockerInfo(ctx *fasthttp.RequestCtx) {
	nodeId, err := h.getNodeIDFromRequest(ctx)
	if err != nil {
		JSONError(ctx, "invalid node id", err)
		return
	}
	node, err := h.nodeRepository.GetByID(nodeId)
	if err != nil {
		JSONError(ctx, "node not found", err)
		return
	}
	nodeState, err := h.nodeManager.AddOrGetNode(node)
	if err != nil {
		JSONError(ctx, "node not found", err)
		return
	}

	info, err := nodeState.DockerClient.Info(ctx)
	if err != nil {
		JSONError(ctx, "failed to get docker info", err)
		return
	}
	JSONSuccess(ctx, info)
}

func (h *BaseHandler) getNodeIDFromRequest(ctx *fasthttp.RequestCtx) (int64, error) {
	nodeIdStr := string(ctx.Request.Header.Peek("X-Node-ID"))
	nodeId, err := strconv.ParseInt(nodeIdStr, 10, 64)
	if err != nil {
		return 0, errors.New("invalid node id")
	}
	return nodeId, nil
}

func (h *BaseHandler) DockerList(ctx *fasthttp.RequestCtx) {
	nodeId, err := h.getNodeIDFromRequest(ctx)
	if err != nil {
		JSONError(ctx, "invalid node id", err)
		return
	}
	node, err := h.nodeRepository.GetByID(nodeId)
	if err != nil {
		JSONError(ctx, "node not found", err)
		return
	}
	nodeState, err := h.nodeManager.AddOrGetNode(node)
	if err != nil {
		JSONError(ctx, "node not found", err)
		return
	}
	containers, err := nodeState.DockerClient.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		JSONError(ctx, "failed to get docker containers", err)
		return
	}
	JSONSuccess(ctx, containers)
}

func (h *BaseHandler) DockerImagePullAuto(ctx *fasthttp.RequestCtx) {
	type dockerImagePullAutoRequest struct {
		Image string `json:"image"`
	}
	var req dockerImagePullAutoRequest
	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		JSONError(ctx, "invalid request", err)
		return
	}
	imageString := req.Image
	if imageString == "" {
		JSONError(ctx, "image is required", nil)
		return
	}
	nodeId, err := h.getNodeIDFromRequest(ctx)
	if err != nil {
		JSONError(ctx, "invalid node id", err)
		return
	}
	dst, err := h.nodeRepository.GetByID(nodeId)
	if err != nil {
		JSONError(ctx, "node not found", err)
		return
	}
	ds, err := h.nodeManager.AddOrGetNode(dst)
	if err != nil {
		JSONError(ctx, "node not found", err)
		return
	}

	// select other nodes have the same image
	nodes, err := h.nodeRepository.List()
	if err != nil {
		JSONError(ctx, "failed to get nodes", err)
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
		JSONError(ctx, "no node has the image", nil)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.Set("Content-Type", "text/event-stream")
	ctx.Response.Header.Set("Cache-Control", "no-cache")
	ctx.Response.Header.Set("Connection", "keep-alive")

	ctx.SetBodyStreamWriter(func(writer *bufio.Writer) {
		err := node.CopyImageBetweenNodes(ctx, ss, ds, imageString, func(ctx context.Context, reader io.ReadCloser) error {
			_, err := io.Copy(writer, reader)
			if err != nil {
				return err
			}
			writer.Flush()
			return nil
		})
		if err != nil {
			writer.WriteString("error\n")
			writer.Flush()
			return
		}
		writer.WriteString("done\n")
		writer.Flush()
	})

}
