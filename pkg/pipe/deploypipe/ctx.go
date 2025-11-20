package deploypipe

import (
	"sync"

	"github.com/benlocal/lai-panel/pkg/model"
	"github.com/benlocal/lai-panel/pkg/node"
	"github.com/benlocal/lai-panel/pkg/options"
	"github.com/cloudwego/hertz/pkg/protocol/sse"
)

type DeployCtx struct {
	options   options.IOptions
	App       *model.App
	Service   *model.Service
	NodeState *node.NodeState
	writer    *sse.Writer
	sendMu    sync.Mutex
	env       map[string]string

	// out
	dockerComposeFile *string
	deployInfo        map[string]string
}

func NewDeployCtx(
	options options.IOptions,
	writer *sse.Writer,
	env map[string]string,
) *DeployCtx {
	return &DeployCtx{
		options:    options,
		writer:     writer,
		env:        env,
		sendMu:     sync.Mutex{},
		deployInfo: make(map[string]string),
	}
}

func (d *DeployCtx) Send(event string, data string) error {
	d.sendMu.Lock()
	defer d.sendMu.Unlock()
	return d.writer.WriteEvent("", event, []byte(data))
}

func (d *DeployCtx) GetDeployInfo() map[string]string {
	return d.deployInfo
}

type DownCtx struct {
	Service    *model.Service
	NodeState  *node.NodeState
	deployInfo map[string]string
}

func NewDownCtx(
	service *model.Service,
	nodeState *node.NodeState,
	deployInfo map[string]string,
) *DownCtx {
	return &DownCtx{
		Service:    service,
		NodeState:  nodeState,
		deployInfo: deployInfo,
	}
}
