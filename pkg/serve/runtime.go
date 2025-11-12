package serve

import (
	"context"
	"fmt"

	"github.com/benlocal/lai-panel/pkg/api"
	"github.com/benlocal/lai-panel/pkg/database"
	"github.com/benlocal/lai-panel/pkg/gracefulshutdown"
	"github.com/benlocal/lai-panel/pkg/handler"
	"github.com/benlocal/lai-panel/pkg/service"
)

type ServeRuntime struct {
}

func NewServeRuntime() *ServeRuntime {
	return &ServeRuntime{}
}

func (r *ServeRuntime) Start() error {
	op := NewServeOptions()
	err := database.InitDB(op.DBPath, op.MigrationsPath)
	if err != nil {
		return err
	}

	g := gracefulshutdown.New()
	g.CatchSignals()

	baseHandler := handler.NewServerHandler()

	apiServer := api.NewApiServer(fmt.Sprintf(":%d", op.Port), baseHandler)
	g.Add(apiServer)

	registryService := service.NewLocalRegistryService(op.Port)
	g.Add(registryService)

	ctx := context.Background()
	return g.Start(ctx)
}
