package hub

import (
	"fmt"
	"sync"

	"github.com/benlocal/lai-panel/pkg/repository"
	"github.com/philippseith/signalr"
	"golang.org/x/crypto/ssh"
)

type SimpleHub struct {
	signalr.Hub
	nodeRepository *repository.NodeRepository

	sshSessions      map[string]*ssh.Session
	sshSessionsMutex sync.Mutex
}

func NewSimpleHub(nodeRepository *repository.NodeRepository) *SimpleHub {
	return &SimpleHub{
		nodeRepository:   nodeRepository,
		sshSessions:      make(map[string]*ssh.Session),
		sshSessionsMutex: sync.Mutex{},
	}
}

func (h *SimpleHub) OnConnected(connectionID string) {
	fmt.Println("OnConnected", connectionID)
}

func (h *SimpleHub) OnDisconnected(connectionID string) {
	fmt.Println("OnDisconnected", connectionID)
}

func (h *SimpleHub) SendChatMessage(message string) {
	h.Clients().All().Send("chatMessageReceived", message)
}
