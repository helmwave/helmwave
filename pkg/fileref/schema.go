package fileref

import (
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/stoewer/go-strcase"
	"gopkg.in/yaml.v3"
)

//nolint:gocritic
func (v Config) JSONSchema() *jsonschema.Schema {
	r := &jsonschema.Reflector{
		DoNotReference:             true,
		RequiredFromJSONSchemaTags: true,
		KeyNamer:                   strcase.SnakeCase, // for action.ChartPathOptions
	}

	type fileRef Config
	schema := r.Reflect(fileRef(v))
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
func (v *Config) UnmarshalYAML(node *yaml.Node) error {
	type raw Config
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
func (v *Config) MarshalYAML() (any, error) {
	return struct {
		Src string
		Dst string
	}{
		Src: v.Src,
		Dst: v.Dst,
	}, nil
}
