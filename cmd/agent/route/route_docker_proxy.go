package route

import (
	"github.com/benlocal/lai-panel/pkg/handler"
	"github.com/benlocal/lai-panel/pkg/route"
	"github.com/fasthttp/router"
)

func init() {
	route.DefaultRegistry.Add(func(router *router.Router) {
		router.ANY("/docker.proxy/{any:*}", handler.HandleDockerProxyWithDI)
	})
}
