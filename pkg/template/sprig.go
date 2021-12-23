package template

import (
	"bytes"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	log "github.com/sirupsen/logrus"
)

type sprigTemplater struct{}

func (t sprigTemplater) Name() string {
	return "sprig"
}

func (t sprigTemplater) Render(src string, data interface{}) ([]byte, error) {
	funcs := t.funcMap()
	tpl, err := template.New("tpl").Funcs(funcs).Parse(src)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, data)

	return buf.Bytes(), err
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

func addToMap(dst, src template.FuncMap) {
	for k, v := range src {
		dst[k] = v
	}
}
