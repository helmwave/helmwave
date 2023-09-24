package template

import (
	"bytes"
	"fmt"
	"io/fs"

	"github.com/helmwave/go-fsimpl"
	"github.com/helmwave/helmwave/pkg/helper"
	log "github.com/sirupsen/logrus"
)

// TemplaterOptions is a function that changes templater options.
type TemplaterOptions func(Templater)

// Templater is interface for using different template function groups.
type Templater interface {
	Name() string
	Delims(string, string)
	Render(string, any) ([]byte, error)
}

func getTemplater(name string) (Templater, error) {
	switch name {
	case TemplaterGomplate:
		return &gomplateTemplater{}, nil
	case TemplaterSprig:
		return &sprigTemplater{}, nil
	case TemplaterNone:
		return noTemplater{}, nil
	case TemplaterSOPS:
		return sopsTemplater{}, nil
	default:
		return nil, fmt.Errorf("templater %s is not registered", name)
	}
}

// Tpl2yml renders 'tpl' file to 'yml' file as go template.
func Tpl2yml(
	tplFS fs.FS,
	ymlFSUntyped fs.FS,
	tpl, yml string,
	data any,
	templaterName string,
	opts ...TemplaterOptions,
) error {
	ymlFS, ok := ymlFSUntyped.(fsimpl.WriteableFS)
	if !ok {
		return ErrInvalidFilesystem
	}

	log.WithFields(log.Fields{
		"from": tpl,
		"to":   yml,
	}).Trace("Render yml file")

	if data == nil {
		data = map[string]any{}
	}

	srcFile, err := tplFS.Open(tpl)
	if err != nil {
		return fmt.Errorf("failed to read template file %s: %w", tpl, err)
	}
	defer func() {
		err := srcFile.Close()
		if err != nil {
			log.WithError(err).WithField("file", tpl).Error("failed to close file")
		}
	}()

	buf := &bytes.Buffer{}
	_, err = buf.ReadFrom(srcFile)
	if err != nil {
		log.WithError(err).WithField("file", tpl).Error("failed to read file")
	}
	src := buf.String()

	templater, err := getTemplater(templaterName)
	if err != nil {
		return err
	}
	log.WithField("template engine", templater.Name()).Debug("Loaded template engine")

	for _, opt := range opts {
		opt(templater)
	}

	d, err := templater.Render(src, data)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	log.Trace(yml, " contents\n", string(d))

	f, err := helper.CreateFile(ymlFS, yml)
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
