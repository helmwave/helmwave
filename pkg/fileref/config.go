package fileref

import (
	"context"
	_ "crypto/md5" // for crypto.MD5.New to work
	"fmt"
	"strings"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/template"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// Config is used to match source values file path and temporary.
//

type Config struct {
	Src            string `yaml:"src" json:"src" jsonschema:"required,description=Source of values. Can be local path or HTTP URL"`
	Dst            string `yaml:"dst" json:"dst" jsonschema:"readOnly"`
	DelimiterLeft  string `yaml:"delimiter_left,omitempty" json:"delimiter_left,omitempty"  jsonschema:"Set left delimiter for template engine,default={{"`
	DelimiterRight string `yaml:"delimiter_right,omitempty" json:"delimiter_right,omitempty" jsonschema:"Set right delimiter for template engine,default=}}"`
	Renderer       string `yaml:"renderer" json:"renderer" jsonschema:"description=How to render the file,enum=sprig,enum=gomplate,enum=copy,enum=sops"`
	Strict         bool   `yaml:"strict" json:"strict" jsonschema:"description=Whether to fail if values is not found,default=false"`
}

func (v *Config) fetch(ctx context.Context) error {
	if v.isURL() {
		err := v.Download(ctx)
		if err != nil {
			log.WithError(err).Warnf("%q skipping: cant download", v.Src)

			return ErrValuesNotExist
		}
	} else if !helper.IsExists(v.Src) {
		log.Warn("skipping: local file not found")

		return ErrValuesNotExist
	}

	return nil
}

func (v *Config) Set(ctx context.Context, filename, templater string, data any, files *renderedFiles) error {
	if v.Renderer == "" {
		v.Renderer = templater
	}

	v.Dst = filename

	log.Trace("Building values reference")

	err := v.fetch(ctx)
	if err != nil {
		return err
	}

	if v.isURL() {
		err = template.Tpl2yml(ctx, v.Dst, v.Dst, data, v.Renderer, v.tplOpts(files)...)
	} else {
		err = template.Tpl2yml(ctx, v.Src, v.Dst, data, v.Renderer, v.tplOpts(files)...)
	}

	if err != nil {
		return fmt.Errorf("failed to render %q file: %w", v.Src, err)
	}

	return nil
}

func (v *Config) tplOpts(files *renderedFiles) (opts []template.TemplaterOptions) {
	opts = []template.TemplaterOptions{
		template.SetDelimiters(v.DelimiterLeft, v.DelimiterRight),
	}

	if files != nil {
		buf := &strings.Builder{}
		defer files.Add(v.Src, buf)
		opts = append(
			opts,
			template.CopyOutput(buf),
			template.AddFunc("getValues",
				func(filename string) (any, error) {
					s := files.Get(filename).String()

					var res any
					err := yaml.Unmarshal([]byte(s), &res)

					//nolint:wrapcheck
					return res, err
				},
			))
	}

	return opts
}

func (v *Config) isURL() bool {
	return helper.IsURL(v.Src)
}

// Download downloads values by source URL and places to destination path.
func (v *Config) Download(ctx context.Context) error {
	if err := helper.Download(ctx, v.Dst, v.Src); err != nil {
		return fmt.Errorf("failed to download values %s -> %s: %w", v.Src, v.Dst, err)
	}

	return nil
}

// ProhibitDst Dst now is public method.
// Dst needs to marshal for export.
// Also, dst needs to unmarshal for import from plan.
func ProhibitDst(f []Config) error {
	for i := range f {
		v := f[i]
		if v.Dst != "" {
			return fmt.Errorf("dst %q not allowed here, this field reserved", v.Dst)
		}
	}

	return nil
}
