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
	dp, _ := docker.NewDockerProxy(dh, "/docker.proxy")
	g := gracefulshutdown.New()
	g.CatchSignals()
	appCtx := ctx.NewAppCtx(r.op, dp)
	baseHandler := handler.NewBaseHandler(appCtx)
	apiServer := api.NewApiServer(fmt.Sprintf(":%d", r.op.Port), baseHandler)
	g.Add(apiServer)

	registryService := service.NewRemoteRegistryService(
		r.op.Name,
		r.op.MasterHost,
		r.op.MasterPort,
		r.op.Address,
		r.op.Port,
	)
	g.Add(registryService)

	log.Println("start agent server on port", r.op.Port, "with name", r.op.Name)

	ctx := context.Background()
	return g.Start(ctx)
}
