package handler

import (
	"context"
	"errors"

	"github.com/cloudwego/hertz/pkg/app"
)

func (b *BaseHandler) HandleDockerProxy(ctx context.Context, c *app.RequestContext) {
	dd := b.dockerProxy
	if dd == nil {
		c.Error(errors.New("Docker proxy not initialized"))
		return
	}
	dd.HandleProxy(ctx, c)
}
