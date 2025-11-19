package main

import (
	"context"
	_ "embed"

	"github.com/benlocal/lai-panel/pkg/api"
	"github.com/benlocal/lai-panel/pkg/handler"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/route"
)

func main() {
	runtime := NewServeRuntime()

	if err := runtime.Start(); err != nil {
		panic(err)
	}
}

//go:embed signalr.html
var signalrHTML []byte

func init() {
	api.DefaultRegistry.Add(func(h *handler.BaseHandler, router *route.Engine) {
		router.GET("/healthz", h.HandleHealthz)
		router.POST("/registry", h.GetRegistryHandler)

		openApi := router.Group("/open")
		// static file
		openApi.Static("/static", h.StaticDataPath())
		// file upload
		openApi.POST("/file/upload", h.HandleFileUpload)

		api := router.Group("/api")
		api.POST("/application/list", h.GetApplicationListHandler)
		api.POST("/application/add", h.AddApplicationHandler)
		api.POST("/application/update", h.UpdateApplicationHandler)
		api.POST("/application/delete", h.DeleteApplicationHandler)
		api.POST("/application/get", h.GetApplicationHandler)
		api.POST("/application/page", h.GetApplicationPageHandler)
		api.POST("/docker/info", h.DockerInfo)
		api.POST("/docker/containers", h.DockerContainers)
		api.POST("/docker/images", h.DockerImages)
		api.POST("/docker/volumes", h.DockerVolumes)
		api.POST("/docker/networks", h.DockerNetworks)
		api.POST("/docker/compose/config", h.HandleDockerComposeConfig)
		api.POST("/docker/compose/deploy", h.HandleDockerComposeDeploy)
		api.POST("/node/add", h.AddNodeHandler)
		api.POST("/node/get", h.GetNodeHandler)
		api.POST("/node/update", h.UpdateNodeHandler)
		api.POST("/node/delete", h.DeleteNodeHandler)
		api.POST("/node/list", h.GetNodeListHandler)
		api.POST("/node/page", h.GetNodePageHandler)
		api.POST("/service/page", h.GetServicePageHandler)
		api.POST("/dashboard/stats", h.DashboardStatsHandler)

		// hub
		handler := h.SignalRServer().Handler("/api/signalr")
		router.Any("/api/signalr", handler)
		router.Any("/api/signalr/*wsPath", handler)
		router.Handle("GET", "/signalr.html", func(ctx context.Context, c *app.RequestContext) {
			c.Response.Header.Set("Content-Type", "text/html; charset=utf-8")
			c.Response.SetBodyString(string(signalrHTML))
		})
	})
}
