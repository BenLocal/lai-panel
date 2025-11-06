package node

import (
	"os"
	"path/filepath"
)

type LocalNodeExec struct {
}

// ExecuteCommand implements NodeExec.
func (l *LocalNodeExec) ExecuteCommand(command string) (string, error) {
	panic("unimplemented")
}

func NewLocalNodeExec() *LocalNodeExec {
	return &LocalNodeExec{}
}

func (l *LocalNodeExec) Init() error {
	return nil
}

func (l *LocalNodeExec) Close() error {
	return nil
}

func (l *LocalNodeExec) WriteFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return os.WriteFile(path, data, 0o644)
}

func (l *LocalNodeExec) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}
