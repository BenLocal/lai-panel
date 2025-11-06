package service

import (
	"context"
	"time"

	"github.com/valyala/fasthttp"
)

type RegistryService struct {
	is_local    bool
	http_client *fasthttp.Client

	context context.Context
	cancel  context.CancelFunc
}

func NewLocalRegistryService() *RegistryService {
	ctx, cancel := context.WithCancel(context.Background())
	return &RegistryService{
		context:  ctx,
		cancel:   cancel,
		is_local: true,
	}
}

func NewRemoteRegistryService() *RegistryService {
	ctx, cancel := context.WithCancel(context.Background())
	client := &fasthttp.Client{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	return &RegistryService{
		context:     ctx,
		cancel:      cancel,
		is_local:    false,
		http_client: client,
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
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI("http://127.0.0.1:8080/registry")
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.SetBodyString(`{"name":"test"}`)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := s.http_client.Do(req, resp); err != nil {
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
