package template

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"maps"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/helmwave/helmwave/pkg/parallel"
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
	additionalOutputs             []io.Writer
	additionalFuncs               map[string]any
}

func (t sprigTemplater) Name() string {
	return TemplaterSprig
}

func (t sprigTemplater) Render(ctx context.Context, src string, data any) ([]byte, error) {
	funcs := t.funcMap()
	if t.additionalFuncs != nil {
		maps.Copy(funcs, t.additionalFuncs)
	}

	tpl, err := template.New("tpl").Delims(t.delimiterLeft, t.delimiterRight).Funcs(funcs).Parse(src)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	var buf bytes.Buffer
	writers := []io.Writer{&buf}
	if t.additionalOutputs != nil {
		writers = append(writers, t.additionalOutputs...)
	}
	w := io.MultiWriter(writers...)

	wg := parallel.NewWaitGroup()
	wg.Add(1)

	go func() {
		defer wg.Done()
		wg.ErrChan() <- tpl.Execute(w, data)
	}()

	err = wg.WaitWithContext(ctx)
	if err != nil {
		return nil, err
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

func (t *sprigTemplater) AddOutput(w io.Writer) {
	if t.additionalOutputs == nil {
		t.additionalOutputs = []io.Writer{}
	}
	t.additionalOutputs = append(t.additionalOutputs, w)
}

func (t *sprigTemplater) AddFunc(name string, f any) {
	if t.additionalFuncs == nil {
		t.additionalFuncs = map[string]any{}
	}

	t.additionalFuncs[name] = f
}
