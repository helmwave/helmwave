package release

import (
	"errors"
	"fmt"

	"github.com/helmwave/helmwave/pkg/release/uniqname"
)

// ErrSkipValues is returned when values cannot be used and are skipped.
var ErrMissingDependency = errors.New("dependency is missing")

// DependsOnReference is used to store release dependencies.
type DependsOnReference struct {
	Name     string `json:"name"`
	Optional bool   `json:"optional"`
}

// UnmarshalYAML is used to implement InterfaceUnmarshaler interface of github.com/goccy/go-yaml.
func (d *DependsOnReference) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := unmarshal(&d.Name); err != nil {
		type raw DependsOnReference
		if err := unmarshal((*raw)(d)); err != nil {
			return fmt.Errorf("failed to decode depends_on reference from YAML: %w", err)
		}
	}

	return nil
}

func (d *DependsOnReference) Uniq() uniqname.UniqName {
	return uniqname.UniqName(d.Name)
}
