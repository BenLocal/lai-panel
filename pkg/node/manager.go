package node

import (
	"sync"

	"github.com/benlocal/lai-panel/pkg/docker"
	"github.com/benlocal/lai-panel/pkg/model"
	dockerClient "github.com/docker/docker/client"
)

type NodeState struct {
	info         model.Node
	Exec         NodeExec
	DockerClient *dockerClient.Client
}

type NodeManager struct {
	nodes map[int64]NodeState

	mu sync.RWMutex
}

func NewNodeManager() *NodeManager {
	return &NodeManager{
		nodes: make(map[int64]NodeState),
	}
}

func (m *NodeManager) AddOrGetNode(node *model.Node) (NodeState, error) {
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

	var exec NodeExec
	if node.IsLocal {
		exec = NewLocalNodeExec()
	} else {
		exec = NewRemoteNodeExec(node)
	}
	if err := exec.Init(); err != nil {
		return NodeState{}, err
	}

	var dockerClient *dockerClient.Client
	if node.IsLocal {
		dockerClient, _ = docker.LocalDockerClient()
	} else {
		dockerClient, _ = docker.AgentDockerClient(node.Address, node.AgentPort)
	}

	state := NodeState{
		info:         *node,
		Exec:         exec,
		DockerClient: dockerClient,
	}
	m.nodes[node.ID] = state
	return state, nil
}

func (m *NodeManager) RemoveNode(nodeID int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	state, ok := m.nodes[nodeID]
	if !ok {
		return nil
	}

	if closer, ok := state.Exec.(interface{ Close() error }); ok {
		if err := closer.Close(); err != nil {
			return err
		}
	}

	delete(m.nodes, nodeID)
	return nil
}
