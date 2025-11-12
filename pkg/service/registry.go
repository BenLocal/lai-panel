package service

import (
	"context"
	"time"

	httpClient "github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/protocol"
)

type RegistryService struct {
	is_local    bool
	http_client *httpClient.Client

	context context.Context
	cancel  context.CancelFunc

	masterHost string
	masterPort int
}

func NewLocalRegistryService(masterPort int) *RegistryService {
	ctx, cancel := context.WithCancel(context.Background())
	return &RegistryService{
		context:    ctx,
		cancel:     cancel,
		is_local:   true,
		masterHost: "127.0.0.1",
		masterPort: masterPort,
	}
}

func NewRemoteRegistryService(masterHost string, masterPort int) *RegistryService {
	ctx, cancel := context.WithCancel(context.Background())
	client, _ := httpClient.NewClient()
	return &RegistryService{
		context:     ctx,
		cancel:      cancel,
		is_local:    false,
		http_client: client,
		masterHost:  masterHost,
		masterPort:  masterPort,
	}
}

func (s *RegistryService) Name() string {
	return "registry-service"
}

func (s *RegistryService) Start(ctx context.Context) error {
	if s.is_local {
		return s.updateLocalRegistry()
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-s.context.Done():
			return nil
		case <-ticker.C:
			s.updateRemoteRegistry()
		}
	}
}

func (s *RegistryService) Shutdown() error {
	s.cancel()
	return nil
}

func (s *RegistryService) updateRemoteRegistry() error {
	req := protocol.AcquireRequest()
	defer protocol.ReleaseRequest(req)

	req.SetRequestURI("http://127.0.0.1:8080/registry")
	req.Header.SetMethod("POST")
	req.Header.SetContentTypeBytes([]byte("application/json"))
	req.SetBodyString(`{"name":"test"}`)

	resp := protocol.AcquireResponse()
	defer protocol.ReleaseResponse(resp)

	if err := s.http_client.Do(context.Background(), req, resp); err != nil {
		return err
	}

	body := resp.Body()
	// 处理响应
	_ = body

	return nil
}

func (s *RegistryService) updateLocalRegistry() error {
	return nil
}
