package agent

import (
	"context"

	"github.com/benlocal/lai-panel/pkg/api"
	"github.com/benlocal/lai-panel/pkg/di"
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
	dp, _ := handler.NewDockerProxy("/var/run/docker.sock", "/docker.proxy")

	err := r.provideDependencies(dp)
	if err != nil {
		return err
	}
	g := gracefulshutdown.New()
	g.CatchSignals()

	routeRouter := r.createApiRouter()
	apiServer := api.NewApiServer(":8081", routeRouter)
	g.Add(apiServer)

	registryService := service.NewRemoteRegistryService()
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

func (r *Runtime) provideDependencies(dockerProxy *handler.DockerProxy) error {
	constructors := []interface{}{
		handler.NewHealthzHandler,
		func() *handler.DockerProxy {
			return dockerProxy
		},
	}
	for _, constructor := range constructors {
		err := di.Provide(constructor)
		if err != nil {
			return err
		}
	}
	return nil
}
