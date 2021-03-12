package helmwave

import (
	"github.com/helmwave/helmwave/pkg/kubedog"
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/helmwave/helmwave/pkg/yml"
	"github.com/urfave/cli/v2"
	helm "helm.sh/helm/v3/pkg/cli"
)

var Version = "dev"

type Config struct {
	Version  string
	Helm     *helm.EnvSettings
	Tags     cli.StringSlice
	Tpl      template.Tpl
	Yml      yml.Config
	PlanPath string
	Logger   *Log
	Parallel bool
	Kubedog  *kubedog.Config
}
