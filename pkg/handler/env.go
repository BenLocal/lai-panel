package handler

import (
	"context"
	"net/http"

	"github.com/benlocal/lai-panel/pkg/model"
	"github.com/cloudwego/hertz/pkg/app"
)

func (b *BaseHandler) GetEnvPage(ctx context.Context, c *app.RequestContext) {
	type getEnvPageRequest struct {
		Scope    string `json:"scope"`
		Page     int    `json:"page"`
		PageSize int    `json:"page_size"`
	}
	type getEnvPageResponse struct {
		Total       int         `json:"total"`
		CurrentPage int         `json:"current_page"`
		PageSize    int         `json:"page_size"`
		List        []model.Env `json:"list"`
	}

	var req getEnvPageRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	total, lst, err := b.EnvRepository().GetPage(req.Scope, req.Page, req.PageSize)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(getEnvPageResponse{
		Total:       total,
		CurrentPage: req.Page,
		PageSize:    req.PageSize,
		List:        lst,
	}))
}

func (b *BaseHandler) GetEnvScopes(ctx context.Context, c *app.RequestContext) {
	scopes, err := b.EnvRepository().GetScopes()
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(scopes))
}

func (b *BaseHandler) DeleteEnv(ctx context.Context, c *app.RequestContext) {
	type deleteEnvRequest struct {
		ID int64 `json:"id"`
	}
	var req deleteEnvRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := b.EnvRepository().Delete(req.ID)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, SuccessResponse(nil))
}

func (b *BaseHandler) AddOrUpdateEnv(ctx context.Context, c *app.RequestContext) {
	type addOrUpdateEnvRequest struct {
		ID    int64  `json:"id"`
		Key   string `json:"key"`
		Value string `json:"value"`
		Scope string `json:"scope"`
	}

	var req addOrUpdateEnvRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if req.ID == 0 {
		m := &model.Env{
			Key:   req.Key,
			Value: req.Value,
			Scope: req.Scope,
		}
		err := b.EnvRepository().Create(m)
		if err != nil {
			c.Error(err)
			return
		}
	} else {
		m := &model.Env{
			ID:    req.ID,
			Key:   req.Key,
			Value: req.Value,
			Scope: req.Scope,
		}
		err := b.EnvRepository().Update(m)
		if err != nil {
			c.Error(err)
			return
		}
	}
	c.JSON(http.StatusOK, SuccessResponse(nil))
}
