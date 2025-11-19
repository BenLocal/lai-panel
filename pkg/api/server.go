package api

import (
	"context"
	"log"

	"github.com/benlocal/lai-panel/pkg/handler"
	hertzServer "github.com/cloudwego/hertz/pkg/app/server"
)

type ApiServer struct {
	listenAddr  string
	server      *hertzServer.Hertz
	baseHandler *handler.BaseHandler
}

func NewApiServer(listenAddr string, baseHandler *handler.BaseHandler) *ApiServer {
	s := &ApiServer{
		listenAddr:  listenAddr,
		baseHandler: baseHandler,
	}
	return s
}

func (h *ApiServer) Name() string {
	return "api-http-server"
}

func (h *ApiServer) registryRouter() {
	for _, opt := range DefaultRegistry.Bindings() {
		opt(h.baseHandler, h.server.Engine)
	}
}

func (h *ApiServer) Start(ctx context.Context) error {
	log.Printf("Starting API server on %s", h.listenAddr)
	h.server = hertzServer.Default(
		hertzServer.WithHostPorts(h.listenAddr),
		hertzServer.WithMaxRequestBodySize(1*1024*1024*1024), // 1GB
	)
	h.registryRouter()
	h.server.Spin()
	return nil
}

func (h *ApiServer) Shutdown() error {
	if h.server != nil {
		h.server.Shutdown(context.Background())
	}
	return nil
}
