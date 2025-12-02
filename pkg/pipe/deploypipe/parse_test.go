package deploypipe

import (
	"testing"

	"github.com/benlocal/lai-panel/pkg/constant"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestEditFile_AddLabelsToServiceWithoutLabels(t *testing.T) {
	dockerCompose := `services:
  nginx:
    image: nginx:latest
    ports:
      - 80:80
`
	pipeline := &DockerComposeFileParsePipeline{}
	labels := map[string]string{
		constant.ManagedByLabel: constant.ProjectId,
		constant.OwnerLabel:     constant.ProjectId,
		constant.ServiceLabel:   "nginx-service",
	}

	res, err := pipeline.editFile(dockerCompose, labels)
	assert.NoError(t, err)
	assert.NotEmpty(t, res)

	// Parse the result to verify labels were added
	var y yaml.Node
	err = yaml.Unmarshal([]byte(res), &y)
	assert.NoError(t, err)

	root := &y
	if y.Kind == yaml.DocumentNode && len(y.Content) > 0 {
		root = y.Content[0]
	}

	services := lookup(root, "services")
	assert.NotNil(t, services)

	nginx := lookup(services, "nginx")
	assert.NotNil(t, nginx)

	labelsNode := lookup(nginx, "labels")
	assert.NotNil(t, labelsNode)
	assert.Equal(t, yaml.MappingNode, labelsNode.Kind)

	// Verify all labels were added
	for k, v := range labels {
		labelValue := lookup(labelsNode, k)
		assert.NotNil(t, labelValue, "label %s should be present", k)
		assert.Equal(t, v, labelValue.Value, "label %s should have value %s", k, v)
	}
}

func TestEditFile_UpdateExistingLabels(t *testing.T) {
	dockerCompose := `services:
  nginx:
    image: nginx:latest
    labels:
      com.lai-panel.managed-by: old-value
      com.lai-panel.owner: old-owner
`
	pipeline := &DockerComposeFileParsePipeline{}
	labels := map[string]string{
		constant.ManagedByLabel: constant.ProjectId,
		constant.OwnerLabel:     constant.ProjectId,
		constant.ServiceLabel:   "nginx-service",
	}

	res, err := pipeline.editFile(dockerCompose, labels)
	assert.NoError(t, err)

	// Parse the result to verify labels were updated
	var y yaml.Node
	err = yaml.Unmarshal([]byte(res), &y)
	assert.NoError(t, err)

	root := &y
	if y.Kind == yaml.DocumentNode && len(y.Content) > 0 {
		root = y.Content[0]
	}

	services := lookup(root, "services")
	nginx := lookup(services, "nginx")
	labelsNode := lookup(nginx, "labels")

	// Verify labels were updated
	assert.Equal(t, constant.ProjectId, lookup(labelsNode, constant.ManagedByLabel).Value)
	assert.Equal(t, constant.ProjectId, lookup(labelsNode, constant.OwnerLabel).Value)
	assert.Equal(t, "nginx-service", lookup(labelsNode, constant.ServiceLabel).Value)
}

func TestEditFile_MultipleServices(t *testing.T) {
	dockerCompose := `services:
  nginx:
    image: nginx:latest
  redis:
    image: redis:latest
    labels:
      existing: label
`
	pipeline := &DockerComposeFileParsePipeline{}
	labels := map[string]string{
		constant.ManagedByLabel: constant.ProjectId,
		constant.OwnerLabel:     constant.ProjectId,
	}

	res, err := pipeline.editFile(dockerCompose, labels)
	assert.NoError(t, err)

	// Parse the result
	var y yaml.Node
	err = yaml.Unmarshal([]byte(res), &y)
	assert.NoError(t, err)

	root := &y
	if y.Kind == yaml.DocumentNode && len(y.Content) > 0 {
		root = y.Content[0]
	}

	services := lookup(root, "services")

	// Verify nginx service has labels
	nginx := lookup(services, "nginx")
	nginxLabels := lookup(nginx, "labels")
	assert.NotNil(t, nginxLabels)
	assert.Equal(t, constant.ProjectId, lookup(nginxLabels, constant.ManagedByLabel).Value)

	// Verify redis service has labels (both new and existing)
	redis := lookup(services, "redis")
	redisLabels := lookup(redis, "labels")
	assert.NotNil(t, redisLabels)
	assert.Equal(t, constant.ProjectId, lookup(redisLabels, constant.ManagedByLabel).Value)
	assert.Equal(t, "label", lookup(redisLabels, "existing").Value)
}

func TestEditFile_InvalidYAML(t *testing.T) {
	invalidYAML := `services:
  nginx:
    image: nginx:latest
    ports:
      - 80:80
    invalid: [unclosed
`
	pipeline := &DockerComposeFileParsePipeline{}
	labels := map[string]string{
		constant.ManagedByLabel: constant.ProjectId,
	}

	res, err := pipeline.editFile(invalidYAML, labels)
	// editFile should return the original file and error
	assert.Error(t, err)
	assert.Equal(t, invalidYAML, res)
}

func TestEditFile_NoServices(t *testing.T) {
	noServices := `version: '3.8'
networks:
  default:
    driver: bridge
`
	pipeline := &DockerComposeFileParsePipeline{}
	labels := map[string]string{
		constant.ManagedByLabel: constant.ProjectId,
	}

	res, err := pipeline.editFile(noServices, labels)
	// Should return error when no services found
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no 'services' map found")
	assert.Equal(t, noServices, res)
}

func TestEditFile_EmptyServices(t *testing.T) {
	emptyServices := `services: {}
`
	pipeline := &DockerComposeFileParsePipeline{}
	labels := map[string]string{
		constant.ManagedByLabel: constant.ProjectId,
	}

	res, err := pipeline.editFile(emptyServices, labels)
	// Should succeed but no services to process
	assert.NoError(t, err)
	assert.NotEmpty(t, res)
}

func TestEditFile_ComplexStructure(t *testing.T) {
	complexYAML := `version: '3.8'
services:
  web:
    image: nginx:latest
    ports:
      - "8080:80"
    environment:
      - ENV=production
    volumes:
      - ./data:/data
  db:
    image: postgres:14
    environment:
      POSTGRES_PASSWORD: secret
    labels:
      custom.label: value
`
	pipeline := &DockerComposeFileParsePipeline{}
	labels := map[string]string{
		constant.ManagedByLabel: constant.ProjectId,
		constant.OwnerLabel:     constant.ProjectId,
		constant.ServiceLabel:   "my-service",
	}

	res, err := pipeline.editFile(complexYAML, labels)
	assert.NoError(t, err)

	// Parse and verify
	var y yaml.Node
	err = yaml.Unmarshal([]byte(res), &y)
	assert.NoError(t, err)

	root := &y
	if y.Kind == yaml.DocumentNode && len(y.Content) > 0 {
		root = y.Content[0]
	}

	services := lookup(root, "services")

	// Verify web service
	web := lookup(services, "web")
	webLabels := lookup(web, "labels")
	assert.NotNil(t, webLabels)
	assert.Equal(t, constant.ProjectId, lookup(webLabels, constant.ManagedByLabel).Value)

	// Verify db service (should preserve existing labels)
	db := lookup(services, "db")
	dbLabels := lookup(db, "labels")
	assert.NotNil(t, dbLabels)
	assert.Equal(t, constant.ProjectId, lookup(dbLabels, constant.ManagedByLabel).Value)
	assert.Equal(t, "value", lookup(dbLabels, "custom.label").Value)
}
