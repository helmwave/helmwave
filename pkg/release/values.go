package release

import (
	"crypto"
	_ "crypto/md5" // for crypto.MD5.New to work
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/invopop/jsonschema"
	log "github.com/sirupsen/logrus"
	"github.com/stoewer/go-strcase"
	"gopkg.in/yaml.v3"
)

// ValuesReference is used to match source values file path and temporary.
//
//nolint:lll
type ValuesReference struct {
	Src            string `yaml:"src" json:"src" jsonschema:"required,description=Source of values. Can be local path or HTTP URL"`
	Dst            string `yaml:"dst" json:"dst" jsonschema:"readOnly"`
	DelimiterLeft  string `yaml:"delimiter_left,omitempty" json:"delimiter_left,omitempty"  jsonschema:"Set left delimiter for template engine,default={{"`
	DelimiterRight string `yaml:"delimiter_right,omitempty" json:"delimiter_right,omitempty" jsonschema:"Set right delimiter for template engine,default=}}"`
	Renderer       string `yaml:"renderer" json:"renderer" jsonschema:"description=How to render the file,enum=sprig,enum=gomplate,enum=copy,enum=sops"`
	Strict         bool   `yaml:"strict" json:"strict" jsonschema:"description=Whether to fail if values is not found,default=false"`
}

func (v *ValuesReference) JSONSchema() *jsonschema.Schema {
	r := &jsonschema.Reflector{
		DoNotReference:             true,
		RequiredFromJSONSchemaTags: true,
		KeyNamer:                   strcase.SnakeCase, // for action.ChartPathOptions
	}

	type values *ValuesReference
	schema := r.Reflect(values(v))
	schema.OneOf = []*jsonschema.Schema{
		{
			Type: "string",
		},
		{
			Type: "object",
		},
	}
	schema.Type = ""

	return schema
}

// UnmarshalYAML flexible config.
func (v *ValuesReference) UnmarshalYAML(node *yaml.Node) error {
	type raw ValuesReference
	var err error
	switch node.Kind {
	// single value or reference to another value
	case yaml.ScalarNode, yaml.AliasNode:
		err = node.Decode(&v.Src)
	case yaml.MappingNode:
		err = node.Decode((*raw)(v))
	default:
		err = ErrUnknownFormat
	}

	if err != nil {
		return fmt.Errorf("failed to decode values reference %q from YAML: %w", node.Value, err)
	}

	return nil
}

// MarshalYAML is used to implement Marshaler interface of gopkg.in/yaml.v3.
func (v *ValuesReference) MarshalYAML() (any, error) {
	return struct {
		Src string
		Dst string
	}{
		Src: v.Src,
		Dst: v.Dst,
	}, nil
}

func (v *ValuesReference) isURL() bool {
	return helper.IsURL(v.Src)
}

// Download downloads values by source URL and places to destination path.
func (v *ValuesReference) Download() error {
	if err := helper.Download(v.Dst, v.Src); err != nil {
		return fmt.Errorf("failed to download values %s -> %s: %w", v.Src, v.Dst, err)
	}

	return nil
}

// Get returns destination path of values.
// func (v *ValuesReference) Get() string {
//	return v.Dst
// }

// SetUniq generates unique file path based on provided base directory, release uniqname and sha1 of source path.
func (v *ValuesReference) SetUniq(dir string, name uniqname.UniqName) *ValuesReference {
	h := crypto.MD5.New()
	h.Write([]byte(v.Src))
	hash := h.Sum(nil)
	s := hex.EncodeToString(hash)

	v.Dst = filepath.Join(dir, "values", name.String(), s+".yml")

	return v
}

// ProhibitDst Dst now is public method.
// Dst needs to marshal for export.
// Also, dst needs to unmarshal for import from plan.
func ProhibitDst(values []ValuesReference) error {
	for i := range values {
		v := values[i]
		if v.Dst != "" {
			return fmt.Errorf("dst %q not allowed here, this field reserved", v.Dst)
		}
	}

	return nil
}

// func (v *ValuesReference) Set(Dst string) *ValuesReference {
//	v.Dst = Dst
//	return v
// }

// SetViaRelease downloads and templates values file.
// Returns ErrValuesNotExist if values can't be downloaded or doesn't exist in local FS.
func (v *ValuesReference) SetViaRelease(rel Config, dir, templater string) error {
	if v.Renderer == "" {
		v.Renderer = templater
	}

	v.SetUniq(dir, rel.Uniq())

	l := rel.Logger().WithField("values src", v.Src).WithField("values Dst", v.Dst)

	l.Trace("Building values reference")

	data := struct {
		Release Config
	}{
		Release: rel,
	}

	err := v.fetch(l)
	if err != nil {
		return err
	}

	delimOption := template.SetDelimiters(v.DelimiterLeft, v.DelimiterRight)
	if v.isURL() {
		err = template.Tpl2yml(v.Dst, v.Dst, data, v.Renderer, delimOption)
	} else {
		err = template.Tpl2yml(v.Src, v.Dst, data, v.Renderer, delimOption)
	}

	if err != nil {
		return fmt.Errorf("failed to render %q values: %w", v.Src, err)
	}

	return nil
}

func (v *ValuesReference) fetch(l *log.Entry) error {
	if v.isURL() {
		err := v.Download()
		if err != nil {
			l.WithError(err).Warnf("%q skipping: cant download", v.Src)

			return ErrValuesNotExist
		}
	} else if !helper.IsExists(v.Src) {
		l.Warn("skipping: local file not found")

		return ErrValuesNotExist
	}

	return nil
}

func (rel *config) BuildValues(dir, templater string) error {
	vals := rel.Values()
	for i := len(vals) - 1; i >= 0; i-- {
		v := vals[i]
		err := v.SetViaRelease(rel, dir, templater)
		switch {
		case !v.Strict && errors.Is(ErrValuesNotExist, err):
			rel.Logger().WithError(err).WithField("values", v).Warn("skipping values...")
			rel.ValuesF = append(rel.ValuesF[:i], rel.ValuesF[i+1:]...)
		case err != nil:
			rel.Logger().WithError(err).WithField("values", v).Error("failed to build values")

			return err
		default:
			rel.Values()[i] = v
		}
	}

	return nil
}
