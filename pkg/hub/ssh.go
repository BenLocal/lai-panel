package hub

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/benlocal/lai-panel/pkg/model"
	"github.com/creack/pty"
	"golang.org/x/crypto/ssh"
)

type sshSessionState struct {
	session   *ssh.Session
	client    *ssh.Client
	cmd       *exec.Cmd
	writer    io.WriteCloser
	readers   []io.ReadCloser
	pty       *os.File
	closeOnce sync.Once
}

type nopReadCloser struct {
	io.Reader
}

func (nopReadCloser) Close() error { return nil }

func (h *SimpleHub) startSshSession(connectionID string, nodeID int64, cols int, rows int) error {
	h.sshSessionsMutex.Lock()
	if _, ok := h.sshSessions[connectionID]; ok {
		h.sshSessionsMutex.Unlock()
		return fmt.Errorf("ssh session already exists for connection %s", connectionID)
	}
	h.sshSessionsMutex.Unlock()

	targetNode, err := h.nodeRepository.GetByID(nodeID)
	if err != nil {
		log.Printf("failed to load node %d: %v\n", nodeID, err)
		return err
	}

	var state *sshSessionState

	if targetNode.IsLocal {
		cmd, ptyFile, err := createLocalSession(cols, rows)
		if err != nil {
			log.Printf("failed to start local pty: %v\n", err)
			return err
		}
		state = &sshSessionState{
			cmd:     cmd,
			writer:  ptyFile,
			readers: []io.ReadCloser{nopReadCloser{Reader: ptyFile}},
			pty:     ptyFile,
		}
	} else {
		sshClient, session, stdin, stdout, stderr, err := createRemoteSession(targetNode, cols, rows)
		if err != nil {
			log.Printf("failed to create ssh session: %v\n", err)
			return err
		}
		state = &sshSessionState{
			session: session,
			client:  sshClient,
			writer:  stdin,
			readers: []io.ReadCloser{stdout, stderr},
		}
	}

	h.sshSessionsMutex.Lock()
	h.sshSessions[connectionID] = state
	h.sshSessionsMutex.Unlock()

	for _, reader := range state.readers {
		go h.streamOutput(connectionID, reader)
	}
	go h.waitSession(connectionID, state)

	return nil
}

func (h *SimpleHub) streamOutput(connectionID string, reader io.ReadCloser) {
	defer reader.Close()

	buf := make([]byte, 4096)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			chunk := string(buf[:n])
			h.Clients().Caller().Send("sshData", chunk)
		}
		if err != nil {
			if err != io.EOF {
				log.Printf("ssh read error (%s): %v\n", connectionID, err)
			}
			break
		}
	}
}

func (h *SimpleHub) waitSession(connectionID string, state *sshSessionState) {
	var err error
	if state.session != nil {
		err = state.session.Wait()
	} else if state.cmd != nil {
		err = state.cmd.Wait()
	}
	if err != nil && err != io.EOF {
		log.Printf("ssh session ended with error (%s): %v\n", connectionID, err)
	}

	if h.stopSshSessionByID(connectionID) {
		h.Clients().Caller().Send("sshClosed", "")
	}
}

func (h *SimpleHub) resizeSshSession(connectionID string, cols int, rows int) {
	if cols <= 0 {
		cols = 120
	}
	if rows <= 0 {
		rows = 32
	}

	h.sshSessionsMutex.Lock()
	state, ok := h.sshSessions[connectionID]
	h.sshSessionsMutex.Unlock()
	if !ok {
		return
	}

	if state.session != nil {
		if err := state.session.WindowChange(rows, cols); err != nil {
			log.Printf("failed to resize remote pty (%s): %v\n", connectionID, err)
		}
	} else if state.pty != nil {
		if err := pty.Setsize(state.pty, &pty.Winsize{Rows: uint16(rows), Cols: uint16(cols)}); err != nil {
			log.Printf("failed to resize local pty (%s): %v\n", connectionID, err)
		}
	}
}

