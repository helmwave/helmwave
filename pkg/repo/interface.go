package repo

import (
	"context"
	"fmt"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/log"
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
}

// UnmarshalYAML is an unmarshaller for gopkg.in/yaml.v3 to parse YAML into `Config` interface.
func UnmarshalYAML(node *yaml.Node) ([]Config, error) {
	r := make([]*config, 0)
	if err := node.Decode(&r); err != nil {
		return nil, fmt.Errorf("failed to decode repository Repository from YAML: %w", err)
	}

	res := make([]Config, len(r))
	for i := range r {
		res[i] = r[i]
	}

	return res, nil
}
