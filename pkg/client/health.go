package client

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/hertz/pkg/protocol"
)

const (
	HealthCheckPath = "/healthz"
)

func (c *BaseClient) HealthCheck(host string, port int) error {
	req := protocol.AcquireRequest()
	defer protocol.ReleaseRequest(req)

	url := fmt.Sprintf("http://%s:%d%s", host, port, HealthCheckPath)
	req.SetRequestURI(url)
	req.Header.SetMethod("GET")
	req.Header.SetContentTypeBytes([]byte("application/json"))

	resp := protocol.AcquireResponse()
	defer protocol.ReleaseResponse(resp)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := c.httpClient.Do(ctx, req, resp); err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return fmt.Errorf("health check failed, status code: %d", resp.StatusCode())
	}

	return nil
}