func (h *SimpleHub) onDataReceived(connectionID string, data []byte) {
	h.sshSessionsMutex.Lock()
	state, ok := h.sshSessions[connectionID]
	h.sshSessionsMutex.Unlock()

	if !ok || state.writer == nil {
		return
	}

	if _, err := state.writer.Write(data); err != nil {
		log.Printf("failed to write to ssh stdin (%s): %v\n", connectionID, err)
	}
}

func (h *SimpleHub) stopSshSessionByID(connectionID string) bool {
	h.sshSessionsMutex.Lock()
	state, ok := h.sshSessions[connectionID]
	if ok {
		delete(h.sshSessions, connectionID)
	}
	h.sshSessionsMutex.Unlock()

	if !ok {
		return false
	}

	state.closeOnce.Do(func() {
		if state.writer != nil {
			_ = state.writer.Close()
		}
		for _, r := range state.readers {
			_ = r.Close()
		}
		if state.session != nil {
			_ = state.session.Close()
		}
		if state.client != nil {
			_ = state.client.Close()
		}
		if state.cmd != nil && state.cmd.Process != nil {
			_ = state.cmd.Process.Kill()
		}
		if state.pty != nil {
			_ = state.pty.Close()
		}
	})
	return true
}

func createRemoteSession(node *model.Node, cols int, rows int) (*ssh.Client, *ssh.Session, io.WriteCloser, io.ReadCloser, io.ReadCloser, error) {
	if cols <= 0 {
		cols = 120
	}
	if rows <= 0 {
		rows = 32
	}

	pwd, err := node.GetDecryptedSSHPassword()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	clientConfig := &ssh.ClientConfig{
		User:            node.SSHUser,
		Auth:            []ssh.AuthMethod{ssh.Password(pwd)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", node.Address, node.SSHPort), clientConfig)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		client.Close()
		return nil, nil, nil, nil, nil, err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err := session.RequestPty("xterm-256color", rows, cols, modes); err != nil {
		session.Close()
		client.Close()
		return nil, nil, nil, nil, nil, err
	}

	stdin, err := session.StdinPipe()
	if err != nil {
		session.Close()
		client.Close()
		return nil, nil, nil, nil, nil, err
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		stdin.Close()
		session.Close()
		client.Close()
		return nil, nil, nil, nil, nil, err
	}

	stderr, err := session.StderrPipe()
	if err != nil {
		stdin.Close()
		session.Close()
		client.Close()
		return nil, nil, nil, nil, nil, err
	}

	if err := session.Shell(); err != nil {
		stdin.Close()
		session.Close()
		client.Close()
		return nil, nil, nil, nil, nil, err
	}

	return client, session, stdin, io.NopCloser(stdout), io.NopCloser(stderr), nil
}

func createLocalSession(cols int, rows int) (*exec.Cmd, *os.File, error) {
	if cols <= 0 {
		cols = 120
	}
	if rows <= 0 {
		rows = 32
	}

	shell := os.Getenv("SHELL")
	if shell == "" {
		if runtime.GOOS == "windows" {
			shell = "powershell.exe"
		} else {
			shell = "/bin/sh"
		}
	}

	cmd := exec.Command(shell)
	env := os.Environ()
	hasTerm := false
	hasHome := false
	var homeDir string
	for _, kv := range env {
		if len(kv) >= 5 && kv[:5] == "TERM=" {
			hasTerm = true
		}
		if len(kv) >= 5 && kv[:5] == "HOME=" {
			hasHome = true
			// 提取 HOME 路径
			homeDir = kv[5:]
		}
	}
	if !hasTerm {
		env = append(env, "TERM=xterm-256color")
	}
	if !hasHome {
		// 确保 HOME 路径被设置
		var err error
		homeDir, err = os.UserHomeDir()
		if err == nil {
			env = append(env, "HOME="+homeDir)
		}
	}
	cmd.Env = env
	// 设置工作目录为 HOME 目录
	if homeDir != "" {
		cmd.Dir = homeDir
	}
	ptyFile, err := pty.StartWithSize(cmd, &pty.Winsize{Rows: uint16(rows), Cols: uint16(cols)})
	if err != nil {
		return nil, nil, err
	}

	return cmd, ptyFile, nil
}
