package template

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	log "github.com/sirupsen/logrus"
)

const (
	TemplaterSprig = "sprig"
)

var (
	sprigAliases = map[string]string{ //nolint:gochecknoglobals // can't make these const
		"get":    "sprigGet",
		"hasKey": "sprigHasKey",
	}

	customFuncs = map[string]any{ //nolint:gochecknoglobals // can't make these const
		"toYaml":         ToYaml,
		"fromYaml":       FromYaml,
		"exec":           Exec,
		"setValueAtPath": SetValueAtPath,
		"requiredEnv":    RequiredEnv,
		"required":       Required,
		"readFile":       ReadFile,
		"get":            Get,
		"hasKey":         HasKey,
	}
)

type sprigTemplater struct {
	delimiterLeft, delimiterRight string
}

func (t sprigTemplater) Name() string {
	return TemplaterSprig
}

func (t sprigTemplater) Render(_ context.Context, src string, data any) ([]byte, error) {
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

func (t sprigTemplater) funcMap() template.FuncMap {
	funcMap := template.FuncMap{}

	log.Debug("Loading sprig template functions")
	sprigFuncMap := sprig.TxtFuncMap()
	for orig, alias := range sprigAliases {
		sprigFuncMap[alias] = sprigFuncMap[orig]
	}

	addToMap(funcMap, sprigFuncMap)
	addToMap(funcMap, customFuncs)

	return funcMap
}

func addToMap(dst, src template.FuncMap) {
	for k, v := range src {
		dst[k] = v
	}
}

func (t *sprigTemplater) Delims(left, right string) {
	t.delimiterLeft = left
	t.delimiterRight = right
}
