package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/benlocal/lai-panel/pkg/client"
	"github.com/benlocal/lai-panel/pkg/handler"
	"github.com/benlocal/lai-panel/pkg/model"

	appCtx "github.com/benlocal/lai-panel/pkg/ctx"
)

const (
	BaseIP = "127.0.0.1"
)

type RegistryService struct {
	context context.Context
	cancel  context.CancelFunc

	baseHandler *handler.BaseHandler
	baseClient  *client.BaseClient
}

func NewLocalRegistryService(baseHandler *handler.BaseHandler, baseClient *client.BaseClient) *RegistryService {
	ctx, cancel := context.WithCancel(context.Background())
	return &RegistryService{
		context:     ctx,
		cancel:      cancel,
		baseClient:  baseClient,
		baseHandler: baseHandler,
	}
}

func NewRemoteRegistryService(baseClient *client.BaseClient) *RegistryService {
	ctx, cancel := context.WithCancel(context.Background())
	return &RegistryService{
		context:    ctx,
		cancel:     cancel,
		baseClient: baseClient,
	}
}

func (s *RegistryService) Name() string {
	return "registry-service"
}

func (s *RegistryService) Start(ctx context.Context) error {
	isLocal := appCtx.GlobalServerStore.IsLocal()
	if isLocal {
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
	nodeName := appCtx.GlobalServerStore.GetName()
	node, err := s.baseHandler.NodeRepository().GetByNodeName(nodeName)
	if err != nil {
		return err
	}
	if node != nil {
		appCtx.GlobalServerStore.SetID(node.ID)
		return nil
	}

	// create local node
	node = &model.Node{
		Name:        nodeName,
		DisplayName: &nodeName,
		Address:     appCtx.GlobalServerStore.GetAddress(),
		SSHPort:     0,
		SSHUser:     "root",
		AgentPort:   0,
		IsLocal:     true,
		Status:      "online",
	}
	if err := s.baseHandler.NodeRepository().Create(node); err != nil {
		return err
	}
	appCtx.GlobalServerStore.SetID(node.ID)
	return nil
}

func (s *RegistryService) updateRegistry() error {
	nodeName := appCtx.GlobalServerStore.GetName()
	masterHost := appCtx.GlobalServerStore.GetMasterHost()
	masterPort := appCtx.GlobalServerStore.GetMasterPort()
	reqBody := model.RegistryRequest{
		Name:      nodeName,
		AgentPort: appCtx.GlobalServerStore.GetAgentPort(),
		IsLocal:   appCtx.GlobalServerStore.IsLocal(),
		Status:    "online",
		Address:   appCtx.GlobalServerStore.GetAddress(),
		DataPath:  appCtx.GlobalServerStore.GetDataPath(),
	}
	resp, err := s.baseClient.Registry(masterHost, masterPort, &reqBody)
	if err != nil {
		return err
	}
	if resp.ID <= 0 {
		return errors.New("registry failed")
	}

	// set service id
	appCtx.GlobalServerStore.SetID(resp.ID)
	return nil
}
