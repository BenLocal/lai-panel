package service

import (
	"context"
	"log"
	"time"

	"github.com/benlocal/lai-panel/pkg/client"
	"github.com/benlocal/lai-panel/pkg/handler"
)

type HealthCheckService struct {
	context     context.Context
	cancel      context.CancelFunc
	baseClient  *client.BaseClient
	baseHandler *handler.BaseHandler
}

func NewHealthCheckService(baseClient *client.BaseClient,
	baseHandler *handler.BaseHandler) *HealthCheckService {
	ctx, cancel := context.WithCancel(context.Background())
	return &HealthCheckService{
		context:     ctx,
		cancel:      cancel,
		baseClient:  baseClient,
		baseHandler: baseHandler,
	}
}

func (s *HealthCheckService) Name() string {
	return "healthcheck"
}

func (s *HealthCheckService) Start(ctx context.Context) error {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-s.context.Done():
			return nil
		case <-ticker.C:
			if err := s.healthCheck(); err != nil {
				log.Println("health check failed", err)
			}
		}
	}
}

func (s *HealthCheckService) Shutdown() error {
	s.cancel()
	return nil
}

func (s *HealthCheckService) healthCheck() error {
	nodes, err := s.baseHandler.NodeRepository().List()
	if err != nil {
		return err
	}
	for _, node := range nodes {
		if !node.IsLocal {
			if err := s.baseClient.HealthCheck(node.Address, node.AgentPort); err != nil {
				log.Println("health check failed", node.Address, node.AgentPort, err)
				// update node status to offline
				s.baseHandler.NodeRepository().UpdateNodeStatus(node.ID, "offline")
			} else {
				// update node status to online
				s.baseHandler.NodeRepository().UpdateNodeStatus(node.ID, "online")
			}
		}
	}
	return nil
}
