package hub

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"

	"github.com/benlocal/lai-panel/pkg/node"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

// readerCloser 将 bufio.Reader 和 net.Conn 组合成 io.ReadCloser
type readerCloser struct {
	reader *bufio.Reader
	conn   net.Conn
}

func (r *readerCloser) Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}

func (r *readerCloser) Close() error {
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}

type dockerSessionState struct {
	sessionID   string
	writer      io.WriteCloser
	nodeID      int64
	containerID string
	nodeState   *node.NodeState
}

func (h *SimpleHub) startDockerExec(connectionID string, nodeID int64, containerID string, cols int, rows int, shell string) error {
	h.dockerSessionsMutex.Lock()
	if _, ok := h.dockerSessions[connectionID]; ok {
		h.dockerSessionsMutex.Unlock()
		return fmt.Errorf("docker session already exists for connection %s", connectionID)
	}
	h.dockerSessionsMutex.Unlock()

	// 默认使用 sh，如果 shell 为空
	if shell == "" || strings.TrimSpace(shell) == "" {
		shell = "sh"
	}
	// 清理 shell 字符串，移除前后空格
	shell = strings.TrimSpace(shell)

	targetNode, err := h.nodeRepository.GetByID(nodeID)
	if err != nil {
		log.Printf("failed to get node %d: %v\n", nodeID, err)
		return err
	}
	nodeState, err := h.nodeManager.AddOrGetNode(targetNode)
	if err != nil {
		log.Printf("failed to add or get node %d: %v\n", nodeID, err)
		return err
	}
	id, resp, err := createDockerExecSession(nodeState, containerID, rows, cols, shell)
	if err != nil {
		log.Printf("failed to create docker exec session: %v\n", err)
		return err
	}
	if id == "" {
		return fmt.Errorf("failed to create docker exec session")
	}

	reader := &readerCloser{
		reader: resp.Reader,
		conn:   resp.Conn,
	}

	go h.streamDockerExecOutput(connectionID, reader)

	state := &dockerSessionState{
		nodeID:      nodeID,
		containerID: containerID,
		writer:      resp.Conn,
		sessionID:   id,
		nodeState:   nodeState,
	}
	h.dockerSessionsMutex.Lock()
	h.dockerSessions[connectionID] = state
	h.dockerSessionsMutex.Unlock()

	return nil
}

func (h *SimpleHub) streamDockerExecOutput(connectionID string, reader io.ReadCloser) {
	defer reader.Close()

	buf := make([]byte, 4096)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			chunk := string(buf[:n])
			h.Clients().Caller().Send("dockerExecData", chunk)
		}
		if err != nil {
			if err != io.EOF {
				log.Printf("docker exec read error (%s): %v\n", connectionID, err)
			}
			h.handleDockerExecDisconnected(connectionID)
			break
		}
	}
}

func (h *SimpleHub) handleDockerExecDisconnected(connectionID string) {
	if h.stopDockerExecByID(connectionID) {
		h.Clients().Caller().Send("dockerExecClosed", "")
	}
}

func (h *SimpleHub) sendDockerExecInput(connectionID string, data string) {
	h.dockerSessionsMutex.Lock()
	state, ok := h.dockerSessions[connectionID]
	h.dockerSessionsMutex.Unlock()

	if !ok || state.writer == nil {
		return
	}

	if _, err := state.writer.Write([]byte(data)); err != nil {
		log.Printf("failed to write to docker exec stdin (%s): %v\n", connectionID, err)
		h.handleDockerExecDisconnected(connectionID)
	}
}

func (h *SimpleHub) resizeDockerExec(connectionID string, cols int, rows int) {
	h.dockerSessionsMutex.Lock()
	state, ok := h.dockerSessions[connectionID]
	h.dockerSessionsMutex.Unlock()

	if !ok || state.sessionID == "" {
		return
	}
	nodeState := state.nodeState
	if nodeState == nil {
		return
	}

	err := nodeState.DockerClient.ContainerExecResize(context.Background(), state.sessionID, container.ResizeOptions{
		Height: uint(rows),
		Width:  uint(cols),
	})
	if err != nil {
		log.Printf("failed to resize docker exec: %v\n", err)
		return
	}

	log.Printf("resize docker exec (%s): cols=%d, rows=%d\n", connectionID, cols, rows)
}

func (h *SimpleHub) stopDockerExecByID(connectionID string) bool {
	h.dockerSessionsMutex.Lock()
	state, ok := h.dockerSessions[connectionID]
	if ok {
		delete(h.dockerSessions, connectionID)
	}
	h.dockerSessionsMutex.Unlock()

	if !ok {
		return false
	}

	if state.writer != nil {
		_ = state.writer.Close()
	}

	return true
}

func createDockerExecSession(state *node.NodeState, containerID string, rows int, cols int, shell string) (string, *types.HijackedResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 默认使用 sh，如果 shell 为空
	if shell == "" || strings.TrimSpace(shell) == "" {
		shell = "sh"
	}
	// 清理 shell 字符串，移除前后空格
	shell = strings.TrimSpace(shell)

	resp, err := state.DockerClient.ContainerExecCreate(ctx, containerID, container.ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		AttachStdin:  true,
		DetachKeys:   "ctrl-p,ctrl-q",
		Env:          []string{},
		Cmd:          []string{shell},
		Tty:          true,
	})
	if err != nil {
		log.Printf("failed to create docker exec: %v\n", err)
		return "", nil, err
	}

	if resp.ID == "" {
		return "", nil, fmt.Errorf("failed to create docker exec session")
	}

	// err = state.DockerClient.ContainerExecStart(ctx, resp.ID, container.ExecStartOptions{
	// 	Detach: false,
	// 	Tty:    false,
	// })
	// if err != nil {
	// 	errMsg := err.Error()
	// 	if strings.Contains(errMsg, "already running") {
	// 		log.Printf("docker exec %s is already running, checking status\n", resp.ID)
	// 		inspect, inspectErr := state.DockerClient.ContainerExecInspect(ctx, resp.ID)
	// 		if inspectErr == nil && inspect.Running {
	// 			log.Printf("docker exec %s is running, will attach to it\n", resp.ID)
	// 			return resp.ID, nil
	// 		}
	// 		log.Printf("docker exec %s status check failed or not running: inspectErr=%v, running=%v\n",
	// 			resp.ID, inspectErr, inspect.Running)
	// 	}
	// 	log.Printf("failed to start docker exec: %v\n", err)
	// 	return "", err
	// }

	v, err := state.DockerClient.ContainerExecAttach(context.Background(), resp.ID, container.ExecAttachOptions{
		Tty:         true,
		ConsoleSize: &[2]uint{uint(rows), uint(cols)},
	})
	if err != nil {
		log.Printf("failed to attach to docker exec: %v\n", err)
		return "", nil, err
	}

	return resp.ID, &v, nil
}
