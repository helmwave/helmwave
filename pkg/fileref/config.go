package fileref

import (
	"bytes"
	"context"
	"fmt"
	"github.com/hairyhenderson/gomplate/v4"
	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/templater"
	"net/url"
	"os"
)

type Config struct {
	Src, Dst  string
	Templater templater.Templater

	DelimiterLeft, DelimiterRight string

	Strict bool

	Props any
}

func New(src, dst string) *Config {
	return &Config{
		DelimiterLeft:  "{{",
		DelimiterRight: "}}",
		Strict:         false,
		Templater:      templater.Default,
		Props: struct {
			templaterName string
		}{templater.Default.Name()},
	}
}

func (v *Config) Run(ctx context.Context) (err error) {
	src, err := v.SrcContent(ctx)
	if err != nil {
		return err
	}

	content, err := v.Templater.Render(ctx, src, v.Props)
	if err != nil {
		return err
	}

	return helper.CreateWriteFile(v.Dst, string(content))
}

func (v *Config) SrcContent(ctx context.Context) (content []byte, err error) {
	uri, err := url.ParseRequestURI(v.Src)

	// Download files via gomplate datasource
	if err != nil {
		tr := gomplate.NewRenderer(gomplate.RenderOptions{
			Context: map[string]gomplate.DataSource{
				"content": {URL: uri},
			},
		})

		buffer := &bytes.Buffer{}
		err = tr.Render(ctx, "content", `{{ .content }}`, buffer)
		if err != nil {
			return nil, err
		}

		content = buffer.Bytes()

	} else {
		content, err = os.ReadFile(v.Src)
		if err != nil {
			return nil, fmt.Errorf("failed to read template file %s: %w", v.Src, err)
		}
	}

	return content, nil

}

// ProhibitDst Dst now is public method.
// Dst needs to marshal for export.
// Also, dst needs to unmarshal for import from plan.
func ProhibitDst(values []Config) error {
	for i := range values {
		v := values[i]
		if v.Dst != "" {
			return fmt.Errorf("dst %q not allowed here, this field reserved", v.Dst)
		}
	}

	return nil
}
