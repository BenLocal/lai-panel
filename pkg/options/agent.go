package options

import (
	"os"
	"path"
	"strings"

	"github.com/google/uuid"
)

type AgentOptions struct {
	Name       string
	Port       int
	MasterHost string
	MasterPort int
	Address    string
	dataPath   string
}

func NewAgentOptions(opts ...func(o *AgentOptions)) *AgentOptions {
	uuid := uuid.New().String()[:8]
	dataPath := getDefaultDataPath("data")

	t := &AgentOptions{
		Port:       8081,
		MasterHost: "127.0.0.1",
		MasterPort: 8080,
		Name:       uuid,
		dataPath:   dataPath,
	}

	for _, f := range opts {
		f(t)
	}

	return t
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

func WithAddress(address string) func(o *AgentOptions) {
	return func(o *AgentOptions) {
		o.Address = address
	}
}

func WithName(name string) func(o *AgentOptions) {
	return func(o *AgentOptions) {
		if name == "" {
			cachePath := path.Join(o.dataPath, "name")
			name, recreate := readNameFromCache(cachePath)
			if name == "" {
				name = uuid.New().String()[:8]
			}
			if recreate {
				os.WriteFile(cachePath, []byte(name), 0644)
			}
		}

		o.Name = name
	}
}

func readNameFromCache(cachePath string) (string, bool) {
	name := ""
	recreate := false
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		recreate = true
	} else {
		data, err := os.ReadFile(cachePath)
		if err == nil {
			name = strings.TrimSpace(string(data))
		} else {
			recreate = true
		}
	}

	return name, recreate
}

func (o *AgentOptions) DataPath() string {
	return o.dataPath
}

func (o *AgentOptions) Agent() bool {
	return true
}
