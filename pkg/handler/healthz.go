package handler

import (
	"github.com/benlocal/lai-panel/pkg/di"
	"github.com/valyala/fasthttp"
)

type HealthzHandler struct {
}

func NewHealthzHandler() *HealthzHandler {
	return &HealthzHandler{}
}

func (h *HealthzHandler) Handle(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBodyString("UP")
}

func HandleHealthzWithDI(ctx *fasthttp.RequestCtx) {
	err := di.Invoke(func(h *HealthzHandler) {
		h.Handle(ctx)
	})
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(err.Error())
	}
}
