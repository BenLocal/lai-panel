package route

import (
	"github.com/benlocal/lai-panel/pkg/handler"
	"github.com/benlocal/lai-panel/pkg/route"
	"github.com/fasthttp/router"
)

func init() {
	route.DefaultRegistry.Add(func(h *handler.BaseHandler, router *router.Router) {
		router.Handle("GET", "/api/docker/info", h.DockerInfo)

		router.Handle("POST", "/api/docker/compose/config", h.HandleDockerComposeConfig)
		router.Handle("POST", "/api/docker/compose/deploy", h.HandleDockerComposeDeploy)
	})
}
