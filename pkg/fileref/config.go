package fileref

import (
	"fmt"
	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/invopop/jsonschema"
	"github.com/stoewer/go-strcase"
	"gopkg.in/yaml.v3"
)

type File struct {
	Src            string `yaml:"src" json:"src" jsonschema:"required,description=Source of values. Can be local path or HTTP URL"`
	Dst            string `yaml:"dst" json:"dst" jsonschema:"readOnly"`
	DelimiterLeft  string `yaml:"delimiter_left,omitempty" json:"delimiter_left,omitempty"  jsonschema:"Set left delimiter for template engine,default={{"`
	DelimiterRight string `yaml:"delimiter_right,omitempty" json:"delimiter_right,omitempty" jsonschema:"Set right delimiter for template engine,default=}}"`
	Renderer       string `yaml:"renderer" json:"renderer" jsonschema:"description=How to render the file,enum=sprig,enum=gomplate,enum=copy,enum=sops"`
	Strict         bool   `yaml:"strict" json:"strict" jsonschema:"description=Whether to fail if values is not found,default=false"`
}

//nolint:gocritic
func (v File) JSONSchema() *jsonschema.Schema {
	r := &jsonschema.Reflector{
		DoNotReference:             true,
		RequiredFromJSONSchemaTags: true,
		KeyNamer:                   strcase.SnakeCase, // for action.ChartPathOptions
	}

	type values File
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
func (v *File) UnmarshalYAML(node *yaml.Node) error {
	type raw File
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
