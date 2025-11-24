package main

import (
	"context"
	"fmt"

	"github.com/benlocal/lai-panel/pkg/api"
	"github.com/benlocal/lai-panel/pkg/ctx"
	"github.com/benlocal/lai-panel/pkg/database"
	"github.com/benlocal/lai-panel/pkg/docker"
	"github.com/benlocal/lai-panel/pkg/gracefulshutdown"
	"github.com/benlocal/lai-panel/pkg/handler"
	"github.com/benlocal/lai-panel/pkg/options"
	"github.com/benlocal/lai-panel/pkg/service"

	myClient "github.com/benlocal/lai-panel/pkg/client"
)

type ServeRuntime struct {
}

func NewServeRuntime() *ServeRuntime {
	return &ServeRuntime{}
}

func (r *ServeRuntime) Start() error {
	op := options.NewServeOptions()
	err := database.InitDB(op.DBPath, op.MigrationsPath)
	if err != nil {
		return err
	}

	localDockerClient, err := docker.LocalDockerClient()
	if err != nil {
		return err
	}

	err = options.InitOptions(op)
	if err != nil {
		return err
	}

	appCtx, err := ctx.NewAppCtx(op, nil)
	if err != nil {
		return err
	}

	g := gracefulshutdown.New()
	g.CatchSignals()

	baseHandler := handler.NewBaseHandler(appCtx)
	baseClient := myClient.NewBaseClient()

	apiServer := api.NewApiServer(fmt.Sprintf(":%d", op.Port), baseHandler)
	g.Add(apiServer)

	registryService := service.NewLocalRegistryService(baseHandler, baseClient)
	g.Add(registryService)

	dockerEventListenerService := service.NewdockerEventListenerService(localDockerClient, baseClient)
	g.Add(dockerEventListenerService)

	ctx := context.Background()
	return g.Start(ctx)
}
