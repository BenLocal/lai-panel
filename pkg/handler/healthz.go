package handler

import (
	"context"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
)

func (h *BaseHandler) HandleHealthz(ctx context.Context, c *app.RequestContext) {
	c.JSON(http.StatusOK, map[string]string{"status": "UP"})
}
