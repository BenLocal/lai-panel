package main

import (
	"fmt"
	"os"

	"github.com/benlocal/lai-panel/pkg/api"
	"github.com/benlocal/lai-panel/pkg/handler"
	"github.com/benlocal/lai-panel/pkg/options"
	"github.com/benlocal/lai-panel/pkg/version"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "lai-agent",
		Short: "lai-panel agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAgent(cmd)
		},
	}

	runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run lai-panel agent",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAgent(cmd)
		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print agent version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.Version)
		},
	}

	masterHost string
	masterPort int
	name       string
	address    string
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	api.DefaultRegistry.Add(func(h *handler.BaseHandler, router *route.Engine) {
		router.Handle("GET", "/healthz", h.HandleHealthz)
		router.Any("/docker.proxy/*path", h.HandleDockerProxy)

		api := router.Group("/open")
		// static files
		api.Static("/static", h.StaticDataPath())
		// file upload
		api.POST("/file/upload", h.HandleFileUpload)
	})
}

func init() {
	// CLI 命令和参数
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(versionCmd)

	runCmd.Flags().StringVar(&masterHost, "master-host", "localhost", "master host")
	runCmd.Flags().IntVar(&masterPort, "master-port", 8080, "master port")
	runCmd.Flags().StringVar(&name, "name", "", "name")
	runCmd.Flags().StringVar(&address, "address", "", "address")
}

func runAgent(_ *cobra.Command) error {
	op := options.NewAgentOptions(
		options.WithMasterHost(masterHost),
		options.WithMasterPort(masterPort),
		options.WithName(name),
		options.WithAddress(address),
	)

	runtime := NewAgentRuntime(op)

	if err := runtime.Start(); err != nil {
		return err
	}
	return nil
}
