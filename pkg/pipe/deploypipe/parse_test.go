package deploypipe

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// normalizeWhitespace removes all whitespace characters for comparison
func normalizeWhitespace(s string) string {
	// Remove all whitespace characters (spaces, tabs, newlines, etc.)
	var result strings.Builder
	for _, r := range s {
		if !strings.ContainsRune(" \t\n\r", r) {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func TestDockerComposeFileParsePipelineProcess(t *testing.T) {
	dockerCompose := `services:
  nginx:
    image: nginx:latest
    ports:
      - 80:80
`
	pipeline := &DockerComposeFileParsePipeline{}
	res, err := pipeline.editFile(dockerCompose, map[string]string{
		ManagedByLabel: "lai-panel",
		OwnerLabel:     "lai-panel",
	})
	if err != nil {
		t.Fatalf("editFile failed: %v", err)
	}

	expected := `services:
  nginx:
    image: nginx:latest
    ports:
      - 80:80
    labels:
      managed-by: lai-panel
      owner: lai-panel
`
	// Normalize whitespace for comparison
	assert.Equal(t, normalizeWhitespace(res), normalizeWhitespace(expected))
}
