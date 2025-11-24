package node

import "io"

type NodeExec interface {
	Init() error
	Close() error
	WriteFile(path string, data []byte) error
	WriteFileStream(path string, reader io.Reader) error
	ReadFile(path string) ([]byte, error)
	ReadFileStream(path string, writer io.Writer) error
	ExecuteOutput(command string, env map[string]string) (string, string, error)
	ExecuteCommand(
		command string,
		env map[string]string,
		onStdout func(string),
		onStderr func(string),
	) error
}
