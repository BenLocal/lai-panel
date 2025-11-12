package agent

type AgentOptions struct {
	Port       int
	MasterHost string
	MasterPort int
}

func NewAgentOptions() *AgentOptions {
	return &AgentOptions{
		Port:       8081,
		MasterHost: "127.0.0.1",
		MasterPort: 8080,
	}
}

func WithPort(port int) func(o *AgentOptions) {
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
