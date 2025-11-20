package nodepipe

import (
	"github.com/benlocal/lai-panel/pkg/model"
	"github.com/benlocal/lai-panel/pkg/node"
)

type NodeCtx struct {
	Node  *model.Node
	state *node.NodeState
}

func NewNodeCtx(node *model.Node, state *node.NodeState) *NodeCtx {
	return &NodeCtx{
		Node:  node,
		state: state,
	}
}
