package templater

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/templater/gomplate"
	"github.com/helmwave/helmwave/pkg/templater/no"
	"github.com/helmwave/helmwave/pkg/templater/sops"
	"github.com/helmwave/helmwave/pkg/templater/sprig"
	log "github.com/sirupsen/logrus"
)

// TemplaterOptions is a function that changes templater options.
type TemplaterOptions func(Templater)

var Default = &sprig.Templater{}

// Templater is interface for using different template function groups.
type Templater interface {
	Name() string
	Delims(string, string)
	Render(context.Context, []byte, any) ([]byte, error)
	AddOutput(io.Writer)
	AddFunc(string, any)
}

func GetTemplater(name string) (Templater, error) {
	switch name {
	case gomplate.TemplaterName:
		return &gomplate.Templater{}, nil
	case sprig.TemplaterName:
		return &sprig.Templater{}, nil
	case no.TemplaterName:
		return &no.Templater{}, nil
	case sops.TemplaterName:
		return &sops.Templater{}, nil
	default:
		return nil, fmt.Errorf("templater %s is not registered", name)
	}
}

// Tpl2yml renders 'tpl' file to 'yml' file as go template.
// Deprecated
func Tpl2yml(ctx context.Context, tpl, yml string, data any, templater Templater, opts ...TemplaterOptions) error {
	log.WithFields(log.Fields{
		"from": tpl,
		"to":   yml,
	}).Trace("Render yml file")

	if data == nil {
		data = map[string]any{}
	}

	src, err := os.ReadFile(tpl)
	if err != nil {
		return fmt.Errorf("failed to read template file %s: %w", tpl, err)
	}

	log.WithField("template engine", templater.Name()).Debug("Loaded template engine")

	for _, opt := range opts {
		opt(templater)
	}

	d, err := templater.Render(ctx, src, data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
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

func SetDelimiters(left, right string) TemplaterOptions {
	return func(s Templater) {
		s.Delims(left, right)
	}
}

func CopyOutput(output io.Writer) TemplaterOptions {
	return func(s Templater) {
		s.AddOutput(output)
	}
}

func AddFunc(name string, f any) TemplaterOptions {
	return func(s Templater) {
		s.AddFunc(name, f)
	}
}
