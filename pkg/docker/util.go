package docker

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/docker/docker/client"
)

func LocalDockerClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.WithAPIVersionNegotiation())
}

func AgentDockerClient(host string, port int) (*client.Client, error) {
	return agentDockerClient(host, port, false)
}

func agentDockerClient(host string, port int, withoutProxy bool) (*client.Client, error) {
	hostURL := fmt.Sprintf("tcp://%s:%d/docker.proxy", host, port)

	opts := []client.Opt{
		client.WithHost(hostURL),
		client.WithAPIVersionNegotiation(),
	}

	if withoutProxy {
		cc, err := customWithoutProxyHTTPClient()
		if err != nil {
			return nil, err
		}

		opts = append(opts, client.WithHTTPClient(cc))
	}
	return client.NewClientWithOpts(opts...)
}

func customWithoutProxyHTTPClient() (*http.Client, error) {
	transport := &http.Transport{}
	transport.MaxIdleConns = 6
	transport.IdleConnTimeout = 30 * time.Second
	transport.DisableCompression = false
	transport.DialContext = (&net.Dialer{
		Timeout: 10 * time.Second,
	}).DialContext
	return &http.Client{
		Transport: transport,
		CheckRedirect: func(_ *http.Request, via []*http.Request) error {
			if via[0].Method == http.MethodGet {
				return http.ErrUseLastResponse
			}
			return errors.New("unexpected redirect in response")
		},
	}, nil
}
