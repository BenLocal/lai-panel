package service

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
)

type dockerEventListenerService struct {
	dockerClient *client.Client
}

func NewdockerEventListenerService(dockerClient *client.Client) *dockerEventListenerService {
	return &dockerEventListenerService{
		dockerClient: dockerClient,
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
			fmt.Println(event)
		case err := <-errChain:
			fmt.Println(err)
		}
	}

	return nil
}

func (s *dockerEventListenerService) Shutdown() error {
	return nil
}
