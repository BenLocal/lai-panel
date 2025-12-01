package tmpl

import (
	"bytes"
	"text/template"
)

func ParseWithEnv(name string, tmpl string, env map[string]string, funcs ...map[string]interface{}) (string, error) {
	funcMap := template.FuncMap{}
	for _, f := range funcs {
		for k, v := range f {
			funcMap[k] = v
		}
	}
	t := template.New(name).Funcs(funcMap)

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
