package sprig

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"maps"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/templater/customs"
	log "github.com/sirupsen/logrus"
)

const (
	TemplaterName = "sprig"
)

var sprigAliases = map[string]string{
	"get":    "sprigGet",
	"hasKey": "sprigHasKey",
}

type Templater struct {
	additionalFuncs               map[string]any
	delimiterLeft, delimiterRight string
	additionalOutputs             []io.Writer
}

func (t Templater) Name() string {
	return TemplaterName
}

func (t Templater) Render(ctx context.Context, src []byte, data any) ([]byte, error) {
	funcs := t.funcMap()
	if t.additionalFuncs != nil {
		maps.Copy(funcs, t.additionalFuncs)
	}

	tpl, err := template.New("tpl").Delims(t.delimiterLeft, t.delimiterRight).Funcs(funcs).Parse(string(src))
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

func (t Templater) funcMap() template.FuncMap {
	funcMap := template.FuncMap{}

	log.Debug("Loading sprig template functions")
	sprigFuncMap := sprig.TxtFuncMap()
	for orig, alias := range sprigAliases {
		sprigFuncMap[alias] = sprigFuncMap[orig]
	}

	maps.Copy(funcMap, sprigFuncMap)
	maps.Copy(funcMap, customs.FuncMap)

	return funcMap
}

func (t *Templater) Delims(left, right string) {
	t.delimiterLeft = left
	t.delimiterRight = right
}

func (t *Templater) AddOutput(w io.Writer) {
	if t.additionalOutputs == nil {
		t.additionalOutputs = []io.Writer{}
	}
	t.additionalOutputs = append(t.additionalOutputs, w)
}

func (t *Templater) AddFunc(name string, f any) {
	if t.additionalFuncs == nil {
		t.additionalFuncs = map[string]any{}
	}

	t.additionalFuncs[name] = f
}
