package server

import (
	"context"
	"fmt"

	"github.com/benlocal/lai-panel/pkg/api"
	"github.com/benlocal/lai-panel/pkg/database"
	"github.com/benlocal/lai-panel/pkg/gracefulshutdown"
	"github.com/benlocal/lai-panel/pkg/handler"
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

	g := gracefulshutdown.New()
	g.CatchSignals()

	baseHandler := handler.NewServerHandler()
	routeRouter := r.createApiRouter(baseHandler)

	listenAddr := fmt.Sprintf(":%d", op.Port)
	apiServer := api.NewApiServer(listenAddr, routeRouter)
	g.Add(apiServer)

	registryService := service.NewLocalRegistryService(op.Port)
	g.Add(registryService)

	ctx := context.Background()
	return g.Start(ctx)
}

func (r *Runtime) createApiRouter(baseHandler *handler.BaseHandler) *router.Router {

	router := router.New()
	for _, opt := range route.DefaultRegistry.Bindings() {
		opt(baseHandler, router)
	}

	return router
}
