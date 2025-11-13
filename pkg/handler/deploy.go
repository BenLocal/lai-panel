package handler

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/benlocal/lai-panel/pkg/tmpl"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/sse"
	"gopkg.in/yaml.v3"
)

func (b *BaseHandler) HandleDockerComposeConfig(ctx context.Context, c *app.RequestContext) {
	type dockerComposeConfigRequest struct {
		DockerCompose string            `json:"docker_compose"`
		Env           map[string]string `json:"env"`
	}
	type dockerComposeConfigResponse struct {
		Config string `json:"config"`
	}
	var req dockerComposeConfigRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.Error(err)
		return
	}
	config, err := tmpl.ParseDockerCompose("test", req.DockerCompose, req.Env)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, dockerComposeConfigResponse{Config: config})
}

func (b *BaseHandler) HandleDockerComposeDeploy(ctx context.Context, c *app.RequestContext) {
	type dockerComposeDeployRequest struct {
		AppId  int64 `json:"app_id"`
		NodeId int64 `json:"node_id"`
	}
	var req dockerComposeDeployRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.Error(err)
		return
	}
	writer := sse.NewWriter(c)
	defer writer.Close()

	var sendMu sync.Mutex
	send := func(event string, data string) bool {
		sendMu.Lock()
		defer sendMu.Unlock()
		if err := writer.WriteEvent("", event, []byte(data)); err != nil {
			return false
		}
		return true
	}

	if req.AppId == 0 || req.NodeId == 0 {
		send("error", "app_id and node_id are required")
		return
	}

	if !send("message", "Fetching application information") {
		return
	}

	app, err := b.appRepository.GetByID(req.AppId)
	if err != nil {
		send("error", fmt.Sprintf("Failed to fetch application: %v", err))
		return
	}
	if app == nil {
		send("error", "Application not found")
		return
	}
	if app.DockerCompose == nil || len(strings.TrimSpace(*app.DockerCompose)) == 0 {
		send("error", "Application does not have docker compose configuration")
		return
	}

	if !send("message", "Loading node information") {
		return
	}

	node, err := b.nodeRepository.GetByID(req.NodeId)
	if err != nil {
		send("error", fmt.Sprintf("Failed to fetch node: %v", err))
		return
	}
	if node == nil {
		send("error", "Node not found")
		return
	}

	env := app.GetDefaultEnv()
	if env == nil {
		env = make(map[string]string)
	}
	env["APP_NAME"] = app.Name
	env["APP_VERSION"] = app.Version
	env["NODE_NAME"] = node.Name

	if !send("message", "Rendering docker compose template") {
		return
	}

	config, err := tmpl.ParseDockerCompose(app.Name, *app.DockerCompose, env)
	if err != nil {
		send("error", fmt.Sprintf("Failed to render docker compose template: %v", err))
		return
	}

	state, err := b.nodeManager.AddOrGetNode(node)
	if err != nil {
		send("error", fmt.Sprintf("Failed to initialize node executor: %v", err))
		return
	}

	_, _ = extractComposeImages(config)
	if err != nil {
		send("error", fmt.Sprintf("Failed to parse docker compose images: %v", err))
		return
	}

	filePath := fmt.Sprintf("/tmp/%s-%d-compose.yaml", sanitizeName(app.Name), time.Now().Unix())
	if !send("message", fmt.Sprintf("Uploading docker compose file to node (%s)", filePath)) {
		return
	}
	if err := state.Exec.WriteFile(filePath, []byte(config)); err != nil {
		send("error", fmt.Sprintf("Failed to upload docker compose file: %v", err))
		return
	}

	commands := []struct {
		Title   string
		Command string
	}{
		{
			Title:   "Pulling latest images",
			Command: fmt.Sprintf("docker compose -f %s pull", filePath),
		},
		{
			Title:   "Starting services",
			Command: fmt.Sprintf("docker compose -f %s up -d --remove-orphans", filePath),
		},
	}

	for _, step := range commands {
		if !send("message", fmt.Sprintf("%s...", step.Title)) {
			return
		}
		if err := state.Exec.ExecuteCommand(
			step.Command,
			env,
			func(line string) {
				if strings.TrimSpace(line) == "" {
					return
				}
				send("log", line)
			},
			func(line string) {
				if strings.TrimSpace(line) == "" {
					return
				}
				send("error", line)
			},
		); err != nil {
			send("error", fmt.Sprintf("Command failed: %v", err))
			return
		}
	}

	send("done", "Deployment completed successfully")

}

func extractComposeImages(config string) ([]string, error) {
	type composeService struct {
		Image string `yaml:"image"`
	}
	type composeDocument struct {
		Services map[string]composeService `yaml:"services"`
	}

	var doc composeDocument
	if err := yaml.Unmarshal([]byte(config), &doc); err != nil {
		return nil, err
	}

	imagesSet := make(map[string]struct{})
	for _, service := range doc.Services {
		image := strings.TrimSpace(service.Image)
		if image == "" {
			continue
		}
		imagesSet[image] = struct{}{}
	}

	images := make([]string, 0, len(imagesSet))
	for image := range imagesSet {
		images = append(images, image)
	}

	sort.Strings(images)
	return images, nil
}

func sanitizeName(name string) string {
	return strings.ReplaceAll(name, " ", "-")
}
