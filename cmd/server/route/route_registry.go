package route

import (
	"github.com/benlocal/lai-panel/pkg/handler"
	"github.com/benlocal/lai-panel/pkg/route"
	"github.com/fasthttp/router"
)

func init() {
	route.DefaultRegistry.Add(func(h *handler.BaseHandler, router *router.Router) {
		router.POST("/registry", h.GetRegistryHandler)
	})
}
