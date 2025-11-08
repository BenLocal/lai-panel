package handler

import (
	"encoding/json"
	"log"

	"github.com/benlocal/lai-panel/pkg/docker"
	"github.com/benlocal/lai-panel/pkg/node"
	"github.com/benlocal/lai-panel/pkg/repository"
	"github.com/valyala/fasthttp"
)

type BaseHandler struct {
	// agent
	dockerProxy *docker.DockerProxy

	// server
	nodeManager    *node.NodeManager
	nodeRepository *repository.NodeRepository
	appRepository  *repository.AppRepository
}

func NewServerHandler() *BaseHandler {
	nodeRepository := repository.NewNodeRepository()
	nodeManager := node.NewNodeManager()
	appRepository := repository.NewAppRepository()

	return &BaseHandler{
		nodeManager:    nodeManager,
		nodeRepository: nodeRepository,
		appRepository:  appRepository,
	}
}

func NewAgentHandler(dp *docker.DockerProxy) *BaseHandler {
	return &BaseHandler{
		dockerProxy: dp,
	}
}

type response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func JSONSuccess(ctx *fasthttp.RequestCtx, data interface{}) {
	resp := response{
		Code:    0,
		Message: "success",
		Data:    data,
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")
	json.NewEncoder(ctx).Encode(resp)
}

func JSONEmptySuccess(ctx *fasthttp.RequestCtx) {
	JSONSuccess(ctx, nil)
}

func JSONError(ctx *fasthttp.RequestCtx, message string, err error) {
	if err != nil {
		log.Printf("Failed with error: %v", err)
	}

	resp := response{
		Code:    -1,
		Message: message,
		Data:    nil,
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("application/json")

	json.NewEncoder(ctx).Encode(resp)
}
