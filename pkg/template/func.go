package template

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig/v3"
	"io/ioutil"
	"os"
	"text/template"
)

func Tpl2yml(tpl string, yml string, debug bool) error {
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
	err = t.Execute(&buf, map[string]interface{}{})
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

	for name, f := range aliased {
		funcMap[name] = f
	}

	return funcMap
}

func Include(filename string) (string, error) {
	return "", nil
}
