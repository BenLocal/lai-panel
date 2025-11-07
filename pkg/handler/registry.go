package handler

import (
	"encoding/json"

	"github.com/benlocal/lai-panel/pkg/model"
	"github.com/valyala/fasthttp"
)

func (h *BaseHandler) GetRegistryHandler(ctx *fasthttp.RequestCtx) {
	var req model.RegistryRequest
	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		JSONError(ctx, "Invalid JSON", err)
		return
	}

	registry, err := h.nodeRepository.GetByNodeName(req.Name)
	if err != nil {
		JSONError(ctx, "Failed to get registry", err)
		return
	}

	if registry == nil {
		JSONError(ctx, "Registry not found", nil)
		return
	}

	node := &model.Node{
		ID:     registry.ID,
		Name:   registry.Name,
		Status: req.Status,
	}
	err = h.nodeRepository.UpdateRegistry(node)
	if err != nil {
		JSONError(ctx, "Failed to update registry", err)
		return
	}

	resp := model.RegistryResponse{
		ID:   node.ID,
		Name: node.Name,
	}

	JSONSuccess(ctx, resp)
}
