package node

import "io"

type NodeExec interface {
	Init() error
	Close() error
	WriteFile(path string, data []byte) error
	WriteFileStream(path string, reader io.Reader) error
	ReadFile(path string) ([]byte, error)
	ReadFileStream(path string, writer io.Writer) error
	ExecuteOutput(command string, opt *NodeExecuteCommandOptions) (string, string, error)
	ExecuteCommand(
		command string,
		opt *NodeExecuteCommandOptions,
		onStdout func(string),
		onStderr func(string),
	) error
}

type NodeExecuteCommandOptions struct {
	Env        map[string]string
	WorkingDir string
}

func NewNodeExecuteCommandOptions() *NodeExecuteCommandOptions {
	return &NodeExecuteCommandOptions{
		Env:        make(map[string]string),
		WorkingDir: "",
	}
}

func (o *NodeExecuteCommandOptions) SetEnv(env map[string]string) {
	o.Env = env
}

func (o *NodeExecuteCommandOptions) SetWorkingDir(workingDir string) {
	o.WorkingDir = workingDir
}
