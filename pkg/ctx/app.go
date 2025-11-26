package ctx

import (
	"context"
	"errors"

	"github.com/benlocal/lai-panel/pkg/docker"
	"github.com/benlocal/lai-panel/pkg/hub"
	"github.com/benlocal/lai-panel/pkg/node"
	"github.com/benlocal/lai-panel/pkg/options"
	"github.com/benlocal/lai-panel/pkg/repository"
)

type AppCtx struct {
	options           options.IOptions
	dockerProxy       *docker.DockerProxy
	nodeManager       *node.NodeManager
	nodeRepository    *repository.NodeRepository
	appRepository     *repository.AppRepository
	serviceRepository *repository.ServiceRepository
	envRepository     *repository.EnvRepository
	signalrServer     *hub.SignalRServer
	kvRepository      *repository.KvRepository
	serverStore       *ServerStore
}

func NewAppCtx(opt options.IOptions, dockerProxy *docker.DockerProxy) (*AppCtx, error) {
	// server
	if !opt.Agent() {
		ss := GetServerStoreForLocal(opt.(*options.ServeOptions))
		if ss == nil {
			return nil, errors.New("failed to get server store for local")
		}

		nodeRepository := repository.NewNodeRepository()
		nodeManager := node.NewNodeManager(nodeRepository)
		appRepository := repository.NewAppRepository()
		serviceRepository := repository.NewServiceRepository()
		kvRepository := repository.NewKvRepository()
		envRepository := repository.NewEnvRepository()
		h := hub.NewSimpleHub(nodeRepository, nodeManager)
		signalrServer, _ := hub.NewSignalRServer(context.Background(), h)

		return &AppCtx{
			kvRepository:      kvRepository,
			nodeManager:       nodeManager,
			nodeRepository:    nodeRepository,
			appRepository:     appRepository,
			signalrServer:     signalrServer,
			serviceRepository: serviceRepository,
			options:           opt,
			dockerProxy:       dockerProxy,
			serverStore:       ss,
			envRepository:     envRepository,
		}, nil
	}

	ss := GetServerStore(opt.(*options.AgentOptions))
	if ss == nil {
		return nil, errors.New("failed to get server store for agent")
	}

	return &AppCtx{
		options:     opt,
		dockerProxy: dockerProxy,
		serverStore: ss,
	}, nil
}

func (a *AppCtx) SignalRServer() *hub.SignalRServer {
	return a.signalrServer
}

func (a *AppCtx) NodeRepository() *repository.NodeRepository {
	return a.nodeRepository
}

func (a *AppCtx) AppRepository() *repository.AppRepository {
	return a.appRepository
}

func (a *AppCtx) ServiceRepository() *repository.ServiceRepository {
	return a.serviceRepository
}

func (a *AppCtx) KvRepository() *repository.KvRepository {
	return a.kvRepository
}

func (a *AppCtx) DockerProxy() *docker.DockerProxy {
	return a.dockerProxy
}

func (a *AppCtx) Options() options.IOptions {
	return a.options
}

func (a *AppCtx) NodeManager() *node.NodeManager {
	return a.nodeManager
}

func (a *AppCtx) EnvRepository() *repository.EnvRepository {
	return a.envRepository
}
