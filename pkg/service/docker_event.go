package service

import (
	"context"
	"log"

	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"

	myClient "github.com/benlocal/lai-panel/pkg/client"
	"github.com/benlocal/lai-panel/pkg/ctx"
	"github.com/benlocal/lai-panel/pkg/model"
)

type dockerEventListenerService struct {
	dockerClient *client.Client
	baseClient   *myClient.BaseClient
}

func NewdockerEventListenerService(
	dockerClient *client.Client,
	baseClient *myClient.BaseClient,
) *dockerEventListenerService {
	return &dockerEventListenerService{
		dockerClient: dockerClient,
		baseClient:   baseClient,
	}
}

func (s *dockerEventListenerService) Name() string {
	return "docker-event-listener-service"
}

func (s *dockerEventListenerService) Start(ctx context.Context) error {
	eventChain, errChain := s.dockerClient.Events(ctx, events.ListOptions{})

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case event := <-eventChain:
			s.handleEvent(event)
		case err := <-errChain:
			log.Println("docker event listener service error: ", err)
		}
	}

	return nil
}

func (s *dockerEventListenerService) handleEvent(event events.Message) {
	id := ctx.GlobalServerStore.GetID()
	masterHost := ctx.GlobalServerStore.GetMasterHost()
	masterPort := ctx.GlobalServerStore.GetMasterPort()

	//log.Println("docker event listener service handleEvent: ", id, masterHost, masterPort)
	if id <= 0 {
		return
	}

	dockerEvent := model.DockerEvent{
		NodeId: id,
		Event:  event,
	}

	err := s.baseClient.DockerEvent(masterHost, masterPort, &dockerEvent)
	if err != nil {
		log.Println("docker event listener service error: ", err)
	}
}

func (s *dockerEventListenerService) Shutdown() error {
	return nil
}
