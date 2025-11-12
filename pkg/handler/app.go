package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/benlocal/lai-panel/pkg/model"
	"github.com/cloudwego/hertz/pkg/app"
)

func (h *BaseHandler) GetApplicationPageHandler(ctx context.Context, c *app.RequestContext) {
	type getApplicationPageRequest struct {
		Page     int `json:"page"`
		PageSize int `json:"page_size"`
	}

	type getApplicationPageResponse struct {
		Total       int              `json:"total"`
		CurrentPage int              `json:"current_page"`
		PageSize    int              `json:"page_size"`
		Apps        []*model.AppView `json:"apps"`
	}

	var req getApplicationPageRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.Error(err)
		return
	}

	total, apps, err := h.appRepository.ListPage(req.Page, req.PageSize)
	if err != nil {
		c.Error(err)
		return
	}
	views := []*model.AppView{}
	for _, app := range apps {
		views = append(views, app.ToView())
	}

	c.JSON(http.StatusOK, getApplicationPageResponse{
		Total:       total,
		CurrentPage: req.Page,
		PageSize:    req.PageSize,
		Apps:        views,
	})
}

func (h *BaseHandler) GetApplicationListHandler(ctx context.Context, c *app.RequestContext) {
	apps, err := h.appRepository.List()
	if err != nil {
		c.Error(err)
		return
	}

	views := []*model.AppView{}
	for _, app := range apps {
		views = append(views, app.ToView())
	}

	c.JSON(http.StatusOK, views)
}

func (h *BaseHandler) AddApplicationHandler(ctx context.Context, c *app.RequestContext) {
	var app model.AppView
	if err := c.BindAndValidate(&app); err != nil {
		c.Error(err)
		return
	}
	appModel := app.ToModel()
	if err := h.appRepository.Create(appModel); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, app)
}

func (h *BaseHandler) UpdateApplicationHandler(ctx context.Context, c *app.RequestContext) {
	var app model.AppView
	if err := c.BindAndValidate(&app); err != nil {
		c.Error(err)
		return
	}
	appModel := app.ToModel()
	if err := h.appRepository.Update(appModel); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, app)
}

func (h *BaseHandler) DeleteApplicationHandler(ctx context.Context, c *app.RequestContext) {
	type deleteApplicationRequest struct {
		ID int64 `json:"id"`
	}

	var req deleteApplicationRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.Error(err)
		return
	}
	if err := h.appRepository.Delete(req.ID); err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, nil)
}

func (h *BaseHandler) GetApplicationHandler(ctx context.Context, c *app.RequestContext) {
	type getApplicationRequest struct {
		ID int64 `json:"id"`
	}

	var req getApplicationRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.Error(err)
		return
	}

	app, err := h.appRepository.GetByID(req.ID)
	if err != nil {
		c.Error(err)
		return
	}
	if app == nil {
		c.Error(errors.New("application not found"))
		return
	}
	c.JSON(http.StatusOK, app)
}
