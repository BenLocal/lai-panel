package route

import (
	_ "embed"

	"github.com/benlocal/lai-panel/pkg/handler"
	"github.com/benlocal/lai-panel/pkg/route"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

//go:embed signalr.html
var signalrHTML []byte

func init() {
	route.DefaultRegistry.Add(func(h *handler.BaseHandler, router *router.Router) {
		handler := h.SignalRServer().Handler("/api/signalr")
		router.ANY("/api/signalr", handler)
		router.ANY("/api/signalr/{filepath:*}", handler)

		router.Handle("GET", "/signalr.html", func(ctx *fasthttp.RequestCtx) {
			ctx.Response.Header.Set("Content-Type", "text/html; charset=utf-8")
			ctx.SuccessString("text/html; charset=utf-8", string(signalrHTML))
		})
	})
}
