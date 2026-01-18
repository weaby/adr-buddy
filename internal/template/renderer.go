package template

import (
	"bytes"
	"text/template"

	"github.com/weaby/adr-buddy/internal/model"
)

// Render renders an ADR using the provided template
func Render(adr *model.ADR, tmplStr string) (string, error) {
	tmpl, err := template.New("adr").Parse(tmplStr)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, adr); err != nil {
		return "", err
	}

	return buf.String(), nil
}
