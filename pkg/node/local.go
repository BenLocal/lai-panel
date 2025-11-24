package node

import (
	"bufio"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
)

type LocalNodeExec struct {
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

func (l *LocalNodeExec) WriteFileStream(path string, reader io.Reader) error {
	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}

	return nil
}

func (l *LocalNodeExec) ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (l *LocalNodeExec) ReadFileStream(path string, writer io.Writer) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(writer, file)
	if err != nil {
		return err
	}
	return nil
}

func (l *LocalNodeExec) ExecuteOutput(command string, env map[string]string) (string, string, error) {
	stdout := ""
	stderr := ""
	err := l.ExecuteCommand(command, env, func(line string) {
		stdout += line + "\n"
	}, func(line string) {
		stderr += line + "\n"
	})
	return stdout, stderr, err
}

func (l *LocalNodeExec) ExecuteCommand(
	command string,
	env map[string]string,
	onStdout func(string),
	onStderr func(string),
) error {
	cmd := exec.Command("bash", "-c", command)

	if len(env) > 0 {
		envList := os.Environ()
		keys := make([]string, 0, len(env))
		for k := range env {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, key := range keys {
			envList = append(envList, key+"="+env[key])
		}
		cmd.Env = envList
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	done := make(chan error, 2)

	go func() {
		done <- streamLines(stdoutPipe, onStdout)
	}()

	go func() {
		done <- streamLines(stderrPipe, onStderr)
	}()

	for i := 0; i < 2; i++ {
		if streamErr := <-done; streamErr != nil {
			if onStderr != nil {
				onStderr(streamErr.Error())
			}
		}
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func streamLines(r io.Reader, handler func(string)) error {
	if handler == nil {
		// consume and discard
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
		}
		return scanner.Err()
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		handler(scanner.Text())
	}
	return scanner.Err()
}
