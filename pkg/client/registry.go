package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/benlocal/lai-panel/pkg/handler"
	"github.com/benlocal/lai-panel/pkg/model"
	"github.com/cloudwego/hertz/pkg/protocol"
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

	if resp.StatusCode() != 200 {
		respBody := resp.Body()
		return nil, fmt.Errorf("registry request failed, status code: %d, response: %s", resp.StatusCode(), string(respBody))
	}

	respBody := resp.Body()
	var respModel registryRespEnvelope
	if err := json.Unmarshal(respBody, &respModel); err != nil {
		return nil, fmt.Errorf("failed to parse registry response: %w, response body: %s", err, string(respBody))
	}

	if respModel.Code != 0 {
		return nil, errors.New(respModel.Message)
	}

	return &respModel.Data, nil
}

type registryRespEnvelope struct {
	handler.ApiResponse
	Data model.RegistryResponse `json:"data"`
}
