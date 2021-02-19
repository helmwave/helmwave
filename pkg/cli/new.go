package cli

import (
	"github.com/zhilyaev/helmwave/pkg/helmwave"
	"helm.sh/helm/v3/pkg/cli"
)

func New() *helmwave.Config {
	return &helmwave.Config{
		Version: "0.8.2",
		Helm:    cli.New(),
		Logger:  &helmwave.Log{},
		Kubedog: &helmwave.KubedogConfig{},
	}
}
