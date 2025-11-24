package main

import (
	"context"
	"fmt"
	"log"

	"github.com/benlocal/lai-panel/pkg/api"
	"github.com/benlocal/lai-panel/pkg/ctx"
	"github.com/benlocal/lai-panel/pkg/docker"
	"github.com/benlocal/lai-panel/pkg/gracefulshutdown"
	"github.com/benlocal/lai-panel/pkg/handler"
	"github.com/benlocal/lai-panel/pkg/options"
	"github.com/benlocal/lai-panel/pkg/service"
	"github.com/docker/docker/client"

	myClient "github.com/benlocal/lai-panel/pkg/client"
)

type AgentRuntime struct {
	op *options.AgentOptions
}

func NewAgentRuntime(options *options.AgentOptions) *AgentRuntime {
	return &AgentRuntime{
		op: options,
	}
}

func (r *AgentRuntime) Start() error {
	err := options.InitOptions(r.op)
	if err != nil {
		return err
	}

	dh := client.DefaultDockerHost
	localDockerClient, err := docker.LocalDockerClient()
	if err != nil {
		return err
	}
	dp, _ := docker.NewDockerProxy(dh, "/docker.proxy")
	baseClient := myClient.NewBaseClient()

	appCtx, err := ctx.NewAppCtx(r.op, dp)
	if err != nil {
		return err
	}

	g := gracefulshutdown.New()
	g.CatchSignals()

	baseHandler := handler.NewBaseHandler(appCtx)
	apiServer := api.NewApiServer(fmt.Sprintf(":%d", r.op.Port), baseHandler)
	g.Add(apiServer)

	dockerEventListenerService := service.NewdockerEventListenerService(
		localDockerClient,
		baseClient,
	)
	g.Add(dockerEventListenerService)

	registryService := service.NewRemoteRegistryService(baseClient)
	g.Add(registryService)

	log.Println("start agent server on port", r.op.Port, "with name", r.op.Name)

	ctx := context.Background()
	return g.Start(ctx)
}
