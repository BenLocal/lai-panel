package api

import (
	"context"
	"net/http"
	"time"

	"github.com/philippseith/signalr"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
)

type SignalRServer struct {
	server  signalr.Server
	handler fasthttp.RequestHandler
}

func NewSignalRServer(ctx context.Context, hub signalr.HubInterface, path string, opts ...func(signalr.Party) error) (*SignalRServer, error) {
	baseOpts := []func(signalr.Party) error{
		signalr.SimpleHubFactory(hub),
		signalr.HTTPTransports(signalr.TransportServerSentEvents),
		signalr.KeepAliveInterval(2 * time.Second),
	}
	opts = append(baseOpts, opts...)

	srv, err := signalr.NewServer(ctx, opts...)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	srv.MapHTTP(signalr.WithHTTPServeMux(mux), path)

	handler := fasthttpadaptor.NewFastHTTPHandler(mux)

	return &SignalRServer{
		server:  srv,
		handler: handler,
	}, nil
}

func (s *SignalRServer) Handle(ctx *fasthttp.RequestCtx) {
	if s == nil || s.handler == nil {
		ctx.Error("signalr handler not initialized", fasthttp.StatusInternalServerError)
		return
	}
	s.handler(ctx)
}

func (s *SignalRServer) HubClients() signalr.HubClients {
	if s == nil || s.server == nil {
		return nil
	}
	return s.server.HubClients()
}
