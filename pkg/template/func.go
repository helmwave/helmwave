package template

import (
	"context"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/hairyhenderson/gomplate/v3"
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

func FuncMap(gomplateConfig *GomplateConfig) template.FuncMap {
	funcMap := template.FuncMap{}

	log.Debug("Loading sprig template functions")
	sprigFuncMap := sprig.TxtFuncMap()
	for orig, alias := range sprigAliases {
		sprigFuncMap[alias] = sprigFuncMap[orig]
	}
	addToMap(funcMap, sprigFuncMap, "sprig")

	if gomplateConfig.Enabled {
		log.Debug("Loading gomplate template functions")
		gomplateFuncMap := gomplate.CreateFuncs(context.Background(), gomplateConfig.data)
		addToMap(funcMap, gomplateFuncMap, "gomplate")
	}

	addToMap(funcMap, customFuncs, "custom overrides")

	return funcMap
}

func addToMap(dst, src template.FuncMap, name string) {
	for k, v := range src {
		log.Trace("Loading function ", k, " out of ", name)
		dst[k] = v
	}
}
