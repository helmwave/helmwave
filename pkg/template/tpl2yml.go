package template

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/helmwave/helmwave/pkg/helper"
	log "github.com/sirupsen/logrus"
)

// TemplaterOptions is a function that changes templater options.
type TemplaterOptions func(Templater)

// Templater is interface for using different template function groups.
type Templater interface {
	Name() string
	Delims(string, string)
	Render(context.Context, string, any) ([]byte, error)
	AddOutput(io.Writer)
	AddFunc(string, any)
}

func getTemplater(name string) (Templater, error) {
	switch name {
	case TemplaterGomplate:
		return &gomplateTemplater{}, nil
	case TemplaterSprig:
		return &sprigTemplater{}, nil
	case TemplaterNone:
		return &noTemplater{}, nil
	case TemplaterSOPS:
		return &sopsTemplater{}, nil
	default:
		return nil, fmt.Errorf("templater %s is not registered", name)
	}
}

// Tpl2yml renders 'tpl' file to 'yml' file as go template.
func Tpl2yml(ctx context.Context, tpl, yml string, data any, templaterName string, opts ...TemplaterOptions) error {
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

	templater, err := getTemplater(templaterName)
	if err != nil {
		return err
	}
	log.WithField("template engine", templater.Name()).Debug("Loaded template engine")

	for _, opt := range opts {
		opt(templater)
	}

	d, err := templater.Render(ctx, string(src), data)
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
