package cli

import (
	"github.com/helmwave/helmwave/pkg/helmwave"
	"github.com/helmwave/helmwave/pkg/kubedog"
	"helm.sh/helm/v3/pkg/cli"
)

func New() *helmwave.Config {
	return &helmwave.Config{
		Version: helmwave.Version,
		Helm:    cli.New(),
		Logger:  &helmwave.Log{},
		Kubedog: &kubedog.Config{},
	}
}
