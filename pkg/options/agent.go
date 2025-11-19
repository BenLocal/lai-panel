package options

import (
	"path"

	"github.com/google/uuid"
)

type AgentOptions struct {
	Name       string
	Port       int
	MasterHost string
	MasterPort int
	dataPath   string
}

func NewAgentOptions() *AgentOptions {
	uuid := uuid.New().String()[:8]
	dataPath := getDefaultDataPath("agent")

	return &AgentOptions{
		Port:       8081,
		MasterHost: "127.0.0.1",
		MasterPort: 8080,
		Name:       uuid,
		dataPath:   dataPath,
	}
}

func WithAgentPort(port int) func(o *AgentOptions) {
	return func(o *AgentOptions) {
		o.Port = port
	}
}

func WithMasterHost(masterHost string) func(o *AgentOptions) {
	return func(o *AgentOptions) {
		o.MasterHost = masterHost
	}
}

func WithMasterPort(masterPort int) func(o *AgentOptions) {
	return func(o *AgentOptions) {
		o.MasterPort = masterPort
	}
}

func (o *AgentOptions) DataPath() string {
	return o.dataPath
}

func (o *AgentOptions) Agent() bool {
	return true
}

func (o *AgentOptions) StaticDataPath() string {
	return path.Join(o.dataPath, "static")
}

func (o *AgentOptions) LogDataPath() string {
	return path.Join(o.dataPath, "logs")
}
