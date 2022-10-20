package release

import (
	"errors"
	"fmt"

	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"gopkg.in/yaml.v3"
)

// ErrSkipValues is returned when values cannot be used and are skipped.
var ErrMissingDependency = errors.New("dependency is missing")

// DependsOnReference is used to store release dependencies.
type DependsOnReference struct {
	Name     string `yaml:"name" json:"name"`
	Optional bool   `yaml:"optional" json:"optional"`
}

// UnmarshalYAML is used to implement Unmarshaler interface of gopkg.in/yaml.v3.
func (d *DependsOnReference) UnmarshalYAML(node *yaml.Node) error {
	type raw DependsOnReference
	var err error
	switch node.Kind {
	// single value or reference to another value
	case yaml.ScalarNode, yaml.AliasNode:
		err = node.Decode(&d.Name)
	case yaml.MappingNode:
		err = node.Decode((*raw)(d))
	default:
		err = fmt.Errorf("unknown format")
	}

	if err != nil {
		return fmt.Errorf("failed to decode depends_on reference %q from YAML: %w", node.Value, err)
	}

	return nil
}

func (d *DependsOnReference) Uniq() uniqname.UniqName {
	return uniqname.UniqName(d.Name)
}
