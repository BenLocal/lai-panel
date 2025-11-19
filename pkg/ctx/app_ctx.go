package ctx

import (
	"context"

	"github.com/benlocal/lai-panel/pkg/docker"
	"github.com/benlocal/lai-panel/pkg/hub"
	"github.com/benlocal/lai-panel/pkg/node"
	"github.com/benlocal/lai-panel/pkg/options"
	"github.com/benlocal/lai-panel/pkg/pipe"
	"github.com/benlocal/lai-panel/pkg/repository"
)

type AppCtx struct {
	options           options.IOptions
	dockerProxy       *docker.DockerProxy
	nodeManager       *node.NodeManager
	nodeRepository    *repository.NodeRepository
	appRepository     *repository.AppRepository
	serviceRepository *repository.ServiceRepository
	signalrServer     *hub.SignalRServer
	kvRepository      *repository.KvRepository
	nodePipeline      *pipe.NodePipeline
}

func NewAppCtx(options options.IOptions, dockerProxy *docker.DockerProxy) *AppCtx {
	// server
	if !options.Agent() {
		nodeRepository := repository.NewNodeRepository()
		nodeManager := node.NewNodeManager()
		appRepository := repository.NewAppRepository()
		serviceRepository := repository.NewServiceRepository()
		kvRepository := repository.NewKvRepository()
		h := hub.NewSimpleHub(nodeRepository)
		signalrServer, _ := hub.NewSignalRServer(context.Background(), h)
		nodePipeline := pipe.NewNodePipeline()

		return &AppCtx{
			kvRepository:      kvRepository,
			nodeManager:       nodeManager,
			nodeRepository:    nodeRepository,
			appRepository:     appRepository,
			signalrServer:     signalrServer,
			serviceRepository: serviceRepository,
			options:           options,
			nodePipeline:      nodePipeline,
			dockerProxy:       dockerProxy,
		}
	}

	return &AppCtx{
		options:     options,
		dockerProxy: dockerProxy,
	}
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

func (a *AppCtx) NodePipeline() *pipe.NodePipeline {
	return a.nodePipeline
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
