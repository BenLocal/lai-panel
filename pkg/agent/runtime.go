package agent

import (
	"context"
	"fmt"

	"github.com/benlocal/lai-panel/pkg/api"
	"github.com/benlocal/lai-panel/pkg/docker"
	"github.com/benlocal/lai-panel/pkg/gracefulshutdown"
	"github.com/benlocal/lai-panel/pkg/handler"
	"github.com/benlocal/lai-panel/pkg/service"
)

type AgentRuntime struct {
}

func NewAgentRuntime() *AgentRuntime {
	return &AgentRuntime{}
}

func (r *AgentRuntime) Start() error {
	op := NewAgentOptions()
	dp, _ := docker.NewDockerProxy("/var/run/docker.sock", "/docker.proxy")

	g := gracefulshutdown.New()
	g.CatchSignals()

	baseHandler := handler.NewAgentHandler(dp)
	apiServer := api.NewApiServer(fmt.Sprintf(":%d", op.Port), baseHandler)
	g.Add(apiServer)

	registryService := service.NewRemoteRegistryService(op.MasterHost, op.MasterPort)
	g.Add(registryService)

	ctx := context.Background()
	return g.Start(ctx)
}
