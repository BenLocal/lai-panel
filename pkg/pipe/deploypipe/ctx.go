package deploypipe

import (
	"github.com/benlocal/lai-panel/pkg/model"
	"github.com/benlocal/lai-panel/pkg/node"
)

type DeployCtx struct {
	App   *model.App
	Node  *model.Node
	state *node.NodeState
}

func NewDeployCtx(
	app *model.App,
	node *model.Node,
	state *node.NodeState,
) *DeployCtx {
	return &DeployCtx{
		Node:  node,
		state: state,
		App:   app,
	}
}
