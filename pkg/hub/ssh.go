package hub

import (
	"fmt"
	"log"

	"golang.org/x/crypto/ssh"
)

func (h *SimpleHub) StartSshSession(connectionID string, nodeID int64) {
	node, err := h.nodeRepository.GetByID(nodeID)
	if err != nil {
		log.Println("get node by id failed", err)
		return
	}
	h.sshSessionsMutex.Lock()
	defer h.sshSessionsMutex.Unlock()

	if _, ok := h.sshSessions[connectionID]; ok {
		log.Println("ssh session already exists", connectionID)
		return
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", node.Address, node.SSHPort), &ssh.ClientConfig{
		User: node.SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(node.SSHPassword),
		},
	})
	if err != nil {
		log.Println("ssh dial failed", err)
		return
	}

	session, err := client.NewSession()
	if err != nil {
		log.Println("new session failed", err)
		return
	}

	log.Println("ssh session started", connectionID, nodeID)
	h.sshSessions[connectionID] = session
}
