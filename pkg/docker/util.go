package docker

import (
	"fmt"

	"github.com/docker/docker/client"
)

func LocalDockerClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}

func AgentDockerClient(host string, port int) (*client.Client, error) {
	return client.NewClientWithOpts(client.WithHost(fmt.Sprintf("tcp://%s:%d/docker.proxy", host, port)),
		client.WithAPIVersionNegotiation())
}
