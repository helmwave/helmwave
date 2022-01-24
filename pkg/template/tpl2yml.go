package template

import (
	"fmt"
	"os"

	"github.com/helmwave/helmwave/pkg/helper"
	log "github.com/sirupsen/logrus"
)

// Templater is interface for using different template function groups.
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
		return nil, fmt.Errorf("templater %s is not registered", name)
	}
}

// Tpl2yml renders 'tpl' file to 'yml' file as go template.
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
		return fmt.Errorf("failed to read template file %s: %w", tpl, err)
	}

	templater, err := getTemplater(templaterName)
	if err != nil {
		return err
	}
	log.WithField("template engine", templater.Name()).Debug("Loaded template engine")

	d, err := templater.Render(string(src), data)
	if err != nil {
		return err //nolint:wrapcheck // we control the interface
	}

	log.Trace(yml, " contents\n", string(d))

	f, err := helper.CreateFile(yml)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", yml, err)
	}

	_, err = f.Write(d)
	if err != nil {
		return fmt.Errorf("failed to write to destination file %s: %w", yml, err)
	}

	err = f.Close()
	if err != nil {
		return fmt.Errorf("failed to close destination file %s: %w", yml, err)
	}

	return nil
}
