package server

import (
	"context"

	"github.com/benlocal/lai-panel/pkg/api"
	"github.com/benlocal/lai-panel/pkg/database"
	"github.com/benlocal/lai-panel/pkg/di"
	"github.com/benlocal/lai-panel/pkg/gracefulshutdown"
	"github.com/benlocal/lai-panel/pkg/handler"
	"github.com/benlocal/lai-panel/pkg/node"
	"github.com/benlocal/lai-panel/pkg/repository"
	"github.com/benlocal/lai-panel/pkg/route"
	"github.com/benlocal/lai-panel/pkg/service"
	"github.com/fasthttp/router"
)

type Runtime struct {
}

func NewRuntime() *Runtime {
	return &Runtime{}
}

func (r *Runtime) Start() error {
	op := NewOptions()
	err := database.InitDB(op.DBPath, op.MigrationsPath)
	if err != nil {
		return err
	}
	err = r.provideDependencies()
	if err != nil {
		return err
	}
	g := gracefulshutdown.New()
	g.CatchSignals()

	routeRouter := r.createApiRouter()
	apiServer := api.NewApiServer(":8080", routeRouter)
	g.Add(apiServer)

	registryService := service.NewLocalRegistryService()
	g.Add(registryService)

	ctx := context.Background()
	return g.Start(ctx)
}

func (r *Runtime) createApiRouter() *router.Router {

	router := router.New()
	for _, opt := range route.DefaultRegistry.Bindings() {
		opt(router)
	}

	return router
}

func (r *Runtime) provideDependencies() error {
	constructors := []interface{}{
		node.NewNodeManager,
		repository.NewNodeRepository,
		handler.NewHealthzHandler,
		handler.NewNodeApiHandler,
		handler.NewRegistryApiHandler,
		handler.NewDockerHandler,
	}
	for _, constructor := range constructors {
		err := di.Provide(constructor)
		if err != nil {
			return err
		}
	}
	return nil
}
