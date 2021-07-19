package template

import (
	"text/template"

	"github.com/Masterminds/sprig/v3"
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
