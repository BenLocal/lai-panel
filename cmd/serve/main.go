package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os"

	"github.com/benlocal/lai-panel/pkg/api"
	"github.com/benlocal/lai-panel/pkg/client"
	"github.com/benlocal/lai-panel/pkg/handler"
	"github.com/benlocal/lai-panel/pkg/version"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "lai-serve",
		Short: "lai-panel server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServe(cmd)
		},
	}

	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run lai-panel server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServe(cmd)
		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print server version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.Version)
		},
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

//go:embed signalr.html
var signalrHTML []byte

func init() {
	// HTTP 路由
	api.DefaultRegistry.Add(func(h *handler.BaseHandler, router *route.Engine) {
		router.GET("/healthz", h.HandleHealthz)
		router.POST(client.RegistryPath, h.GetRegistryHandler)
		router.POST(client.DockerEventPath, h.GetDockerEventHandler)

		openApi := router.Group("/open")
		// static files
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
		api.Static("/workspace", h.WorkSpaceDataPath())
		api.POST("/workspace/upload", h.HandleWorkspaceUpload)
		api.POST("/workspace/list", h.WorkspaceListHandler)
		api.POST("/workspace/read", h.WorkspaceReadHandler)
		api.POST("/workspace/save", h.WorkspaceSaveHandler)
		api.POST("/workspace/delete", h.WorkspaceDeleteHandler)
		api.POST("/workspace/mkdir", h.WorkspaceMkdirHandler)
		api.POST("/env/page", h.GetEnvPage)
		api.POST("/env/scopes", h.GetEnvScopes)
		api.POST("/env/addOrUpdate", h.AddOrUpdateEnv)
		api.POST("/env/delete", h.DeleteEnv)

		// hub
		sp := "/api/signalr"
		handler := h.SignalRServer().Handler(sp)
		router.Any(sp, handler)
		router.Any(sp+"/*wsPath", handler)
		// only for debug
		if version.Version == "dev" {
			router.Handle("GET", "/signalr.html", func(ctx context.Context, c *app.RequestContext) {
				c.Response.Header.Set("Content-Type", "text/html; charset=utf-8")
				c.Response.SetBodyString(string(signalrHTML))
			})
		}
	})

	// CLI 命令
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(versionCmd)
}

func runServe(_ *cobra.Command) error {
	log.Println("version:", version.Version)
	runtime := NewServeRuntime()

	if err := runtime.Start(); err != nil {
		return err
	}
	return nil
}
