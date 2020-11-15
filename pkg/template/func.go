package template

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"io/ioutil"
	"os"
	"text/template"
)

func Tpl2yml(tpl string, yml string, data interface{}, debug bool) error {
	if data == nil {
		data = map[string]interface{}{}
	}

	if debug {
		fmt.Println("📄 Render", tpl, "->", yml)
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

	if debug {
		fmt.Printf("%+v\n", buf.String())
	}

	f, err := os.Create(yml)
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
