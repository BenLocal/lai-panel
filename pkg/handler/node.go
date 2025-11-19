package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/benlocal/lai-panel/pkg/model"
	"github.com/cloudwego/hertz/pkg/app"
)

func (h *BaseHandler) AddNodeHandler(ctx context.Context, c *app.RequestContext) {
	var node model.NodeView
	if err := c.BindAndValidate(&node); err != nil {
		c.Error(err)
		return
	}
	modelNode := node.ToModel()
	if err := h.NodeRepository().Create(modelNode); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(node))
}

func (h *BaseHandler) GetNodeHandler(ctx context.Context, c *app.RequestContext) {
	type getNodeRequest struct {
		ID int64 `json:"id"`
	}

	var req getNodeRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.Error(err)
		return
	}
	if req.ID <= 0 {
		c.Error(errors.New("ID is required"))
		return
	}

	node, err := h.NodeRepository().GetByID(req.ID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(node))
}

func (h *BaseHandler) UpdateNodeHandler(ctx context.Context, c *app.RequestContext) {
	var node model.NodeView
	if err := c.BindAndValidate(&node); err != nil {
		c.Error(err)
		return
	}
	if node.ID <= 0 {
		c.Error(errors.New("ID is required"))
		return
	}

	modelNode := node.ToModel()
	if err := h.NodeRepository().Update(modelNode); err != nil {
		c.Error(err)
		return
	}
	h.NodeManager().RemoveNode(node.ID)

	c.JSON(http.StatusOK, SuccessResponse(node))
}

func (h *BaseHandler) DeleteNodeHandler(ctx context.Context, c *app.RequestContext) {
	type deleteNodeRequest struct {
		ID int64 `json:"id"`
	}

	var req deleteNodeRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.Error(err)
		return
	}

	if req.ID <= 0 {
		c.Error(errors.New("ID is required"))
		return
	}
	h.NodeManager().RemoveNode(req.ID)
	if err := h.NodeRepository().Delete(req.ID); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, EmptyResponse())
}

func (h *BaseHandler) GetNodeListHandler(ctx context.Context, c *app.RequestContext) {
	nodes, err := h.NodeRepository().List()
	if err != nil {
		c.Error(err)
		return
	}

	var nodesView []*model.NodeView
	for _, node := range nodes {
		nodesView = append(nodesView, node.ToView())
	}

	c.JSON(http.StatusOK, SuccessResponse(nodesView))
}

func (h *BaseHandler) GetNodePageHandler(ctx context.Context, c *app.RequestContext) {
	type getNodePageRequest struct {
		Page     int `json:"page"`
		PageSize int `json:"page_size"`
	}

	type getNodePageResponse struct {
		Total    int               `json:"total"`
		Page     int               `json:"page"`
		PageSize int               `json:"page_size"`
		Nodes    []*model.NodeView `json:"nodes"`
	}

	var req getNodePageRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.Error(err)
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	total, nodes, err := h.NodeRepository().Page(req.Page, req.PageSize)
	if err != nil {
		c.Error(err)
		return
	}

	nodesView := make([]*model.NodeView, len(nodes))
	for i, node := range nodes {
		nodesView[i] = node.ToView()
	}

	resp := getNodePageResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Nodes:    nodesView,
	}

	c.JSON(http.StatusOK, SuccessResponse(resp))
}
