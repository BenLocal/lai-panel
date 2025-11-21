package node

import (
	"fmt"
	"sync"

	"github.com/benlocal/lai-panel/pkg/docker"
	"github.com/benlocal/lai-panel/pkg/model"
	dockerClient "github.com/docker/docker/client"
)

type NodeState struct {
	info         model.Node
	exec         NodeExec
	dockerClient *dockerClient.Client

	execMu         sync.RWMutex
	dockerClientMu sync.RWMutex
}

func (n *NodeState) GetNodeInfo() string {
	return fmt.Sprintf("Node ID: %d, Node Name: %s, Node Address: %s", n.info.ID, n.info.Name, n.info.Address)
}

func (n *NodeState) GetNodeID() int64 {
	return n.info.ID
}

func (n *NodeState) GetDockerClient() (*dockerClient.Client, error) {
	n.dockerClientMu.RLock()
	if n.dockerClient != nil {
		n.dockerClientMu.RUnlock()
		return n.dockerClient, nil
	}
	n.dockerClientMu.RUnlock()

	var dockerClient *dockerClient.Client
	var err error
	if n.info.IsLocal {
		dockerClient, err = docker.LocalDockerClient()
	} else {
		dockerClient, err = docker.AgentDockerClient(n.info.Address, n.info.AgentPort)
	}
	if err != nil {
		return nil, err
	}

	n.dockerClientMu.Lock()
	defer n.dockerClientMu.Unlock()
	n.dockerClient = dockerClient
	return dockerClient, nil
}

func (n *NodeState) GetExec() (NodeExec, error) {
	n.execMu.RLock()
	if n.exec != nil {
		n.execMu.RUnlock()
		return n.exec, nil
	}
	n.execMu.RUnlock()

	var exec NodeExec
	if n.info.IsLocal {
		exec = NewLocalNodeExec()
	} else {
		exec = NewRemoteNodeExec(&n.info)
	}
	if err := exec.Init(); err != nil {
		return nil, err
	}

	n.execMu.Lock()
	defer n.execMu.Unlock()
	n.exec = exec
	return exec, nil
}

func (n *NodeState) Close() error {
	if n.exec != nil {
		if err := n.exec.Close(); err != nil {
			return err
		}
	}
	if n.dockerClient != nil {
		if err := n.dockerClient.Close(); err != nil {
			return err
		}
	}
	return nil
}
