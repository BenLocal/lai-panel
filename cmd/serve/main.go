package main

import (
	"context"
	_ "embed"

	"github.com/benlocal/lai-panel/pkg/api"
	"github.com/benlocal/lai-panel/pkg/client"
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
		router.POST(client.RegistryPath, h.GetRegistryHandler)
		router.POST(client.DockerEventPath, h.GetDockerEventHandler)

		openApi := router.Group("/open")
		// workspace files
		openApi.Static("/workspace", h.WorkSpaceDataPath())
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
		api.POST("/docker/container/start", h.DockerContainerStart)
		api.POST("/docker/container/stop", h.DockerContainerStop)
		api.POST("/docker/container/restart", h.DockerContainerRestart)
		api.POST("/docker/container/remove", h.DockerContainerRemove)
		api.POST("/docker/container/log", h.DockerContainerLog)
		api.POST("/docker/container/inspect", h.DockerContainerInspect)
		api.POST("/docker/images", h.DockerImages)
		api.POST("/docker/image/inspect", h.DockerImageInspect)
		api.POST("/docker/image/pushTo", h.DockerImagePushTo)
		api.POST("/docker/volumes", h.DockerVolumes)
		api.POST("/docker/networks", h.DockerNetworks)
		api.POST("/docker/compose/config", h.HandleDockerComposeConfig)
		api.POST("/docker/compose/deploy", h.HandleDockerComposeDeploy)
		api.POST("/docker/compose/undeploy", h.HandleDockerComposeUndeploy)
		api.POST("/node/add", h.AddNodeHandler)
		api.POST("/node/get", h.GetNodeHandler)
		api.POST("/node/update", h.UpdateNodeHandler)
		api.POST("/node/delete", h.DeleteNodeHandler)
		api.POST("/node/list", h.GetNodeListHandler)
		api.POST("/node/page", h.GetNodePageHandler)
		api.POST("/service/page", h.GetServicePageHandler)
		api.POST("/service/save", h.SaveServiceHandler)
		api.POST("/service/delete", h.DeleteServiceHandler)
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
