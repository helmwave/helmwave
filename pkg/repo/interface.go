package repo

import (
	"context"
	"fmt"

	"github.com/helmwave/helmwave/pkg/helper"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/repo"
)

// Config is an interface to manage particular helm repository.
type Config interface {
	helper.EqualChecker[Config]
	Install(context.Context, *helm.EnvSettings, *repo.File) error
	Name() string
	URL() string
	Logger() *log.Entry
}

// UnmarshalYAML is an unmarshaller for gopkg.in/yaml.v3 to parse YAML into `Config` interface.
func UnmarshalYAML(node *yaml.Node) ([]Config, error) {
	r := make([]*config, 0)
	if err := node.Decode(&r); err != nil {
		return nil, fmt.Errorf("failed to decode repository config from YAML: %w", err)
	}

	res := make([]Config, len(r))
	for i := range r {
		res[i] = r[i]
	}

	return res, nil
}
