package hub

import (
	"log"
	"sync"

	"github.com/benlocal/lai-panel/pkg/node"
	"github.com/benlocal/lai-panel/pkg/repository"
	"github.com/philippseith/signalr"
)

type SimpleHub struct {
	signalr.Hub
	nodeRepository *repository.NodeRepository
	nodeManager    *node.NodeManager

	sshSessions      map[string]*sshSessionState
	sshSessionsMutex sync.Mutex

	dockerSessions      map[string]*dockerSessionState
	dockerSessionsMutex sync.Mutex
}

func NewSimpleHub(nodeRepository *repository.NodeRepository,
	nodeManager *node.NodeManager) *SimpleHub {
	return &SimpleHub{
		nodeRepository:      nodeRepository,
		nodeManager:         nodeManager,
		sshSessions:         make(map[string]*sshSessionState),
		sshSessionsMutex:    sync.Mutex{},
		dockerSessions:      make(map[string]*dockerSessionState),
		dockerSessionsMutex: sync.Mutex{},
	}
}

func (h *SimpleHub) OnConnected(connectionID string) {
	log.Printf("signalr connection connected: %s\n", connectionID)
}

func (h *SimpleHub) OnDisconnected(connectionID string) {
	log.Printf("signalr connection disconnected: %s\n", connectionID)
	h.stopSshSessionByID(connectionID)
	h.stopDockerExecByID(connectionID)
}

func (h *SimpleHub) SendChatMessage(message string) {
	h.Clients().All().Send("chatMessageReceived", message)
}

// StartSshSession establishes an interactive SSH session for the current SignalR connection.
// nodeID identifies the target node, cols/rows configure the PTY size.
func (h *SimpleHub) StartSshSession(nodeID int64, cols int, rows int) error {
	connectionID := h.ConnectionID()
	return h.startSshSession(connectionID, nodeID, cols, rows)
}

// SendSshInput writes user input from the client to the remote SSH session.
func (h *SimpleHub) SendSshInput(data string) {
	connectionID := h.ConnectionID()
	h.onDataReceived(connectionID, []byte(data))
}

// ResizeSshSession resizes the remote PTY when the terminal window changes.
func (h *SimpleHub) ResizeSshSession(cols int, rows int) {
	connectionID := h.ConnectionID()
	h.resizeSshSession(connectionID, cols, rows)
}

// StopSshSession terminates the SSH session associated with the current connection.
func (h *SimpleHub) StopSshSession() {
	connectionID := h.ConnectionID()
	if h.stopSshSessionByID(connectionID) {
		h.Clients().Caller().Send("sshClosed", "")
	}
}

func (h *SimpleHub) StartDockerExec(nodeID int64, containerID string, cols int, rows int, shell string) error {
	connectionID := h.ConnectionID()
	return h.startDockerExec(connectionID, nodeID, containerID, cols, rows, shell)
}

func (h *SimpleHub) SendDockerExecInput(data string) {
	connectionID := h.ConnectionID()
	h.sendDockerExecInput(connectionID, data)
}

func (h *SimpleHub) ResizeDockerExec(cols int, rows int) {
	connectionID := h.ConnectionID()
	h.resizeDockerExec(connectionID, cols, rows)
}

func (h *SimpleHub) StopDockerExec() {
	connectionID := h.ConnectionID()
	if h.stopDockerExecByID(connectionID) {
		h.Clients().Caller().Send("dockerExecClosed", "")
	}
}
