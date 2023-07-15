package template

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/hairyhenderson/gomplate/v3"
	gomplateData "github.com/hairyhenderson/gomplate/v3/data"
	"github.com/hairyhenderson/gomplate/v3/tmpl"
	log "github.com/sirupsen/logrus"
)

const (
	TemplaterGomplate = "gomplate"
)

type gomplateTemplater struct {
	delimiterLeft, delimiterRight string
}

func (t gomplateTemplater) Name() string {
	return TemplaterGomplate
}

func (t gomplateTemplater) Render(src string, data any) ([]byte, error) {
	tpl := template.New("tpl")
	funcs := t.funcMap(tpl, data)
	tpl, err := tpl.Delims(t.delimiterLeft, t.delimiterRight).Funcs(funcs).Parse(src)
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

func (t gomplateTemplater) funcMap(tpl *template.Template, data any) template.FuncMap {
	funcMap := template.FuncMap{}

	log.Debug("Loading gomplate template functions")
	ctx := context.Background()
	gomplateFuncMap := gomplate.CreateFuncs(ctx, &gomplateData.Data{Ctx: ctx})

	addToMap(funcMap, gomplateFuncMap)
	addToMap(funcMap, customFuncs)

	tp := tmpl.New(tpl, data, tpl.Name())
	funcMap["tmpl"] = func() *tmpl.Template { return tp }
	funcMap["tpl"] = tp.Inline

	return funcMap
}

func (t *gomplateTemplater) Delims(left, right string) {
	t.delimiterLeft = left
	t.delimiterRight = right
}
