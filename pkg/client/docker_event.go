package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/benlocal/lai-panel/pkg/model"
	"github.com/cloudwego/hertz/pkg/protocol"
)

func (c *BaseClient) DockerEvent(host string, port int, body *model.DockerEvent) error {
	req := protocol.AcquireRequest()
	defer protocol.ReleaseRequest(req)

	url := fmt.Sprintf("http://%s:%d%s", host, port, DockerEventPath)
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.SetContentTypeBytes([]byte("application/json"))
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req.SetBody(jsonBody)

	resp := protocol.AcquireResponse()
	defer protocol.ReleaseResponse(resp)

	if err := c.httpClient.Do(context.Background(), req, resp); err != nil {
		return err
	}

	return nil
}
