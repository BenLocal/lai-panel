package hub

import (
	"context"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/adaptor"
	"github.com/philippseith/signalr"
)

type SignalRServer struct {
	server signalr.Server
}

func NewSignalRServer(ctx context.Context, hub signalr.HubInterface, opts ...func(signalr.Party) error) (*SignalRServer, error) {
	baseOpts := []func(signalr.Party) error{
		signalr.HubFactory(func() signalr.HubInterface {
			return hub
		}),
		signalr.HTTPTransports(signalr.TransportWebSockets),
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

func (s *SignalRServer) Handler(path string) app.HandlerFunc {
	mux := http.NewServeMux()
	s.server.MapHTTP(signalr.WithHTTPServeMux(mux), path)

	return adaptor.HertzHandler(mux)
}

func (s *SignalRServer) HubClients() signalr.HubClients {
	if s == nil || s.server == nil {
		return nil
	}
	return s.server.HubClients()
}
