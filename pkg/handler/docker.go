package handler

import (
	"errors"
	"strconv"

	"github.com/benlocal/lai-panel/pkg/di"
	"github.com/benlocal/lai-panel/pkg/node"
	"github.com/benlocal/lai-panel/pkg/repository"
	"github.com/valyala/fasthttp"
)

type DockerHandler struct {
	nodeManager *node.NodeManager
	nodeRepo    *repository.NodeRepository
}

func NewDockerHandler(nodeManager *node.NodeManager,
	nodeRepo *repository.NodeRepository,
) *DockerHandler {
	return &DockerHandler{
		nodeManager: nodeManager,
		nodeRepo:    nodeRepo,
	}
}

func HandleDockerInfoWithDI(ctx *fasthttp.RequestCtx) {
	err := di.Invoke(func(h *DockerHandler) {
		h.Info(ctx)
	})
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(err.Error())
	}
}

func (h *DockerHandler) Info(ctx *fasthttp.RequestCtx) {
	nodeId, err := h.getNodeID(ctx)
	if err != nil {
		JSONError(ctx, "invalid node id", err)
		return
	}
	node, err := h.nodeRepo.GetByID(nodeId)
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

func (h *DockerHandler) getNodeID(ctx *fasthttp.RequestCtx) (int64, error) {
	nodeIdStr := string(ctx.Request.Header.Peek("X-Node-ID"))
	nodeId, err := strconv.ParseInt(nodeIdStr, 10, 64)
	if err != nil {
		return 0, errors.New("invalid node id")
	}
	return nodeId, nil
}
