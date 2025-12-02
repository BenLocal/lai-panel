package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/benlocal/lai-panel/pkg/constant"
	"github.com/benlocal/lai-panel/pkg/handler"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
)

type ServicesStateUpdater struct {
	context     context.Context
	cancel      context.CancelFunc
	baseHandler *handler.BaseHandler
}

func NewServicesStateService(baseHandler *handler.BaseHandler) *ServicesStateUpdater {
	ctx, cancel := context.WithCancel(context.Background())

	return &ServicesStateUpdater{
		context:     ctx,
		cancel:      cancel,
		baseHandler: baseHandler,
	}
}

func (s *ServicesStateUpdater) Name() string {
	return "services-state-updater"
}

func (s *ServicesStateUpdater) Start(ctx context.Context) error {
	ticker := time.NewTicker(25 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-s.context.Done():
			return nil
		case <-ticker.C:
			if err := s.updateServicesState(); err != nil {
				log.Println("update services state failed", err)
			}
		}
	}
}

func (s *ServicesStateUpdater) Shutdown() error {
	s.cancel()
	return nil
}

func (s *ServicesStateUpdater) updateServicesState() error {
	repo := s.baseHandler.ServiceRepository()
	services, err := repo.List()
	if err != nil {
		return err
	}
	for _, service := range services {
		state, err := s.baseHandler.NodeManager().GetNodeState(service.NodeID)
		if err != nil {
			continue
		}
		dc, err := state.GetDockerClient()
		if err != nil {
			continue
		}
		filters := filters.NewArgs()
		filters.Add("label", fmt.Sprintf("%s=%s", constant.OwnerLabel, constant.ProjectId))
		filters.Add("label", fmt.Sprintf("%s=%s", constant.ManagedByLabel, constant.ProjectId))
		filters.Add("label", fmt.Sprintf("%s=%s", constant.ServiceLabel, service.Name))

		containers, err := dc.ContainerList(context.Background(), container.ListOptions{
			All:     true,
			Filters: filters,
		})
		if err != nil {
			continue
		}

		status := "running"
		for _, container := range containers {
			if container.State != "running" {
				status = "stopped"
				break
			}
		}
		service.Status = status
		repo.UpdateStatusInfo(service)
	}

	return nil
}
