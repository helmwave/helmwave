package monitor

import (
	"context"

	"github.com/helmwave/helmwave/pkg/log"
	"github.com/invopop/jsonschema"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// SubConfig is an interface to manage particular typed monitor.
type SubConfig interface {
	Init(context.Context, *logrus.Entry) error
	Run(context.Context) error
	Validate() error
}

// Config is an interface to manage particular monitor.
type Config interface {
	log.LoggerGetter
	Name() string
	Run(context.Context) error
	Validate() error
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
