package template

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/hairyhenderson/gomplate/v3"
	gomplateData "github.com/hairyhenderson/gomplate/v3/data"
	log "github.com/sirupsen/logrus"
)

type gomplateTemplater struct{}

func (t gomplateTemplater) Name() string {
	return "gomplate"
}

func (t gomplateTemplater) Render(src string, data interface{}) ([]byte, error) {
	funcs := t.funcMap()
	tpl, err := template.New("tpl").Funcs(funcs).Parse(src)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)
	if err != nil {
		return nil, fmt.Errorf("failed to render template: %w", err)
	}

	return buf.Bytes(), nil
}

func (t gomplateTemplater) funcMap() template.FuncMap {
	funcMap := template.FuncMap{}

	log.Debug("Loading gomplate template functions")
	gomplateFuncMap := gomplate.CreateFuncs(context.Background(), &gomplateData.Data{})

	addToMap(funcMap, gomplateFuncMap)
	addToMap(funcMap, customFuncs)

	return funcMap
}
