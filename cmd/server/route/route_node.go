package route

import (
	"github.com/benlocal/lai-panel/pkg/handler"
	"github.com/benlocal/lai-panel/pkg/route"
	"github.com/fasthttp/router"
)

func init() {
	route.DefaultRegistry.Add(func(h *handler.BaseHandler, router *router.Router) {
		router.Handle("POST", "/api/node/add", h.AddNodeHandler)
		router.Handle("POST", "/api/node/get", h.GetNodeHandler)
		router.Handle("POST", "/api/node/update", h.UpdateNodeHandler)
		router.Handle("POST", "/api/node/delete", h.DeleteNodeHandler)
		router.Handle("POST", "/api/node/list", h.GetNodeListHandler)
	})
}
