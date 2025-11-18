package handler

import (
	"context"
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

	total, services, err := h.serviceRepository.GetPage(req.Page, req.PageSize)
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
