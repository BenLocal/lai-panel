package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/benlocal/lai-panel/pkg/model"
	"github.com/benlocal/lai-panel/pkg/pipe/deploypipe"
	"github.com/benlocal/lai-panel/pkg/tmpl"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/sse"
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
		ServiceId int64             `json:"service_id"`
		AppId     int64             `json:"app_id"`
		NodeId    int64             `json:"node_id"`
		QAValues  map[string]string `json:"qa_values"`
	}
	var req dockerComposeDeployRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.Error(err)
		return
	}
	writer := sse.NewWriter(c)
	defer writer.Close()

	deployCtx := deploypipe.NewDeployCtx(
		b.options,
		writer,
		req.QAValues,
	)

	service, err := b.ServiceRepository().GetByID(req.ServiceId)
	if err != nil {
		deployCtx.Send("error", err.Error())
		return
	}
	if service == nil {
		deployCtx.Send("error", "service not found")
		return
	}
	deployCtx.Service = service

	app, err := b.AppRepository().GetByID(req.AppId)
	if err != nil {
		deployCtx.Send("error", err.Error())
		return
	}
	if app == nil {
		deployCtx.Send("error", "app not found")
		return
	}
	deployCtx.App = app

	node, err := b.NodeRepository().GetByID(req.NodeId)
	if err != nil {
		deployCtx.Send("error", err.Error())
		return
	}
	if node == nil {
		deployCtx.Send("error", "node not found")
		return
	}
	state, err := b.NodeManager().AddOrGetNode(node)
	if err != nil {
		deployCtx.Send("error", err.Error())
		return
	}
	deployCtx.NodeState = state

	res, err := b.deployPipeline.Up(ctx, deployCtx)
	if err != nil {
		deployCtx.Send("error", err.Error())
		return
	}

	// update service deploy info and status
	err = b.updateServiceDeployInfo(service, res.GetDeployInfo())
	if err != nil {
		deployCtx.Send("error", err.Error())
		return
	}
}

func (b *BaseHandler) HandleDockerComposeUndeploy(ctx context.Context, c *app.RequestContext) {
	type dockerComposeUndeployRequest struct {
		ServiceId int64 `json:"service_id"`
	}
	var req dockerComposeUndeployRequest
	if err := c.BindAndValidate(&req); err != nil {
		c.Error(err)
		return
	}

	service, err := b.ServiceRepository().GetByID(req.ServiceId)
	if err != nil {
		c.Error(err)
		return
	}
	if service == nil {
		return
	}

	_, err = b.dockerComposeUndeploy(ctx, service)
	if err != nil {
		c.Error(err)
		return
	}

}

func (b *BaseHandler) updateServiceDeployInfo(service *model.Service, deployInfo map[string]string) error {
	jsonStr, err := json.Marshal(deployInfo)
	if err != nil {
		return err
	}
	j := string(jsonStr)
	db := model.Service{
		DeployInfo: &j,
		Status:     "running",
		ID:         service.ID,
	}

	return b.ServiceRepository().UpdateDeployInfo(&db)
}

func (b *BaseHandler) dockerComposeUndeploy(ctx context.Context, service *model.Service) (*deploypipe.DownCtx, error) {
	var deployInfo map[string]string
	err := json.Unmarshal([]byte(*service.DeployInfo), &deployInfo)
	if err != nil {
		return nil, err
	}
	node, err := b.NodeRepository().GetByID(service.NodeID)
	if err != nil {
		return nil, err
	}
	if node == nil {
		return nil, errors.New("node not found")
	}
	state, err := b.NodeManager().AddOrGetNode(node)
	if err != nil {
		return nil, err
	}

	downCtx := deploypipe.NewDownCtx(service, state, deployInfo)
	res, err := b.deployPipeline.Down(ctx, downCtx)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// func extractComposeImages(config string) ([]string, error) {
// 	type composeService struct {
// 		Image string `yaml:"image"`
// 	}
// 	type composeDocument struct {
// 		Services map[string]composeService `yaml:"services"`
// 	}

// 	var doc composeDocument
// 	if err := yaml.Unmarshal([]byte(config), &doc); err != nil {
// 		return nil, err
// 	}

// 	imagesSet := make(map[string]struct{})
// 	for _, service := range doc.Services {
// 		image := strings.TrimSpace(service.Image)
// 		if image == "" {
// 			continue
// 		}
// 		imagesSet[image] = struct{}{}
// 	}

// 	images := make([]string, 0, len(imagesSet))
// 	for image := range imagesSet {
// 		images = append(images, image)
// 	}

// 	sort.Strings(images)
// 	return images, nil
// }

// func sanitizeName(name string) string {
// 	return strings.ReplaceAll(name, " ", "-")
// }
