package handler

import (
	"encoding/json"

	"github.com/benlocal/lai-panel/pkg/model"
	"github.com/valyala/fasthttp"
)

func (h *BaseHandler) AddNodeHandler(ctx *fasthttp.RequestCtx) {
	var node model.Node
	if err := json.Unmarshal(ctx.PostBody(), &node); err != nil {
		JSONError(ctx, "Invalid JSON", err)
		return
	}

	if err := h.nodeRepository.Create(&node); err != nil {
		JSONError(ctx, "Failed to create node", err)
		return
	}

	JSONSuccess(ctx, node)
}

func (h *BaseHandler) GetNodeHandler(ctx *fasthttp.RequestCtx) {
	type getNodeRequest struct {
		ID int64 `json:"id"`
	}

	var req getNodeRequest
	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		JSONError(ctx, "Invalid JSON", err)
		return
	}

	if req.ID <= 0 {
		JSONError(ctx, "ID is required", nil)
		return
	}

	node, err := h.nodeRepository.GetByID(req.ID)
	if err != nil {
		JSONError(ctx, "Node not found", err)
		return
	}

	JSONSuccess(ctx, node)
}

func (h *BaseHandler) UpdateNodeHandler(ctx *fasthttp.RequestCtx) {
	var node model.Node
	if err := json.Unmarshal(ctx.PostBody(), &node); err != nil {
		JSONError(ctx, "Invalid JSON", err)
		return
	}

	if node.ID <= 0 {
		JSONError(ctx, "ID is required", nil)
		return
	}

	if err := h.nodeRepository.Update(&node); err != nil {
		JSONError(ctx, "Failed to update node", err)
		return
	}
	h.nodeManager.RemoveNode(node.ID)

	JSONSuccess(ctx, node)
}

func (h *BaseHandler) DeleteNodeHandler(ctx *fasthttp.RequestCtx) {
	type deleteNodeRequest struct {
		ID int64 `json:"id"`
	}

	var req deleteNodeRequest
	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		JSONError(ctx, "Invalid JSON", err)
		return
	}

	if req.ID <= 0 {
		JSONError(ctx, "ID is required", nil)
		return
	}
	h.nodeManager.RemoveNode(req.ID)
	if err := h.nodeRepository.Delete(req.ID); err != nil {
		JSONError(ctx, "Failed to delete node", err)
		return
	}

	JSONSuccess(ctx, nil)
}

func (h *BaseHandler) GetNodeListHandler(ctx *fasthttp.RequestCtx) {
	nodes, err := h.nodeRepository.List()
	if err != nil {
		JSONError(ctx, "Failed to get nodes", err)
		return
	}

	if nodes == nil {
		nodes = []model.Node{}
	}

	JSONSuccess(ctx, nodes)
}
