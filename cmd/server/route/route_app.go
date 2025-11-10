package route

import (
	"github.com/benlocal/lai-panel/pkg/handler"
	"github.com/benlocal/lai-panel/pkg/route"
	"github.com/fasthttp/router"
)

func init() {
	route.DefaultRegistry.Add(func(h *handler.BaseHandler, router *router.Router) {
		router.Handle("POST", "/api/application/list", h.GetApplicationListHandler)
		router.Handle("POST", "/api/application/add", h.AddApplicationHandler)
		router.Handle("POST", "/api/application/update", h.UpdateApplicationHandler)
		router.Handle("POST", "/api/application/delete", h.DeleteApplicationHandler)
		router.Handle("POST", "/api/application/get", h.GetApplicationHandler)
		router.Handle("POST", "/api/application/page", h.GetApplicationPageHandler)
	})
}
