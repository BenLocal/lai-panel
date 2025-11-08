package api

import (
	"context"
	"log"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type ApiServer struct {
	listenAddr string
	router     *router.Router

	server *fasthttp.Server
}

func NewApiServer(listenAddr string, router *router.Router) *ApiServer {
	return &ApiServer{
		listenAddr: listenAddr,
		router:     router,
	}
}

func (h *ApiServer) Name() string {
	return "api-http-server"
}

func (h *ApiServer) Start(ctx context.Context) error {
	corsHandler := func(ctx *fasthttp.RequestCtx) {
		// 设置 CORS 头
		ctx.Response.Header.Set("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "GET,POST")
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")

		// 处理预检请求
		if string(ctx.Method()) == fasthttp.MethodOptions {
			ctx.SetStatusCode(fasthttp.StatusNoContent)
			return
		}

		// 交给原路由处理
		h.router.Handler(ctx)
	}
	h.server = &fasthttp.Server{
		Handler: corsHandler,
	}

	log.Printf("Starting API server on %s", h.listenAddr)
	if err := h.server.ListenAndServe(h.listenAddr); err != nil {
		log.Fatalf("Error in ListenAndServe: %v", err)
		return err
	}

	return nil
}

func (h *ApiServer) Shutdown() error {
	if h.server != nil {
		if err := h.server.Shutdown(); err != nil {
			log.Printf("Error shutting down server: %v", err)
			return err
		}
	}
	return nil
}
