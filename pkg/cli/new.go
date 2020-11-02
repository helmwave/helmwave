package cli

import (
	"github.com/zhilyaev/helmwave/pkg/helmwave"
	"helm.sh/helm/v3/pkg/cli"
)

func New() *helmwave.Config {
	return &helmwave.Config{
		Version: "0.1.4",
		Helm:    cli.New(),
	}
}
