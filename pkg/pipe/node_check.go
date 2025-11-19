package pipe

import (
	"context"
	"errors"
)

type NodeCheckPipeline struct {
}

func (p *NodeCheckPipeline) Process(ctx context.Context, nodeCtx *NodeCtx) (*NodeCtx, error) {
	if nodeCtx.Node.ID <= 0 {
		return nil, errors.New("node is invalid")
	}

	out, errout, err := nodeCtx.state.Exec.ExecuteOutput("docker ps", nil)
	if err != nil {
		return nil, err
	}

	if errout != "" {
		return nil, errors.New(errout)
	}

	if out == "" {
		return nil, errors.New("docker is not running")
	}

	return nodeCtx, nil
}

func (p *NodeCheckPipeline) Cancel(nodeCtx *NodeCtx, err error) {
	// do nothing
}
