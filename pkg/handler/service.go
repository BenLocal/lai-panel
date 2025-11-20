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
	type saveServiceResponse struct {
		ID int64 `json:"id"`
	}

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
	var id int64
	if req.ID == 0 {
		// add new service
		i, err := h.ServiceRepository().Create(service)
		if err != nil {
			c.Error(err)
			return
		}
		id = i
	} else {
		// update existing service
		err := h.ServiceRepository().Update(service)
		if err != nil {
			c.Error(err)
			return
		}

		id = service.ID
	}

	c.JSON(http.StatusOK, SuccessResponse(saveServiceResponse{
		ID: id,
	}))
}

func (h *BaseHandler) DeleteServiceHandler(ctx context.Context, c *app.RequestContext) {
	type deleteServiceRequest struct {
		ID    int64 `json:"id"`
		Force bool  `json:"force"`
	}

	var req deleteServiceRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.Error(err)
		return
	}

	currentService, err := h.ServiceRepository().GetByID(req.ID)
	if err != nil {
		c.Error(err)
		return
	}
	if currentService == nil {
		c.Error(errors.New("service not found"))
		return
	}

	// check if service is deployed
	if currentService.DeployInfo != nil {
		if !req.Force {
			c.Error(errors.New("service is deployed, use force to undeploy"))
			return
		}

		_, err = h.dockerComposeUndeploy(ctx, currentService)
		if err != nil {
			c.Error(err)
			return
		}
	}

	err = h.ServiceRepository().Delete(req.ID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(nil))
}
