package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/benlocal/lai-panel/pkg/client"
	"github.com/benlocal/lai-panel/pkg/handler"
	"github.com/benlocal/lai-panel/pkg/model"
	nodePkg "github.com/benlocal/lai-panel/pkg/node"
	"github.com/benlocal/lai-panel/pkg/options"
)

type HealthCheckService struct {
	context     context.Context
	cancel      context.CancelFunc
	baseClient  *client.BaseClient
	baseHandler *handler.BaseHandler
}

func NewHealthCheckService(baseClient *client.BaseClient,
	baseHandler *handler.BaseHandler) *HealthCheckService {
	ctx, cancel := context.WithCancel(context.Background())
	return &HealthCheckService{
		context:     ctx,
		cancel:      cancel,
		baseClient:  baseClient,
		baseHandler: baseHandler,
	}
}

func (s *HealthCheckService) Name() string {
	return "healthcheck"
}

func (s *HealthCheckService) Start(ctx context.Context) error {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-s.context.Done():
			return nil
		case <-ticker.C:
			if err := s.healthCheck(); err != nil {
				log.Println("health check failed", err)
			}
		}
	}
}

func (s *HealthCheckService) Shutdown() error {
	s.cancel()
	return nil
}

func (s *HealthCheckService) healthCheck() error {
	nodes, err := s.baseHandler.NodeRepository().List()
	if err != nil {
		return err
	}
	for _, node := range nodes {
		if !node.IsLocal {
			if err := s.baseClient.HealthCheck(node.Address, node.AgentPort); err != nil {
				log.Println("health check failed", node.Address, node.AgentPort, err)
				// update node status to offline
				s.baseHandler.NodeRepository().UpdateNodeStatus(node.ID, "offline")
				s.tryStartupAgentService(&node)
			} else {
				// update node status to online
				s.baseHandler.NodeRepository().UpdateNodeStatus(node.ID, "online")
			}
		}
	}
	return nil
}

func (s *HealthCheckService) tryStartupAgentService(node *model.Node) {
	state, err := s.baseHandler.NodeManager().GetNodeState(node.ID)
	if err != nil {
		return
	}
	exec, err := state.GetExec()
	if err != nil {
		return
	}

	currentBase := path.Join(s.baseHandler.StaticDataPath(), options.INSTALL_BASE_PATH)
	distBase := path.Join(*state.GetDataPath(), options.STATIC_BASE_PATH, options.INSTALL_BASE_PATH)

	// copy agent binary to data path
	agentPath := path.Join(currentBase, "agent")
	fi, err := os.Open(agentPath)
	if err != nil {
		return
	}
	defer fi.Close()
	disAgentPath := path.Join(distBase, "agent")
	err = exec.WriteFileStream(disAgentPath, fi)
	if err != nil {
		return
	}

	// copy install script to data path
	installScriptPath := path.Join(currentBase, "install.sh")
	fi, err = os.Open(installScriptPath)
	if err != nil {
		return
	}
	defer fi.Close()
	disPath := path.Join(distBase, "install.sh")
	err = exec.WriteFileStream(disPath, fi)
	if err != nil {
		return
	}

	// Make script executable
	opt := nodePkg.NewNodeExecuteCommandOptions()
	opt.WorkingDir = distBase
	if _, _, err := exec.ExecuteOutput("chmod +x install.sh", opt); err != nil {
		return
	}

	// Get master host and port
	masterHost := s.baseHandler.Options().MasterHost()
	masterPort := s.baseHandler.Options().MasterPort()

	// Build install command with parameters
	installCmd := fmt.Sprintf(
		"./install.sh --master-host %s --master-port %d --binary-path %s",
		masterHost,
		masterPort,
		disAgentPath,
	)

	// Add agent name if provided
	if node.Name != "" {
		installCmd += fmt.Sprintf(" --name %s", node.Name)
	}

	// Add agent address if provided
	if node.Address != "" {
		installCmd += fmt.Sprintf(" --address %s", node.Address)
	}

	// Add agent port if provided
	if node.AgentPort > 0 {
		installCmd += fmt.Sprintf(" --port %d", node.AgentPort)
	}

	// Execute install script
	err = exec.ExecuteCommand(installCmd, opt, func(s string) {
		log.Println("install script stdout", s)
	}, func(s string) {
		log.Println("install script stderr", s)
	})
	if err != nil {
		log.Println("failed to execute install script", err)
		return
	}
}
