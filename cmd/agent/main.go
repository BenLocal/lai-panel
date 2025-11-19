package main

import (
	"github.com/benlocal/lai-panel/pkg/agent"
	"github.com/benlocal/lai-panel/pkg/api"
	"github.com/benlocal/lai-panel/pkg/handler"
	"github.com/cloudwego/hertz/pkg/route"
)

func main() {
	runtime := agent.NewAgentRuntime()

	if err := runtime.Start(); err != nil {
		panic(err)
	}
}

func init() {
	api.DefaultRegistry.Add(func(h *handler.BaseHandler, router *route.Engine) {
		router.Handle("GET", "/healthz", h.HandleHealthz)
		router.Any("/docker.proxy/*path", h.HandleDockerProxy)

		api := router.Group("/open")
		// static file
		api.Static("/static", h.StaticDataPath())
		// file upload
		api.POST("/file/upload", h.HandleFileUpload)
	})
}
