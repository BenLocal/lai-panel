package node

import (
	"sync"

	"github.com/benlocal/lai-panel/pkg/model"
	"github.com/benlocal/lai-panel/pkg/repository"
)

type NodeManager struct {
	nodes          map[int64]*NodeState
	nodeRepository *repository.NodeRepository

	mu sync.RWMutex
}

func NewNodeManager(nodeRepository *repository.NodeRepository) *NodeManager {
	return &NodeManager{
		nodeRepository: nodeRepository,
		nodes:          make(map[int64]*NodeState),
	}
}

func (m *NodeManager) GetNodeState(nodeID int64) (*NodeState, error) {
	m.mu.RLock()
	if state, ok := m.nodes[nodeID]; ok {
		m.mu.RUnlock()
		return state, nil
	}
	m.mu.RUnlock()

	node, err := m.nodeRepository.GetByID(nodeID)
	if err != nil {
		return nil, err
	}

	return m.addNode(node)
}

func (m *NodeManager) AddOrGetNode(node *model.Node) (*NodeState, error) {
	m.mu.RLock()
	if state, ok := m.nodes[node.ID]; ok {
		m.mu.RUnlock()
		return state, nil
	}
	m.mu.RUnlock()

	return m.addNode(node)
}

func (m *NodeManager) addNode(node *model.Node) (*NodeState, error) {
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
