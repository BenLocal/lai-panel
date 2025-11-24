package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/benlocal/lai-panel/pkg/client"
	"github.com/benlocal/lai-panel/pkg/handler"
	"github.com/benlocal/lai-panel/pkg/model"
)

const (
	BaseIP = "127.0.0.1"
)

type RegistryService struct {
	context context.Context
	cancel  context.CancelFunc

	baseHandler *handler.BaseHandler
	baseClient  *client.BaseClient

	masterHost string
	masterPort int
	name       string
	agentPort  int
	is_local   bool
	address    string
	dataPath   *string
}

func NewLocalRegistryService(masterPort int, baseHandler *handler.BaseHandler, baseClient *client.BaseClient) *RegistryService {
	ctx, cancel := context.WithCancel(context.Background())
	return &RegistryService{
		context:     ctx,
		cancel:      cancel,
		baseClient:  baseClient,
		is_local:    true,
		masterPort:  masterPort,
		name:        "local",
		masterHost:  BaseIP,
		baseHandler: baseHandler,
		address:     BaseIP,
	}
}

func NewRemoteRegistryService(name string,
	dataPath *string,
	masterHost string,
	masterPort int,
	agentAddress string,
	agentPort int,
	baseClient *client.BaseClient) *RegistryService {
	ctx, cancel := context.WithCancel(context.Background())
	return &RegistryService{
		context:    ctx,
		cancel:     cancel,
		is_local:   false,
		masterHost: masterHost,
		masterPort: masterPort,
		name:       name,
		baseClient: baseClient,
		address:    agentAddress,
		agentPort:  agentPort,
		dataPath:   dataPath,
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
	reqBody := model.RegistryRequest{
		Name:      s.name,
		AgentPort: s.agentPort,
		IsLocal:   s.is_local,
		Status:    "online",
		Address:   s.address,
		DataPath:  s.dataPath,
	}
	resp, err := s.baseClient.Registry(s.masterHost, s.masterPort, &reqBody)
	if err != nil {
		return err
	}
	if resp.ID == 0 {
		return errors.New("registry failed")
	}
	return nil
}
