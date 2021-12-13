package template

import (
	"bytes"
	"os"
	"text/template"

	"github.com/helmwave/helmwave/pkg/helper"
	log "github.com/sirupsen/logrus"
)

func Tpl2yml(tpl, yml string, data interface{}, gomplateConfig *GomplateConfig) error {
	log.WithFields(log.Fields{
		"from": tpl,
		"to":   yml,
	}).Trace("Render yml file")

	if data == nil {
		data = map[string]interface{}{}
	}

	src, err := os.ReadFile(tpl)
	if err != nil {
		return err
	}

	// Template
	funcs := FuncMap(gomplateConfig)
	t, err := template.New("tpl").Funcs(funcs).Parse(string(src))
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		return err
	}

	log.Trace(yml, " contents\n", buf.String())

	f, err := helper.CreateFile(yml)
	if err != nil {
		return err
	}

	_, err = f.WriteString(buf.String())
	if err != nil {
		return err
	}

	return f.Close()
}
