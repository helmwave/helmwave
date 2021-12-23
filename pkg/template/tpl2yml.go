package template

import (
	"errors"
	"os"

	"github.com/helmwave/helmwave/pkg/helper"
	log "github.com/sirupsen/logrus"
)

type Templater interface {
	Name() string
	Render(string, interface{}) ([]byte, error)
}

func getTemplater(name string) (Templater, error) { //nolint:ireturn
	switch name {
	case gomplateTemplater{}.Name():
		return gomplateTemplater{}, nil
	case sprigTemplater{}.Name():
		return sprigTemplater{}, nil
	default:
		return nil, errors.New("Templater not found")
	}
}

func Tpl2yml(tpl, yml string, data interface{}, templaterName string) error {
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

	templater, err := getTemplater(templaterName)
	if err != nil {
		return err
	}
	log.WithField("template engine", templater.Name()).Debug("Loaded template engine")

	d, err := templater.Render(string(src), data)
	if err != nil {
		return err
	}

	log.Trace(yml, " contents\n", d)

	f, err := helper.CreateFile(yml)
	if err != nil {
		return err
	}

	_, err = f.Write(d)
	if err != nil {
		return err
	}

	return f.Close()
}
