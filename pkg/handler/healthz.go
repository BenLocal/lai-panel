package handler

import (
	"github.com/valyala/fasthttp"
)

func (h *BaseHandler) HandleHealthz(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("UP")
}
