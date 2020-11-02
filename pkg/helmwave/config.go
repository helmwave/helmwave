package helmwave

import (
	"github.com/urfave/cli/v2"
	"github.com/zhilyaev/helmwave/pkg/template"
	"github.com/zhilyaev/helmwave/pkg/yml"
	helm "helm.sh/helm/v3/pkg/cli"
)

type Config struct {
	Version  string
	Debug    bool
	Helm     *helm.EnvSettings
	Tags     cli.StringSlice
	Tpl      template.Tpl
	Yml      yml.Config
	Plan     yml.Config
	Parallel bool
}
