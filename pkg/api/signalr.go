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
	server signalr.Server
}

func NewSignalRServer(ctx context.Context, hub signalr.HubInterface, opts ...func(signalr.Party) error) (*SignalRServer, error) {
	baseOpts := []func(signalr.Party) error{
		signalr.SimpleHubFactory(hub),
		signalr.HTTPTransports(signalr.TransportServerSentEvents),
		signalr.KeepAliveInterval(2 * time.Second),
		signalr.TimeoutInterval(6 * time.Second),
		signalr.HandshakeTimeout(15 * time.Second),
	}
	opts = append(baseOpts, opts...)

	srv, err := signalr.NewServer(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &SignalRServer{
		server: srv,
	}, nil
}

func (s *SignalRServer) Handler(path string) fasthttp.RequestHandler {
	mux := http.NewServeMux()
	s.server.MapHTTP(signalr.WithHTTPServeMux(mux), path)

	return fasthttpadaptor.NewFastHTTPHandler(mux)
}

func (s *SignalRServer) HubClients() signalr.HubClients {
	if s == nil || s.server == nil {
		return nil
	}
	return s.server.HubClients()
}
