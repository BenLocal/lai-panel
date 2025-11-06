package route

import (
	"github.com/benlocal/lai-panel/pkg/handler"
	"github.com/benlocal/lai-panel/pkg/route"
	"github.com/fasthttp/router"
)

func init() {
	route.DefaultRegistry.Add(func(router *router.Router) {
		router.Handle("POST", "/node/add", handler.HandleAddNodeWithDI)
		router.Handle("POST", "/node/get", handler.HandleGetNodeWithDI)
		router.Handle("POST", "/node/update", handler.HandleUpdateNodeWithDI)
		router.Handle("POST", "/node/delete", handler.HandleDeleteNodeWithDI)
		router.Handle("POST", "/node/list", handler.HandleGetNodeListWithDI)
	})
}
