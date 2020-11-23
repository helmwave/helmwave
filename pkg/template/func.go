package template

import (
	"bytes"
	"github.com/Masterminds/sprig/v3"
	log "github.com/sirupsen/logrus"
	"github.com/zhilyaev/helmwave/pkg/helper"
	"io/ioutil"
	"text/template"
)

func Tpl2yml(tpl string, yml string, data interface{}) error {
	log.Infof("📄 Render %s -> %s", tpl, yml)
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

	log.Debugf("Content of %s:\n %+v\n", yml, buf.String())

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
	aliased := template.FuncMap{}

	aliases := map[string]string{
		"get": "sprigGet",
	}

	funcMap := sprig.TxtFuncMap()

	for orig, alias := range aliases {
		aliased[alias] = funcMap[orig]
	}

	funcMap["toYaml"] = ToYaml
	funcMap["fromYaml"] = FromYaml
	funcMap["exec"] = Exec
	funcMap["setValueAtPath"] = SetValueAtPath
	funcMap["requiredEnv"] = RequiredEnv
	funcMap["required"] = Required
	funcMap["readFile"] = ReadFile
	funcMap["get"] = Get

	for name, f := range aliased {
		funcMap[name] = f
	}

	return funcMap
}
