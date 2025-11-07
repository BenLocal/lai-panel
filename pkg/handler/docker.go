package handler

import (
	"errors"
	"strconv"

	"github.com/valyala/fasthttp"
)

func (h *BaseHandler) DockerInfo(ctx *fasthttp.RequestCtx) {
	nodeId, err := h.getNodeID(ctx)
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

func (h *BaseHandler) getNodeID(ctx *fasthttp.RequestCtx) (int64, error) {
	nodeIdStr := string(ctx.Request.Header.Peek("X-Node-ID"))
	nodeId, err := strconv.ParseInt(nodeIdStr, 10, 64)
	if err != nil {
		return 0, errors.New("invalid node id")
	}
	return nodeId, nil
}
