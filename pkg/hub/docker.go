package hub

import (
	"fmt"
)

type dockerSessionState struct {
	sessionID string
}

func (h *SimpleHub) startDockerExec(connectionID string, nodeID int64, containerID string, cols int, rows int) error {
	h.dockerSessionsMutex.Lock()
	if _, ok := h.dockerSessions[connectionID]; ok {
		h.dockerSessionsMutex.Unlock()
		return fmt.Errorf("docker session already exists for connection %s", connectionID)
	}
	h.dockerSessions[connectionID] = &dockerSessionState{}
	h.dockerSessionsMutex.Unlock()

	// targetNode, err := h.nodeRepository.GetByID(nodeID)
	// if err != nil {
	// 	log.Printf("failed to load node %d: %v\n", nodeID, err)
	// 	return err
	// }
	// state, err := h.nodeManager.AddOrGetNode(targetNode)
	// if err != nil {
	// 	log.Printf("failed to add or get node %d: %v\n", nodeID, err)
	// 	return err
	// }
	// resp, err := state.DockerClient.ContainerExecCreate(context.Background(), containerID, container.ExecOptions{
	// 	AttachStdout: true,
	// 	AttachStderr: true,
	// 	AttachStdin:  true,
	// 	Detach:       false,
	// 	DetachKeys:   "ctrl-p,ctrl-q",
	// 	Env:          []string{},
	// 	Cmd:          []string{"sh"},
	// })
	// if err != nil {
	// 	log.Printf("failed to create docker exec: %v\n", err)
	// 	return err
	// }

	// err = state.DockerClient.ContainerExecStart(context.Background(), resp.ID, container.ExecStartOptions{
	// 	Detach:      false,
	// 	Tty:         true,
	// 	ConsoleSize: &[2]uint{uint(rows), uint(cols)},
	// })
	// if err != nil {
	// 	log.Printf("failed to start docker exec: %v\n", err)
	// 	return err
	// }

	// re, err := state.DockerClient.ContainerExecAttach(context.Background(), resp.ID, container.ExecAttachOptions{
	// 	Tty:         true,
	// 	ConsoleSize: &[2]uint{uint(rows), uint(cols)},
	// })

	// go func() {
	// 	for {
	// 		buf := make([]byte, 4096)
	// 		n, err := re.Reader.Read(buf)
	// 		if err != nil {
	// 			log.Printf("failed to read from docker exec: %v\n", err)
	// 			return
	// 		}
	// 	}
	// }()

	// if err != nil {
	// 	log.Printf("failed to attach to docker exec: %v\n", err)
	// 	return err
	// }

	return nil
}
