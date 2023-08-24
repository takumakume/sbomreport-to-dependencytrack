package cmd

import (
	"bytes"
	"text/template"
)

type templateEngine struct {
	data map[string]interface{}
}

func newTemplateEngine(sbomReport map[string]interface{}) *templateEngine {
	data := map[string]interface{}{
		"input": sbomReport,
	}
	return &templateEngine{
		data: data,
	}
}

func (t *templateEngine) render(tmpl string) (string, error) {
	tpl, err := template.New("").Delims("[[", "]]").Option("missingkey=error").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, t.data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
