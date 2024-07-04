package fileref

import (
	"context"
	"errors"
	"fmt"

	"github.com/helmwave/helmwave/pkg/helper"
	log "github.com/sirupsen/logrus"
)

var (
	// ErrNotExist is returned when files can't be used and are skipped.
	ErrNotExist = errors.New("file reference doesn't exist")

	ErrUnknownFormat = errors.New("unknown format")
)

// Config is used to match source values file path and temporary.
//
//nolint:lll
type Config struct {
	Src            string `yaml:"src" json:"src" jsonschema:"required,description=Source of values. Can be local path or HTTP URL"`
	Dst            string `yaml:"dst" json:"dst" jsonschema:"readOnly"`
	DelimiterLeft  string `yaml:"delimiter_left,omitempty" json:"delimiter_left,omitempty"  jsonschema:"Set left delimiter for template engine,default={{"`
	DelimiterRight string `yaml:"delimiter_right,omitempty" json:"delimiter_right,omitempty" jsonschema:"Set right delimiter for template engine,default=}}"`
	Renderer       string `yaml:"renderer" json:"renderer" jsonschema:"description=How to render the file,enum=sprig,enum=gomplate,enum=copy,enum=sops,default=sprig"`
	Strict         bool   `yaml:"strict" json:"strict" jsonschema:"description=Whether to fail if values is not found,default=false"`
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

// Get returns destination path of values.
// func (v *Config) Get() string {
//	return v.Dst
// }

// ProhibitDst Dst now is public method.
// Dst needs to marshal for export.
// Also, dst needs to unmarshal for import from plan.
func ProhibitDst(files []Config) error {
	for i := range files {
		v := files[i]
		if v.Dst != "" {
			return fmt.Errorf("dst %q not allowed here, this field reserved", v.Dst)
		}
	}

	return nil
}

func (v *Config) fetch(ctx context.Context, l *log.Entry) error {
	if v.isURL() {
		err := v.Download(ctx)
		if err != nil {
			l.WithError(err).Warnf("%q skipping: can't download", v.Src)

			return ErrNotExist
		}
	} else if !helper.IsExists(v.Src) {
		l.Warn("skipping: local file not found")
		return ErrNotExist
	}

	return nil
}
