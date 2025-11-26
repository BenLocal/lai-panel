package tmpl

import (
	"bytes"
	"text/template"
)

func ParseDockerCompose(name string, tmpl string, env map[string]string) (string, error) {
	t := template.New(name).Funcs(template.FuncMap{})

	t, err := t.Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, env)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func ParseWithEnv(name string, tmpl string, env map[string]string) (string, error) {
	t := template.New(name).Funcs(template.FuncMap{})

	t, err := t.Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, env)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
