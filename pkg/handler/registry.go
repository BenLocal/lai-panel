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

	if req.IsLocal {
		resp, err := h.local(&req)
		if err != nil {
			c.Error(err)
			return
		}
		c.JSON(http.StatusOK, SuccessResponse(resp))
		return
	}

	// add remote registry

	resp, err := h.remote(&req)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(resp))
}

func (h *BaseHandler) remote(req *model.RegistryRequest) (*model.RegistryResponse, error) {
	registry, err := h.NodeRepository().GetByNodeName(req.Name)
	if err != nil {
		return nil, err
	}

	if registry == nil {
		// create new node
		node := &model.Node{
			Name:      req.Name,
			Status:    req.Status,
			IsLocal:   req.IsLocal,
			AgentPort: req.AgentPort,
			Address:   req.Address,
		}
		err := h.NodeRepository().Create(node)
		if err != nil {
			return nil, err
		}
		return &model.RegistryResponse{
			ID:   node.ID,
			Name: node.Name,
		}, nil
	} else {
		// update node
		node := &model.Node{
			ID:        registry.ID,
			Name:      registry.Name,
			Status:    req.Status,
			Address:   req.Address,
			AgentPort: req.AgentPort,
		}
		if registry.Status != req.Status || registry.Address != req.Address || registry.AgentPort != req.AgentPort {
			err = h.NodeRepository().UpdateRegistry(node)
			if err != nil {
				return nil, err
			}
		}

		return &model.RegistryResponse{
			ID:   node.ID,
			Name: node.Name,
		}, nil
	}
}

func (h *BaseHandler) local(req *model.RegistryRequest) (*model.RegistryResponse, error) {
	registry, err := h.NodeRepository().GetByNodeName(req.Name)
	if err != nil {
		return nil, err
	}
	if registry == nil {
		return nil, errors.New("node not found")
	}

	if registry.Status != req.Status {
		err = h.NodeRepository().UpdateNodeStatus(registry.ID, req.Status)
		if err != nil {
			return nil, err
		}
	}

	return &model.RegistryResponse{
		ID:   registry.ID,
		Name: registry.Name,
	}, nil
}
