package hub

import (
	"fmt"

	"github.com/philippseith/signalr"
)

type SimpleHub struct {
	signalr.Hub
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
