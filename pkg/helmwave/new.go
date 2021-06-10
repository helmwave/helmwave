package helmwave

import (
	"github.com/helmwave/helmwave/pkg/kubedog"
	"helm.sh/helm/v3/pkg/cli"
)

func New() *Config {
	return &Config{
		Version: Version,
		Helm:    cli.New(),
		Logger:  &Log{},
		Kubedog: &kubedog.Config{},
	}
}
