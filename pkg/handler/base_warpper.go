package handler

import (
	"github.com/benlocal/lai-panel/pkg/ctx"
	"github.com/benlocal/lai-panel/pkg/docker"
	"github.com/benlocal/lai-panel/pkg/hub"
	"github.com/benlocal/lai-panel/pkg/node"
	"github.com/benlocal/lai-panel/pkg/options"
	"github.com/benlocal/lai-panel/pkg/repository"
)

type BaseHandler struct {
	options options.IOptions
	appCtx  *ctx.AppCtx
}

func NewBaseHandler(appCtx *ctx.AppCtx) *BaseHandler {
	return &BaseHandler{
		appCtx:  appCtx,
		options: appCtx.Options(),
	}
}

func (h *BaseHandler) SignalRServer() *hub.SignalRServer {
	return h.appCtx.SignalRServer()
}

func (h *BaseHandler) NodeRepository() *repository.NodeRepository {
	return h.appCtx.NodeRepository()
}

func (h *BaseHandler) AppRepository() *repository.AppRepository {
	return h.appCtx.AppRepository()
}

func (h *BaseHandler) ServiceRepository() *repository.ServiceRepository {
	return h.appCtx.ServiceRepository()
}

func (h *BaseHandler) KvRepository() *repository.KvRepository {
	return h.appCtx.KvRepository()
}

func (h *BaseHandler) DockerProxy() *docker.DockerProxy {
	return h.appCtx.DockerProxy()
}

func (h *BaseHandler) NodeManager() *node.NodeManager {
	return h.appCtx.NodeManager()
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
