package pipe

import (
	"context"

	"github.com/deliveryhero/pipeline/v2"
)

type NodePipeline struct {
	pipeline.Processor[*NodeCtx, *NodeCtx]
}

func NewNodePipeline() *NodePipeline {
	p := pipeline.Sequence(&NodeCheckPipeline{})

	return &NodePipeline{
		Processor: p,
	}
}

func (p *NodePipeline) Run(ctx context.Context, nodeCtx *NodeCtx) (*NodeCtx, error) {
	return p.Processor.Process(ctx, nodeCtx)
}
