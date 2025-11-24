package node

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/benlocal/lai-panel/pkg/model"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type RemoteNodeExec struct {
	node *model.Node

	sftpClient *sftp.Client
	sshClient  *ssh.Client
}

func NewRemoteNodeExec(node *model.Node) *RemoteNodeExec {
	return &RemoteNodeExec{
		node: node,
	}
}

func (r *RemoteNodeExec) Init() error {
	password, err := r.node.GetDecryptedSSHPassword()
	if err != nil {
		return err
	}

	config := &ssh.ClientConfig{
		User: r.node.SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%d", r.node.Address, r.node.SSHPort)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return err
	}
	r.sshClient = sshClient

	client, err := sftp.NewClient(sshClient)
	if err != nil {
		return err
	}
	r.sftpClient = client

	return nil
}

func (r *RemoteNodeExec) Close() error {
	var err error
	if r.sshClient != nil {
		err = r.sshClient.Close()
	}
	if r.sftpClient != nil {
		err = r.sftpClient.Close()
	}
	return err
}

func (r *RemoteNodeExec) WriteFile(path string, data []byte) error {
	if r.sftpClient == nil {
		return fmt.Errorf("SFTP client not initialized")
	}

	dir := filepath.Dir(path)
	if dir != "" && dir != "." {
		if err := r.sftpClient.MkdirAll(dir); err != nil {
			return fmt.Errorf("failed to create remote directory: %w", err)
		}
	}

	file, err := r.sftpClient.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (r *RemoteNodeExec) ReadFile(path string) ([]byte, error) {
	if r.sftpClient == nil {
		return nil, fmt.Errorf("SFTP client not initialized")
	}

	file, err := r.sftpClient.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return data, nil
}

func (r *RemoteNodeExec) ExecuteOutput(command string, env map[string]string) (string, string, error) {
	stdout := ""
	stderr := ""
	err := r.ExecuteCommand(command, env, func(line string) {
		stdout += line + "\n"
	}, func(line string) {
		stderr += line + "\n"
	})
	return stdout, stderr, err
}

func (r *RemoteNodeExec) ExecuteCommand(
	command string,
	env map[string]string,
	onStdout func(string),
	onStderr func(string),
) error {
	if r.sshClient == nil {
		return fmt.Errorf("SSH client not initialized")
	}

	session, err := r.sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create SSH session: %w", err)
	}
	defer session.Close()

	stdoutPipe, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %w", err)
	}
	stderrPipe, err := session.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	fullCommand := buildRemoteCommand(command, env)
	if err := session.Start(fullCommand); err != nil {
		return fmt.Errorf("failed to start remote command: %w", err)
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if err := streamLines(stdoutPipe, onStdout); err != nil && onStderr != nil {
			onStderr(err.Error())
		}
	}()

	go func() {
		defer wg.Done()
		if err := streamLines(stderrPipe, onStderr); err != nil && onStderr != nil {
			onStderr(err.Error())
		}
	}()

	wg.Wait()

	if err := session.Wait(); err != nil {
		return err
	}

	return nil
}

func buildRemoteCommand(command string, env map[string]string) string {
	if len(env) == 0 {
		return fmt.Sprintf("bash -lc %q", command)
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var exports []string
	for _, key := range keys {
		value := env[key]
		exports = append(exports, fmt.Sprintf("export %s='%s'", key, escapeSingleQuotes(value)))
	}

	var buffer bytes.Buffer
	buffer.WriteString(strings.Join(exports, "; "))
	buffer.WriteString("; ")
	buffer.WriteString(command)

	return fmt.Sprintf("bash -lc %q", buffer.String())
}

func escapeSingleQuotes(value string) string {
	return strings.ReplaceAll(value, `'`, `'\''`)
}
