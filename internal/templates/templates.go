package templates

import (
	"bytes"
	"embed"
	"html/template"
	"io"
)

//go:embed data
var Templates embed.FS

type TemplateData interface {
	Patterns() []string
}

func GetTemplateData(templateData TemplateData) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})

	if t, err := template.ParseFS(Templates, templateData.Patterns()...); err != nil {
		return nil, err
	} else {
		if err := t.Execute(buf, templateData); err != nil {
			return nil, err
		}
	}

	return io.ReadAll(buf)
}
