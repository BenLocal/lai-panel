package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/benlocal/lai-panel/pkg/model"
	"github.com/cloudwego/hertz/pkg/app"
)

func (h *BaseHandler) GetRegistryHandler(ctx context.Context, c *app.RequestContext) {
	var req model.RegistryRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.Error(err)
		return
	}

	registry, err := h.NodeRepository().GetByNodeName(req.Name)
	if err != nil {
		c.Error(err)
		return
	}

	if registry == nil {
		c.Error(errors.New("node not found"))
		return
	}

	node := &model.Node{
		ID:     registry.ID,
		Name:   registry.Name,
		Status: req.Status,
	}
	err = h.NodeRepository().UpdateRegistry(node)
	if err != nil {
		c.Error(err)
		return
	}

	resp := model.RegistryResponse{
		ID:   node.ID,
		Name: node.Name,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}
