package deploypipe

import (
	"errors"
	"fmt"
	"path"
	"sync"

	"github.com/benlocal/lai-panel/pkg/ctx"
	"github.com/benlocal/lai-panel/pkg/model"
	"github.com/benlocal/lai-panel/pkg/node"
	"github.com/benlocal/lai-panel/pkg/options"
	"github.com/cloudwego/hertz/pkg/protocol/sse"
)

type DeployCtx struct {
	options     options.IOptions
	App         *model.App
	Service     *model.Service
	NodeState   *node.NodeState
	appCtx      *ctx.AppCtx
	writer      *sse.Writer
	sendMu      sync.Mutex
	env         map[string]string
	tmplFuncMap map[string]interface{}

	// out
	dockerComposeFile *string
	deployInfo        map[string]string
}

func NewDeployCtx(
	options options.IOptions,
	writer *sse.Writer,
	env map[string]string,
	appCtx *ctx.AppCtx,
) *DeployCtx {
	tmplFuncMap := builtinFuncMap(appCtx)
	return &DeployCtx{
		options:     options,
		appCtx:      appCtx,
		writer:      writer,
		env:         env,
		sendMu:      sync.Mutex{},
		deployInfo:  make(map[string]string),
		tmplFuncMap: tmplFuncMap,
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

func (d *DeployCtx) GetServicePath() (string, error) {
	return getPath(d.NodeState, d.options, d.Service.Name)
}

type DownCtx struct {
	options    options.IOptions
	Service    *model.Service
	NodeState  *node.NodeState
	deployInfo map[string]string
}

func NewDownCtx(
	options options.IOptions,
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

func (d *DownCtx) GetServicePath() (string, error) {
	return getPath(d.NodeState, d.options, d.Service.Name)
}

func getPath(nodeState *node.NodeState, opt options.IOptions, name string) (string, error) {
	if nodeState.IsLocal() {
		return path.Join(opt.DataPath(), options.SERVICE_BASE_PATH, name), nil
	}

	dataPath := nodeState.GetDataPath()
	if dataPath == nil {
		return "", errors.New("data path is not set")
	}
	return path.Join(*dataPath, options.SERVICE_BASE_PATH, name), nil
}

func builtinFuncMap(appCtx *ctx.AppCtx) map[string]interface{} {
	return map[string]interface{}{
		"panel_env": func(key interface{}, defaultValue interface{}) string {
			// Convert key to string
			keyStr := ""
			switch v := key.(type) {
			case string:
				keyStr = v
			default:
				keyStr = fmt.Sprintf("%v", v)
			}

			// Convert defaultValue to string
			defaultValueStr := ""
			switch v := defaultValue.(type) {
			case string:
				defaultValueStr = v
			default:
				defaultValueStr = fmt.Sprintf("%v", v)
			}

			e, err := appCtx.EnvRepository().GetByKey(keyStr)
			if err != nil {
				return defaultValueStr
			}

			if e == nil {
				return defaultValueStr
			}

			return e.Value
		},
		"is_agent": func() bool {
			return appCtx.Options().Agent()
		},
		"master_host": func() string {
			return appCtx.Options().MasterHost()
		},
		"master_port": func() int {
			return appCtx.Options().MasterPort()
		},
	}
}
