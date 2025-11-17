package main

import (
	"context"
	_ "embed"

	"github.com/benlocal/lai-panel/pkg/api"
	"github.com/benlocal/lai-panel/pkg/handler"
	"github.com/benlocal/lai-panel/pkg/serve"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/route"
)

func main() {
	runtime := serve.NewServeRuntime()

	if err := runtime.Start(); err != nil {
		panic(err)
	}
}

//go:embed signalr.html
var signalrHTML []byte

func init() {
	api.DefaultRegistry.Add(func(h *handler.BaseHandler, router *route.Engine) {
		router.Handle("GET", "/healthz", h.HandleHealthz)
		router.POST("/registry", h.GetRegistryHandler)

		router.Handle("POST", "/api/application/list", h.GetApplicationListHandler)
		router.Handle("POST", "/api/application/add", h.AddApplicationHandler)
		router.Handle("POST", "/api/application/update", h.UpdateApplicationHandler)
		router.Handle("POST", "/api/application/delete", h.DeleteApplicationHandler)
		router.Handle("POST", "/api/application/get", h.GetApplicationHandler)
		router.Handle("POST", "/api/application/page", h.GetApplicationPageHandler)

		router.Handle("POST", "/api/docker/info", h.DockerInfo)
		router.Handle("POST", "/api/docker/containers", h.DockerContainers)
		router.Handle("POST", "/api/docker/images", h.DockerImages)
		router.Handle("POST", "/api/docker/volumes", h.DockerVolumes)
		router.Handle("POST", "/api/docker/networks", h.DockerNetworks)

		router.Handle("POST", "/api/docker/compose/config", h.HandleDockerComposeConfig)
		router.Handle("POST", "/api/docker/compose/deploy", h.HandleDockerComposeDeploy)
		router.Handle("POST", "/api/node/add", h.AddNodeHandler)
		router.Handle("POST", "/api/node/get", h.GetNodeHandler)
		router.Handle("POST", "/api/node/update", h.UpdateNodeHandler)
		router.Handle("POST", "/api/node/delete", h.DeleteNodeHandler)
		router.Handle("POST", "/api/node/list", h.GetNodeListHandler)
		router.Handle("POST", "/api/node/page", h.GetNodePageHandler)

		handler := h.SignalRServer().Handler("/api/signalr")
		router.Any("/api/signalr", handler)
		router.Any("/api/signalr/*wsPath", handler)

		router.Handle("GET", "/signalr.html", func(ctx context.Context, c *app.RequestContext) {
			c.Response.Header.Set("Content-Type", "text/html; charset=utf-8")
			c.Response.SetBodyString(string(signalrHTML))
		})
	})
}
