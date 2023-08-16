package release

import (
	"fmt"

	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"gopkg.in/yaml.v3"
)

type DependencyType int

const (
	DependencyRelease DependencyType = iota
	DependencyTag
	DependencyInvalid
)

// DependsOnReference is used to store release dependencies.
//
// nolintlint:lll
type DependsOnReference struct {
	Name     string `yaml:"name" json:"name" jsonschema:"description=Uniqname (or just name if in same namespace) of dependency release"`                    //nolint:lll
	Tag      string `yaml:"tag,omitempty" json:"tag,omitempty" jsonschema:"description=All available releases with the tag will be applied as dependencies"` //nolint:lll
	Optional bool   `yaml:"optional" json:"optional" jsonschema:"description=Whether the dependency is required to be present in plan,default=false"`        //nolint:lll
}

// UnmarshalYAML is used to implement InterfaceUnmarshaler interface of gopkg.in/yaml.v3.
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
		err = ErrUnknownFormat
	}

	if err != nil {
		return fmt.Errorf("failed to decode depends_on reference %q from YAML: %w", node.Value, err)
	}

	return nil
}

func (d *DependsOnReference) Uniq() uniqname.UniqName {
	return uniqname.UniqName(d.Name)
}

func (d *DependsOnReference) Type() DependencyType {
	if d.Name != "" {
		return DependencyRelease
	}

	if d.Tag != "" {
		return DependencyTag
	}

	return DependencyInvalid
}
