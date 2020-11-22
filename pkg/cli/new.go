package cli

import (
	"github.com/zhilyaev/helmwave/pkg/helmwave"
	"helm.sh/helm/v3/pkg/cli"
)

func New() *helmwave.Config {
	return &helmwave.Config{
		Version: "0.4.0",
		Helm:    cli.New(),
		Logger:  &helmwave.Log{},
		//Logger: &helmwave.Log{Engine: logrus.New()},
	}
}
