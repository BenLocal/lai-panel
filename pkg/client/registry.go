package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/benlocal/lai-panel/pkg/model"
	"github.com/cloudwego/hertz/pkg/protocol"
)

const (
	RegistryPath = "/registry"
)

func (c *BaseClient) Registry(host string, port int, body *model.RegistryRequest) (*model.RegistryResponse, error) {
	req := protocol.AcquireRequest()
	defer protocol.ReleaseRequest(req)

	url := fmt.Sprintf("http://%s:%d%s", host, port, RegistryPath)
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.SetContentTypeBytes([]byte("application/json"))
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req.SetBody(jsonBody)

	resp := protocol.AcquireResponse()
	defer protocol.ReleaseResponse(resp)

	if err := c.httpClient.Do(context.Background(), req, resp); err != nil {
		return nil, err
	}

	respBody := resp.Body()
	var respModel model.RegistryResponse
	if err := json.Unmarshal(respBody, &respModel); err != nil {
		return nil, err
	}
	return &respModel, nil
}
