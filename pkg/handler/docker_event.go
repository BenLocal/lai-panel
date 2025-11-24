package handler

import (
	"context"
	"net/http"

	"github.com/benlocal/lai-panel/pkg/model"
	"github.com/cloudwego/hertz/pkg/app"
)

func (h *BaseHandler) GetDockerEventHandler(ctx context.Context, c *app.RequestContext) {
	var req model.DockerEvent
	if err := c.BindAndValidate(&req); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, SuccessResponse(nil))
}
