package repo

import (
	"context"
	"fmt"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/log"
	"github.com/invopop/jsonschema"
	"gopkg.in/yaml.v3"
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/repo"
)

// Config is an interface to manage particular helm repository.
type Config interface {
	helper.EqualChecker[Config]
	log.LoggerGetter
	Install(context.Context, *helm.EnvSettings, *repo.File) error
	Name() string
	URL() string
	Validate() error
}

// Configs type of array Config.
type Configs []Config

// UnmarshalYAML is an unmarshaller for gopkg.in/yaml.v3 to parse YAML into `Config` interface.
func (r *Configs) UnmarshalYAML(node *yaml.Node) error {
	rr := make([]*config, 0)
	if err := node.Decode(&rr); err != nil {
		return fmt.Errorf("failed to decode repository config from YAML: %w", err)
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
