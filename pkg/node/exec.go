package node

type NodeExec interface {
	Init() error
	Close() error
	WriteFile(path string, data []byte) error
	ReadFile(path string) ([]byte, error)
	ExecuteOutput(command string, env map[string]string) (string, string, error)
	ExecuteCommand(
		command string,
		env map[string]string,
		onStdout func(string),
		onStderr func(string),
	) error
}
