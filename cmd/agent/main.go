package main

import (
	"flag"

	"github.com/benlocal/lai-panel/pkg/api"
	"github.com/benlocal/lai-panel/pkg/handler"
	"github.com/benlocal/lai-panel/pkg/options"
	"github.com/cloudwego/hertz/pkg/route"
)

func main() {
	var (
		masterHost = flag.String("master-host", "localhost", "master host")
		masterPort = flag.Int("master-port", 8080, "master port")
		name       = flag.String("name", "", "name")
		address    = flag.String("address", "", "address")
	)

	flag.Parse()

	op := options.NewAgentOptions(
		options.WithMasterHost(*masterHost),
		options.WithMasterPort(*masterPort),
		options.WithName(*name),
		options.WithAddress(*address),
	)

	runtime := NewAgentRuntime(op)

	if err := runtime.Start(); err != nil {
		panic(err)
	}
}

func init() {
	api.DefaultRegistry.Add(func(h *handler.BaseHandler, router *route.Engine) {
		router.Handle("GET", "/healthz", h.HandleHealthz)
		router.Any("/docker.proxy/*path", h.HandleDockerProxy)

		api := router.Group("/open")
		// workspace files
		api.Static("/workspace", h.WorkSpaceDataPath())
		// file upload
		api.POST("/file/upload", h.HandleFileUpload)
	})
}
