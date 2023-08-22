package registry

import (
	"github.com/helmwave/helmwave/pkg/log"
	"github.com/invopop/jsonschema"
	"gopkg.in/yaml.v3"
)

// Config is an interface to manage particular helm reg.
type Config interface {
	log.LoggerGetter
	Install() error
	Host() string
	Validate() error
	// Username() string
	// Password() string
}

// Configs type of array Config.
type Configs []Config

// UnmarshalYAML is an unmarshaller for gopkg.in/yaml.v3 to parse YAML into `Config` interface.
func (r *Configs) UnmarshalYAML(node *yaml.Node) error {
	rr := make([]*config, 0)
	err := node.Decode(&rr)
	if err != nil {
		return NewYAMLDecodeError(err)
	}

	*r = make([]Config, len(rr))
	for i := range rr {
		(*r)[i] = rr[i]
	}

	return nil
}

func (Configs) JSONSchema() *jsonschema.Schema {
	r := &jsonschema.Reflector{
		DoNotReference:             true,
		RequiredFromJSONSchemaTags: true,
	}
	var l []*config

	return r.Reflect(&l)
}
