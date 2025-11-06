package node

import (
	"fmt"
	"io"

	"github.com/benlocal/lai-panel/pkg/model"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type RemoteNodeExec struct {
	node *model.Node

	sftpClient *sftp.Client
	sshClient  *ssh.Client
}

// ExecuteCommand implements NodeExec.
func (r *RemoteNodeExec) ExecuteCommand(command string) (string, error) {
	panic("unimplemented")
}

func NewRemoteNodeExec(node *model.Node) *RemoteNodeExec {
	return &RemoteNodeExec{
		node: node,
	}
}

func (r *RemoteNodeExec) Init() error {
	config := &ssh.ClientConfig{
		User: r.node.SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(r.node.SSHPassword),
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
