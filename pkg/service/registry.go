package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/benlocal/lai-panel/pkg/handler"
	"github.com/benlocal/lai-panel/pkg/model"
	httpClient "github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/protocol"
)

type RegistryService struct {
	context    context.Context
	cancel     context.CancelFunc
	httpClient *httpClient.Client

	baseHandler *handler.BaseHandler

	masterHost string
	masterPort int
	name       string
	agentPort  int
	is_local   bool
	address    string
}

func NewLocalRegistryService(masterPort int, baseHandler *handler.BaseHandler) *RegistryService {
	ctx, cancel := context.WithCancel(context.Background())
	httpClient, _ := httpClient.NewClient()
	return &RegistryService{
		context:     ctx,
		cancel:      cancel,
		httpClient:  httpClient,
		is_local:    true,
		masterPort:  masterPort,
		name:        "local",
		masterHost:  "127.0.0.1",
		baseHandler: baseHandler,
		address:     "127.0.0.1",
	}
}

func NewRemoteRegistryService(name string, masterHost string, masterPort int, agentAddress string, agentPort int) *RegistryService {
	ctx, cancel := context.WithCancel(context.Background())
	httpClient, _ := httpClient.NewClient()
	return &RegistryService{
		context:    ctx,
		cancel:     cancel,
		is_local:   false,
		masterHost: masterHost,
		masterPort: masterPort,
		name:       name,
		httpClient: httpClient,
		address:    agentAddress,
		agentPort:  agentPort,
	}
}

func (s *RegistryService) Name() string {
	return "registry-service"
}

func (s *RegistryService) Start(ctx context.Context) error {
	if s.is_local {
		if err := s.tryAddLocalRegistry(); err != nil {
			log.Println("try add local registry failed", err)
		}
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
			if err := s.updateRegistry(); err != nil {
				log.Println("update registry failed", err)
			}
		}
	}
}

func (s *RegistryService) Shutdown() error {
	s.cancel()
	return nil
}

func (s *RegistryService) tryAddLocalRegistry() error {
	node, err := s.baseHandler.NodeRepository().GetByNodeName(s.name)
	if err != nil {
		return err
	}
	if node != nil {
		return nil
	}

	// create local node
	node = &model.Node{
		Name:        s.name,
		DisplayName: &s.name,
		Address:     "127.0.0.1",
		SSHPort:     s.masterPort,
		SSHUser:     "root",
		AgentPort:   0,
		IsLocal:     true,
		Status:      "online",
	}
	if err := s.baseHandler.NodeRepository().Create(node); err != nil {
		return err
	}
	return nil
}

func (s *RegistryService) updateRegistry() error {
	req := protocol.AcquireRequest()
	defer protocol.ReleaseRequest(req)

	url := fmt.Sprintf("http://%s:%d/registry", s.masterHost, s.masterPort)
	log.Println("update registry to", url)
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.SetContentTypeBytes([]byte("application/json"))

	reqBody := model.RegistryRequest{
		Name:      s.name,
		AgentPort: s.agentPort,
		IsLocal:   s.is_local,
		Status:    "online",
		Address:   s.address,
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req.SetBody(jsonBody)

	resp := protocol.AcquireResponse()
	defer protocol.ReleaseResponse(resp)

	if err := s.httpClient.Do(context.Background(), req, resp); err != nil {
		return err
	}

	body := resp.Body()
	// 处理响应
	_ = body

	return nil
}
