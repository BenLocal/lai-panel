package handler

import (
	"encoding/json"

	"github.com/benlocal/lai-panel/pkg/di"
	"github.com/benlocal/lai-panel/pkg/model"
	"github.com/benlocal/lai-panel/pkg/node"
	"github.com/benlocal/lai-panel/pkg/repository"
	"github.com/valyala/fasthttp"
)

type NodeApiHandler struct {
	nodeRepository *repository.NodeRepository
	manager        *node.NodeManager
}

func NewNodeApiHandler(nodeRepository *repository.NodeRepository,
	manager *node.NodeManager,
) *NodeApiHandler {
	return &NodeApiHandler{
		nodeRepository: nodeRepository,
		manager:        manager,
	}
}

func (h *NodeApiHandler) AddNodeHandler(ctx *fasthttp.RequestCtx) {
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

func (h *NodeApiHandler) GetNodeHandler(ctx *fasthttp.RequestCtx) {
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

func (h *NodeApiHandler) UpdateNodeHandler(ctx *fasthttp.RequestCtx) {
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
	h.manager.RemoveNode(node.ID)

	JSONSuccess(ctx, node)
}

func (h *NodeApiHandler) DeleteNodeHandler(ctx *fasthttp.RequestCtx) {
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
	h.manager.RemoveNode(req.ID)
	if err := h.nodeRepository.Delete(req.ID); err != nil {
		JSONError(ctx, "Failed to delete node", err)
		return
	}

	JSONSuccess(ctx, nil)
}

func (h *NodeApiHandler) GetNodeListHandler(ctx *fasthttp.RequestCtx) {
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

func HandleAddNodeWithDI(ctx *fasthttp.RequestCtx) {
	err := di.Invoke(func(h *NodeApiHandler) {
		h.AddNodeHandler(ctx)
	})
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(err.Error())
	}
}

func HandleGetNodeWithDI(ctx *fasthttp.RequestCtx) {
	err := di.Invoke(func(h *NodeApiHandler) {
		h.GetNodeHandler(ctx)
	})
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(err.Error())
	}
}

func HandleUpdateNodeWithDI(ctx *fasthttp.RequestCtx) {
	err := di.Invoke(func(h *NodeApiHandler) {
		h.UpdateNodeHandler(ctx)
	})
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(err.Error())
	}
}

func HandleDeleteNodeWithDI(ctx *fasthttp.RequestCtx) {
	err := di.Invoke(func(h *NodeApiHandler) {
		h.DeleteNodeHandler(ctx)
	})
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(err.Error())
	}
}

func HandleGetNodeListWithDI(ctx *fasthttp.RequestCtx) {
	err := di.Invoke(func(h *NodeApiHandler) {
		h.GetNodeListHandler(ctx)
	})
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(err.Error())
	}
}
