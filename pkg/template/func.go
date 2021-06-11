package template

import (
	"bytes"
	"io/ioutil"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/helmwave/helmwave/pkg/helper"
	log "github.com/sirupsen/logrus"
)

var (
	sprigAliases = map[string]string{
		"get":    "sprigGet",
		"hasKey": "sprigHasKey",
	}

	customFuncs = map[string]interface{}{
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

func Tpl2yml(tpl string, yml string, data interface{}) error {
	log.WithFields(log.Fields{
		"from": tpl,
		"to":   yml,
	}).Info("📄 Render file")

	if data == nil {
		data = map[string]interface{}{}
	}

	src, err := ioutil.ReadFile(tpl)
	if err != nil {
		return err
	}

	// Template
	t := template.Must(template.New("tpl").Funcs(FuncMap()).Parse(string(src)))
	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		return err
	}

	log.Debug(yml, " contents\n", buf.String())

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

func FuncMap() template.FuncMap {
	funcMap := sprig.TxtFuncMap()

	for orig, alias := range sprigAliases {
		funcMap[alias] = funcMap[orig]
	}

	for name, f := range customFuncs {
		funcMap[name] = f
	}

	return funcMap
}
