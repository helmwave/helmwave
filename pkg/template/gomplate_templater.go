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

type gomplateTemplater struct {
	delimiterLeft, delimiterRight string
}

func (t gomplateTemplater) Name() string {
	return "gomplate"
}

//nolint:dupl
func (t gomplateTemplater) Render(src string, data interface{}) ([]byte, error) {
	funcs := t.funcMap()
	tpl, err := template.New("tpl").Delims(t.delimiterLeft, t.delimiterRight).Funcs(funcs).Parse(src)
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
	ctx := context.Background()
	gomplateFuncMap := gomplate.CreateFuncs(ctx, &gomplateData.Data{Ctx: ctx})

	addToMap(funcMap, gomplateFuncMap)
	addToMap(funcMap, customFuncs)

	return funcMap
}

func (t *gomplateTemplater) Delims(left, right string) {
	t.delimiterLeft = left
	t.delimiterRight = right
}
