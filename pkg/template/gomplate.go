package template

import (
	"bytes"
	"context"
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
		return nil, err
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)

	return buf.Bytes(), err
}

func (t gomplateTemplater) funcMap() template.FuncMap {
	log.Debug("Loading gomplate template functions")

	return gomplate.CreateFuncs(context.Background(), &gomplateData.Data{})
}
