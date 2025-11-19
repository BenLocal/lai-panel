package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/benlocal/lai-panel/pkg/model"
	"github.com/cloudwego/hertz/pkg/app"
)

func (h *BaseHandler) GetServicePageHandler(ctx context.Context, c *app.RequestContext) {
	type getServicePageRequest struct {
		Page     int `json:"page"`
		PageSize int `json:"page_size"`
	}

	var req getServicePageRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.Error(err)
		return
	}

	total, services, err := h.ServiceRepository().GetPage(req.Page, req.PageSize)
	if err != nil {
		c.Error(err)
		return
	}

	type getServicePageResponse struct {
		Total    int                  `json:"total"`
		Page     int                  `json:"page"`
		PageSize int                  `json:"page_size"`
		Services []*model.ServiceView `json:"services"`
	}

	servicesView := make([]*model.ServiceView, 0)
	for _, service := range services {
		servicesView = append(servicesView, service.ToView())
	}

	c.JSON(http.StatusOK, SuccessResponse(getServicePageResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		Services: servicesView,
	}))
}

func (h *BaseHandler) SaveServiceHandler(ctx context.Context, c *app.RequestContext) {
	var req model.ServiceView
	if err := c.BindAndValidate(&req); err != nil {
		c.Error(err)
		return
	}

	if req.ID < 0 {
		c.Error(errors.New("ID is required"))
		return
	}

	service := req.ToModel()
	if req.ID == 0 {
		// add new service
		err := h.ServiceRepository().Create(service)
		if err != nil {
			c.Error(err)
			return
		}
	} else {
		// update existing service
		err := h.ServiceRepository().Update(service)
		if err != nil {
			c.Error(err)
			return
		}
	}

	c.JSON(http.StatusOK, SuccessResponse(nil))
}
