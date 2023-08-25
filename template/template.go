package template

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig/v3"
)

type Template struct {
	values map[string]interface{}
}

func New(sbomReport map[string]interface{}) *Template {
	values := map[string]interface{}{
		"sbomReport": sbomReport,
	}
	return &Template{
		values: values,
	}
}

func (t *Template) Render(tmpl string) (string, error) {
	tpl, err := template.New("").Funcs(sprig.FuncMap()).Delims("[[", "]]").Option("missingkey=error").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, t.values); err != nil {
		return "", err
	}
	return buf.String(), nil
}
