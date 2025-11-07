package handler

import (
	"github.com/valyala/fasthttp"
)

func (b *BaseHandler) HandleDockerProxy(ctx *fasthttp.RequestCtx) {
	dd := b.dockerProxy
	if dd == nil {
		JSONError(ctx, "Docker proxy not initialized", nil)
		return
	}
	dd.HandleProxy(ctx)
}
