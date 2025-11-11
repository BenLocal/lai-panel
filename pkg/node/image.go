package node

import (
	"context"
	"fmt"
	"io"
)

func CopyImageBetweenNodes(
	ctx context.Context,
	sourceState *NodeState,
	destState *NodeState,
	image string,
	cb func(ctx context.Context, reader io.ReadCloser) error,
) error {
	reader, err := sourceState.DockerClient.ImageSave(ctx, []string{image})
	if err != nil {
		return fmt.Errorf("failed to export image from source node: %w", err)
	}
	defer reader.Close()

	loadResp, err := destState.DockerClient.ImageLoad(ctx, reader)
	if err != nil {
		return fmt.Errorf("failed to import image to destination node: %w", err)
	}

	if cb != nil {
		if err := cb(ctx, loadResp.Body); err != nil {
			return fmt.Errorf("failed to copy image: %w", err)
		}
	}
	defer loadResp.Body.Close()

	return nil
}
