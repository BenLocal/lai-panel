package client

import httpClient "github.com/cloudwego/hertz/pkg/app/client"

type BaseClient struct {
	httpClient *httpClient.Client
}

func NewBaseClient() *BaseClient {
	httpClient, _ := httpClient.NewClient()

	return &BaseClient{
		httpClient: httpClient,
	}
}
