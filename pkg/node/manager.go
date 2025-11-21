package node

import (
	"sync"

	"github.com/benlocal/lai-panel/pkg/model"
)

type NodeManager struct {
	nodes map[int64]*NodeState

	mu sync.RWMutex
}

func NewNodeManager() *NodeManager {
	return &NodeManager{
		nodes: make(map[int64]*NodeState),
	}
}

func (m *NodeManager) AddOrGetNode(node *model.Node) (*NodeState, error) {
	m.mu.RLock()
	if state, ok := m.nodes[node.ID]; ok {
		m.mu.RUnlock()
		return state, nil
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()
	if state, ok := m.nodes[node.ID]; ok {
		return state, nil
	}

	state := NodeState{
		info: *node,
	}
	m.nodes[node.ID] = &state
	return &state, nil
}

func (m *NodeManager) RemoveNode(nodeID int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	state, ok := m.nodes[nodeID]
	if !ok {
		return nil
	}

	_ = state.Close()

	delete(m.nodes, nodeID)
	return nil
}
