package hub

import (
	"log"
	"sync"

	"github.com/benlocal/lai-panel/pkg/repository"
	"github.com/philippseith/signalr"
)

type SimpleHub struct {
	signalr.Hub
	nodeRepository *repository.NodeRepository

	sshSessions      map[string]*sshSessionState
	sshSessionsMutex sync.Mutex
}

func NewSimpleHub(nodeRepository *repository.NodeRepository) *SimpleHub {
	return &SimpleHub{
		nodeRepository:   nodeRepository,
		sshSessions:      make(map[string]*sshSessionState),
		sshSessionsMutex: sync.Mutex{},
	}
}

func (h *SimpleHub) OnConnected(connectionID string) {
	log.Printf("signalr connection connected: %s\n", connectionID)
}

func (h *SimpleHub) OnDisconnected(connectionID string) {
	log.Printf("signalr connection disconnected: %s\n", connectionID)
	h.stopSshSessionByID(connectionID)
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
