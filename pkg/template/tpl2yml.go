package template

import (
	"bytes"
	"io/ioutil"
	"text/template"

	"github.com/helmwave/helmwave/pkg/helper"
	log "github.com/sirupsen/logrus"
)

func Tpl2yml(tpl, yml string, data interface{}) error {
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
