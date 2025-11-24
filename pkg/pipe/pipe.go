package pipe

import (
	"context"

	"github.com/benlocal/lai-panel/pkg/pipe/deploypipe"
	"github.com/benlocal/lai-panel/pkg/pipe/nodepipe"
	"github.com/deliveryhero/pipeline/v2"
)

type NodePipeline struct {
	pipeline.Processor[*nodepipe.NodeCtx, *nodepipe.NodeCtx]
}

func NewNodePipeline() *NodePipeline {
	p := pipeline.Sequence(&nodepipe.NodeCheckPipeline{})

	return &NodePipeline{
		Processor: p,
	}
}

func (p *NodePipeline) Run(ctx context.Context, nodeCtx *nodepipe.NodeCtx) (*nodepipe.NodeCtx, error) {
	return p.Processor.Process(ctx, nodeCtx)
}

type DeployPipeline struct {
	upPipeline   pipeline.Processor[*deploypipe.DeployCtx, *deploypipe.DeployCtx]
	downPipeline pipeline.Processor[*deploypipe.DownCtx, *deploypipe.DownCtx]
}

func NewDeployPipeline() *DeployPipeline {
	up := pipeline.Sequence(
		&deploypipe.CopyInstallerPipeline{},
		&deploypipe.DockerComposeFileParsePipeline{},
		&deploypipe.LoadImagePipeline{},
		&deploypipe.DockerComposeUpPipeline{},
	)

	down := pipeline.Sequence(
		&deploypipe.DockerComposeDownPipeline{},
	)

	return &DeployPipeline{
		upPipeline:   up,
		downPipeline: down,
	}
}

func (p *DeployPipeline) Up(ctx context.Context, deployCtx *deploypipe.DeployCtx) (*deploypipe.DeployCtx, error) {
	return p.upPipeline.Process(ctx, deployCtx)
}

func (p *DeployPipeline) Down(ctx context.Context, downCtx *deploypipe.DownCtx) (*deploypipe.DownCtx, error) {
	return p.downPipeline.Process(ctx, downCtx)
}
