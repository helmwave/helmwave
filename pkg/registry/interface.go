package registry

import (
	"fmt"

	"github.com/helmwave/helmwave/pkg/log"
	"github.com/invopop/jsonschema"
)

// Config is an interface to manage particular helm reg.
type Config interface {
	log.LoggerGetter
	Install() error
	Host() string
	// Username() string
	// Password() string
}

// Configs type of array Config.
type Configs []Config

// UnmarshalYAML is an unmarshaller for github.com/goccy/go-yaml to parse YAML into `Config` interface.
func (r *Configs) UnmarshalYAML(unmarshal func(interface{}) error) error {
	rr := make([]*config, 0)
	if err := unmarshal(&rr); err != nil {
		return fmt.Errorf("failed to decode registry config from YAML: %w", err)
	}

	*r = make([]Config, len(rr))
	for i := range rr {
		(*r)[i] = rr[i]
	}

	return nil
}

func (Configs) JSONSchema() *jsonschema.Schema {
	r := &jsonschema.Reflector{DoNotReference: true}
	var l []*config

	return r.Reflect(&l)
}
