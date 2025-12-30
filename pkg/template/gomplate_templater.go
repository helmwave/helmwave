package template

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"maps"
	"text/template"

	gomplateOldData "github.com/hairyhenderson/gomplate/v3/data"
	gomplateOldFuncs "github.com/hairyhenderson/gomplate/v3/funcs" //nolint:staticcheck
	"github.com/hairyhenderson/gomplate/v4"
	"github.com/hairyhenderson/gomplate/v4/tmpl"
	"github.com/helmwave/helmwave/pkg/parallel"
	log "github.com/sirupsen/logrus"
)

const (
	TemplaterGomplate = "gomplate"
)

type gomplateTemplater struct {
	additionalFuncs               map[string]any
	delimiterLeft, delimiterRight string
	additionalOutputs             []io.Writer
}

func (t gomplateTemplater) Name() string {
	return TemplaterGomplate
}

func (t gomplateTemplater) Render(ctx context.Context, src string, data any) ([]byte, error) {
	tpl := template.New("tpl")
	funcs := t.funcMap(ctx, tpl, data)
	if t.additionalFuncs != nil {
		maps.Copy(funcs, t.additionalFuncs)
	}

	tpl, err := tpl.Delims(t.delimiterLeft, t.delimiterRight).Funcs(funcs).Parse(src)
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

func (t gomplateTemplater) funcMap(ctx context.Context, tpl *template.Template, data any) template.FuncMap {
	funcMap := template.FuncMap{}

	log.Debug("Loading gomplate template functions")
	gomplateFuncMap := gomplate.CreateFuncs(ctx)

	addToMap(funcMap, gomplateOldFuncs.CreateDataFuncs(ctx, &gomplateOldData.Data{Ctx: ctx}))
	addToMap(funcMap, gomplateFuncMap)
	addToMap(funcMap, customFuncs)
	addToMap(funcMap, deprecatedFuncs)

	tp := tmpl.New(tpl, data, tpl.Name())
	funcMap["tmpl"] = func() *tmpl.Template { return tp }
	funcMap["tpl"] = tp.Inline

	return funcMap
}

func (t *gomplateTemplater) Delims(left, right string) {
	t.delimiterLeft = left
	t.delimiterRight = right
}

func (t *gomplateTemplater) AddOutput(w io.Writer) {
	if t.additionalOutputs == nil {
		t.additionalOutputs = []io.Writer{}
	}
	t.additionalOutputs = append(t.additionalOutputs, w)
}

func (t *gomplateTemplater) AddFunc(name string, f any) {
	if t.additionalFuncs == nil {
		t.additionalFuncs = map[string]any{}
	}

	t.additionalFuncs[name] = f
}
