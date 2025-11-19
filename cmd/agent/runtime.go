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
)

type AgentRuntime struct {
}

func NewAgentRuntime() *AgentRuntime {
	return &AgentRuntime{}
}

func (r *AgentRuntime) Start() error {
	op := options.NewAgentOptions()
	err := options.InitOptions(op)
	if err != nil {
		return err
	}

	dh := client.DefaultDockerHost
	dp, _ := docker.NewDockerProxy(dh, "/docker.proxy")
	g := gracefulshutdown.New()
	g.CatchSignals()
	appCtx := ctx.NewAppCtx(op, dp)
	baseHandler := handler.NewBaseHandler(appCtx)
	apiServer := api.NewApiServer(fmt.Sprintf(":%d", op.Port), baseHandler)
	g.Add(apiServer)

	registryService := service.NewRemoteRegistryService(op.Name, op.MasterHost, op.MasterPort, op.Port)
	g.Add(registryService)

	log.Println("start agent server on port", op.Port, "with name", op.Name)

	ctx := context.Background()
	return g.Start(ctx)
}
