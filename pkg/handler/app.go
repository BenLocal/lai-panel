package handler

import (
	"encoding/json"

	"github.com/benlocal/lai-panel/pkg/model"
	"github.com/valyala/fasthttp"
)

func (h *BaseHandler) GetApplicationPageHandler(ctx *fasthttp.RequestCtx) {
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
	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		JSONError(ctx, "Invalid JSON", err)
		return
	}

	total, apps, err := h.appRepository.ListPage(req.Page, req.PageSize)
	if err != nil {
		JSONError(ctx, "Failed to get application page", err)
		return
	}
	views := []*model.AppView{}
	for _, app := range apps {
		views = append(views, app.ToView())
	}
	JSONSuccess(ctx, getApplicationPageResponse{
		Total:       total,
		CurrentPage: req.Page,
		PageSize:    req.PageSize,
		Apps:        views,
	})
}

func (h *BaseHandler) GetApplicationListHandler(ctx *fasthttp.RequestCtx) {
	apps, err := h.appRepository.List()
	if err != nil {
		JSONError(ctx, "Failed to get application list", err)
		return
	}

	views := []*model.AppView{}
	for _, app := range apps {
		views = append(views, app.ToView())
	}

	JSONSuccess(ctx, views)
}

func (h *BaseHandler) AddApplicationHandler(ctx *fasthttp.RequestCtx) {
	var app model.AppView
	if err := json.Unmarshal(ctx.PostBody(), &app); err != nil {
		JSONError(ctx, "Invalid JSON", err)
		return
	}
	appModel := app.ToModel()
	if err := h.appRepository.Create(appModel); err != nil {
		JSONError(ctx, "Failed to create application", err)
		return
	}
	JSONSuccess(ctx, app)
}

func (h *BaseHandler) UpdateApplicationHandler(ctx *fasthttp.RequestCtx) {
	var app model.AppView
	if err := json.Unmarshal(ctx.PostBody(), &app); err != nil {
		JSONError(ctx, "Invalid JSON", err)
		return
	}
	appModel := app.ToModel()
	if err := h.appRepository.Update(appModel); err != nil {
		JSONError(ctx, "Failed to update application", err)
		return
	}
	JSONSuccess(ctx, app)
}

func (h *BaseHandler) DeleteApplicationHandler(ctx *fasthttp.RequestCtx) {
	type deleteApplicationRequest struct {
		ID int64 `json:"id"`
	}

	var req deleteApplicationRequest
	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		JSONError(ctx, "Invalid JSON", err)
		return
	}
	if err := h.appRepository.Delete(req.ID); err != nil {
		JSONError(ctx, "Failed to delete application", err)
		return
	}
	JSONEmptySuccess(ctx)
}

func (h *BaseHandler) GetApplicationHandler(ctx *fasthttp.RequestCtx) {
	type getApplicationRequest struct {
		ID int64 `json:"id"`
	}

	var req getApplicationRequest
	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		JSONError(ctx, "Invalid JSON", err)
		return
	}

	app, err := h.appRepository.GetByID(req.ID)
	if err != nil {
		JSONError(ctx, "Failed to get application", err)
		return
	}
	if app == nil {
		JSONError(ctx, "Application not found", nil)
		return
	}
	JSONSuccess(ctx, app)
}
