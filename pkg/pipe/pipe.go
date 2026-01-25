package pipe

import (
	"context"

	"github.com/benlocal/lai-panel/pkg/pipe/deploypipe"
	"github.com/benlocal/lai-panel/pkg/pipe/nodepipe"
)

// Processor defines the interface for pipeline processors
type Processor[T any] interface {
	Process(ctx context.Context, input T) (T, error)
	Cancel(input T, err error)
}

// sequenceProcessor chains multiple processors together
type sequenceProcessor[T any] struct {
	processors []Processor[T]
}

func (s *sequenceProcessor[T]) Process(ctx context.Context, input T) (T, error) {
	var err error
	result := input
	for _, p := range s.processors {
		result, err = p.Process(ctx, result)
		if err != nil {
			// Cancel all previous processors in reverse order
			for i := len(s.processors) - 1; i >= 0; i-- {
				if s.processors[i] == p {
					break
				}
				s.processors[i].Cancel(result, err)
			}
			p.Cancel(result, err)
			return result, err
		}
	}
	return result, nil
}

func (s *sequenceProcessor[T]) Cancel(input T, err error) {
	// Cancel all processors in reverse order
	for i := len(s.processors) - 1; i >= 0; i-- {
		s.processors[i].Cancel(input, err)
	}
}

// Sequence creates a processor that executes multiple processors in sequence
func Sequence[T any](processors ...Processor[T]) Processor[T] {
	return &sequenceProcessor[T]{
		processors: processors,
	}
}

type NodePipeline struct {
	Processor Processor[*nodepipe.NodeCtx]
}

func NewNodePipeline() *NodePipeline {
	p := Sequence(&nodepipe.NodeCheckPipeline{})

	return &NodePipeline{
		Processor: p,
	}
}

func (p *NodePipeline) Run(ctx context.Context, nodeCtx *nodepipe.NodeCtx) (*nodepipe.NodeCtx, error) {
	return p.Processor.Process(ctx, nodeCtx)
}

type DeployPipeline struct {
	upPipeline   Processor[*deploypipe.DeployCtx]
	downPipeline Processor[*deploypipe.DownCtx]
}

func NewDeployPipeline() *DeployPipeline {
	up := Sequence(
		&deploypipe.CleanupWorkspacePipeline{},
		&deploypipe.CopyWorkspacePipeline{},
		&deploypipe.DownloadInstallerPipeline{},
		&deploypipe.DockerComposeFileParsePipeline{},
		&deploypipe.LoadImagePipeline{},
		&deploypipe.DockerComposeUpPipeline{},
	)

	down := Sequence(
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
