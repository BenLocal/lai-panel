package handler

import (
	"bufio"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/benlocal/lai-panel/pkg/tmpl"
	"github.com/valyala/fasthttp"
)

func (b *BaseHandler) HandleDockerComposeConfig(ctx *fasthttp.RequestCtx) {
	type dockerComposeConfigRequest struct {
		DockerCompose string            `json:"docker_compose"`
		Env           map[string]string `json:"env"`
	}
	type dockerComposeConfigResponse struct {
		Config string `json:"config"`
	}
	var req dockerComposeConfigRequest
	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		JSONError(ctx, "invalid request", err)
		return
	}
	config, err := tmpl.ParseDockerCompose("test", req.DockerCompose, req.Env)
	if err != nil {
		JSONError(ctx, "failed to parse docker compose template", err)
		return
	}
	JSONSuccess(ctx, dockerComposeConfigResponse{Config: config})
}

func (b *BaseHandler) HandleDockerComposeByApp(ctx *fasthttp.RequestCtx) {
	type dockerComposeConfigRequest struct {
		AppId int64 `json:"app_id"`
	}
	type dockerComposeConfigResponse struct {
		Config string `json:"config"`
	}
	var req dockerComposeConfigRequest
	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		JSONError(ctx, "invalid request", err)
		return
	}
	app, err := b.appRepository.GetByID(req.AppId)
	if err != nil {
		JSONError(ctx, "app not found", err)
		return
	}

	env := app.GetEnv()
	config, err := tmpl.ParseDockerCompose(app.Name, *app.DockerCompose, env)
	if err != nil {
		JSONError(ctx, "failed to parse docker compose template", err)
		return
	}
	JSONSuccess(ctx, dockerComposeConfigResponse{Config: config})
}

func (b *BaseHandler) HandleDockerComposeDeploy(ctx *fasthttp.RequestCtx) {
	type dockerComposeDeployRequest struct {
		AppId  int64 `json:"app_id"`
		NodeId int64 `json:"node_id"`
	}
	var req dockerComposeDeployRequest
	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		JSONError(ctx, "invalid request", err)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.Set("Content-Type", "text/event-stream")
	ctx.Response.Header.Set("Cache-Control", "no-cache")
	ctx.Response.Header.Set("Connection", "keep-alive")

	ctx.SetBodyStreamWriter(func(writer *bufio.Writer) {
		var sendMu sync.Mutex
		send := func(event string, data string) bool {
			sendMu.Lock()
			defer sendMu.Unlock()
			if err := writeSSE(writer, event, data); err != nil {
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

		env := app.GetEnv()
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
	})
}

func writeSSE(writer *bufio.Writer, event, data string) error {
	if event != "" {
		if _, err := writer.WriteString(fmt.Sprintf("event: %s\n", event)); err != nil {
			return err
		}
	}

	lines := strings.Split(data, "\n")
	if len(lines) == 0 {
		lines = []string{""}
	}
	for _, line := range lines {
		if _, err := writer.WriteString(fmt.Sprintf("data: %s\n", line)); err != nil {
			return err
		}
	}
	if _, err := writer.WriteString("\n"); err != nil {
		return err
	}
	return writer.Flush()
}

var invalidNamePattern = regexp.MustCompile(`[^a-zA-Z0-9_-]`)

func sanitizeName(name string) string {
	if name == "" {
		return "app"
	}
	clean := invalidNamePattern.ReplaceAllString(name, "_")
	clean = strings.Trim(clean, "_")
	if clean == "" {
		return "app"
	}
	return clean
}
