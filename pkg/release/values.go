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

// ErrSkipValues is returned when values cannot be used and are skipped.
var ErrSkipValues = errors.New("values have been skipped")

// ValuesReference is used to match source values file path and temporary.
//
// nolintlint:lll
type ValuesReference struct {
	Src            string `yaml:"src" json:"src" jsonschema:"required,description=Source of values. Can be local path or HTTP URL"` //nolint:lll
	Dst            string `yaml:"dst" json:"dst"`
	DelimiterLeft  string `yaml:"delimiter_left,omitempty" json:"delimiter_left,omitempty" jsonschema:"Set left delimiter for template engine,default={{"`    //nolint:lll
	DelimiterRight string `yaml:"delimiter_right,omitempty" json:"delimiter_right,omitempty" jsonschema:"Set right delimiter for template engine,default=}}"` //nolint:lll
	Strict         bool   `yaml:"strict" json:"strict" jsonschema:"description=Whether to fail if values is not found,default=false"`                         //nolint:lll
	Render         bool   `yaml:"render" json:"render" jsonschema:"description=Whether to use templater to render values,default=true"`                       //nolint:lll
}

func (v ValuesReference) JSONSchema() *jsonschema.Schema {
	r := &jsonschema.Reflector{
		DoNotReference:             true,
		RequiredFromJSONSchemaTags: true,
		KeyNamer:                   strcase.SnakeCase, // for action.ChartPathOptions
	}

	type values ValuesReference
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
	// render by default
	v.Render = true

	type raw ValuesReference
	var err error
	switch node.Kind {
	// single value or reference to another value
	case yaml.ScalarNode, yaml.AliasNode:
		err = node.Decode(&v.Src)
	case yaml.MappingNode:
		err = node.Decode((*raw)(v))
	default:
		err = fmt.Errorf("unknown format")
	}

	if err != nil {
		return fmt.Errorf("failed to decode values reference %q from YAML: %w", node.Value, err)
	}

	return nil
}

// MarshalYAML is used to implement Marshaler interface of gopkg.in/yaml.v3.
func (v ValuesReference) MarshalYAML() (any, error) {
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
	for _, v := range values {
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
// Returns ErrSkipValues if values cannot be downloaded or doesn't exist in local FS.
func (v *ValuesReference) SetViaRelease(rel Config, dir, templater string) error {
	if !v.Render {
		templater = "copy"
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
		err = template.Tpl2yml(v.Dst, v.Dst, data, templater, delimOption)
	} else {
		err = template.Tpl2yml(v.Src, v.Dst, data, templater, delimOption)
	}

	if err != nil {
		return fmt.Errorf("failed to render %q values: %w", v.Src, err)
	}

	return nil
}

// nolintlint:nestif // it is still pretty easy to understand
func (v *ValuesReference) fetch(l *log.Entry) error {
	if v.isURL() {
		err := v.Download()
		if err != nil {
			l.WithError(err).Warnf("%q skipping: cant download", v.Src)

			if v.Strict {
				return ErrSkipValues
			}
		}
	} else if !helper.IsExists(v.Src) {
		l.Warn("skipping: local file not found")

		return ErrSkipValues
	}

	return nil
}

func (rel *config) BuildValues(dir, templater string) error {
	for i := len(rel.Values()) - 1; i >= 0; i-- {
		err := rel.Values()[i].SetViaRelease(rel, dir, templater)
		if errors.Is(ErrSkipValues, err) {
			rel.ValuesF = append(rel.ValuesF[:i], rel.ValuesF[i+1:]...)
		} else if err != nil {
			rel.Logger().WithError(err).WithField("values", rel.Values()[i]).Fatal("failed to build values")

			return err
		}
	}

	return nil
}
