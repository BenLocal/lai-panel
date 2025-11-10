package tmpl

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestParseDockerCompose(t *testing.T) {
	env := map[string]string{
		"PORT": "8080",
	}
	dockerCompose, err := ParseDockerCompose("test", "{{.PORT}}", env)
	if err != nil {
		t.Fatalf("failed to parse docker compose: %v", err)
	}
	assert.Equal(t, dockerCompose, "8080")
}

func TestParseDockerComposeWithEnv(t *testing.T) {
	env := map[string]string{
		"PORT_HTTP":  "8080",
		"PORT_HTTPS": "8443",
	}
	dockerCompose, err := ParseDockerCompose("test", `
		services:
			web:
				ports:
				{{- if .PORT_HTTP }}
					- "{{.PORT_HTTP}}:8080"
				{{- end }}
				{{- if .PORT_HTTPS }}
					- "{{.PORT_HTTPS}}:8443"
				{{- end }}
	`, env)
	if err != nil {
		t.Fatalf("failed to parse docker compose: %v", err)
	}
	assert.Equal(t, dockerCompose, `
		services:
			web:
				ports:
					- "8080:8080"
					- "8443:8443"
	`)
}
