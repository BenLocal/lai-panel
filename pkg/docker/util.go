package docker

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/docker/docker/client"
)

func LocalDockerClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}

func AgentDockerClient(host string, port int) (*client.Client, error) {
	hostURL := fmt.Sprintf("http://%s:%d", host, port)
	httpClient := &http.Client{
		Transport: &proxyTransport{
			base: &http.Transport{
				Proxy: nil,
			},
			baseURL:   hostURL,
			proxyPath: "/docker.proxy",
		},
	}
	return client.NewClientWithOpts(
		client.FromEnv,
		client.WithHost(hostURL),
		client.WithAPIVersionNegotiation(),
		client.WithHTTPClient(httpClient),
	)
}

type proxyTransport struct {
	base      http.RoundTripper
	baseURL   string
	proxyPath string
}

func (t *proxyTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	originalPath := req.URL.Path
	req.URL.Path = t.proxyPath + originalPath
	if req.URL.Scheme == "" {
		req.URL.Scheme = "http"
	}
	if req.URL.Host == "" {
		baseURL, err := url.Parse(t.baseURL)
		if err != nil {
			return nil, err
		}
		req.URL.Host = baseURL.Host
	}

	return t.base.RoundTrip(req)
}
