package handler

import (
	"context"

	"github.com/benlocal/lai-panel/pkg/docker"
	"github.com/benlocal/lai-panel/pkg/hub"
	"github.com/benlocal/lai-panel/pkg/node"
	"github.com/benlocal/lai-panel/pkg/repository"
)

type BaseHandler struct {
	// agent
	dockerProxy *docker.DockerProxy

	// server
	nodeManager       *node.NodeManager
	nodeRepository    *repository.NodeRepository
	appRepository     *repository.AppRepository
	serviceRepository *repository.ServiceRepository
	signalrServer     *hub.SignalRServer
}

func NewServerHandler() *BaseHandler {
	nodeRepository := repository.NewNodeRepository()
	nodeManager := node.NewNodeManager()
	appRepository := repository.NewAppRepository()
	serviceRepository := repository.NewServiceRepository()
	h := hub.NewSimpleHub(nodeRepository)
	signalrServer, _ := hub.NewSignalRServer(context.Background(), h)

	return &BaseHandler{
		nodeManager:       nodeManager,
		nodeRepository:    nodeRepository,
		appRepository:     appRepository,
		signalrServer:     signalrServer,
		serviceRepository: serviceRepository,
	}
}

func (h *BaseHandler) SignalRServer() *hub.SignalRServer {
	return h.signalrServer
}

func (h *BaseHandler) NodeRepository() *repository.NodeRepository {
	return h.nodeRepository
}

func NewAgentHandler(dp *docker.DockerProxy) *BaseHandler {
	return &BaseHandler{
		dockerProxy: dp,
	}
}

type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewApiResponse(code int, message string, data interface{}) *ApiResponse {
	return &ApiResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func SuccessResponse(data interface{}) *ApiResponse {
	return NewApiResponse(0, "success", data)
}

func ErrorResponse(code int, message string) *ApiResponse {
	return NewApiResponse(code, message, nil)
}

func FailResponse(message string) *ApiResponse {
	return NewApiResponse(-1, message, nil)
}

func EmptyResponse() *ApiResponse {
	return NewApiResponse(0, "success", nil)
}
