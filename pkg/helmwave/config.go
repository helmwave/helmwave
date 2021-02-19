package helmwave

import (
	"github.com/urfave/cli/v2"
	"github.com/zhilyaev/helmwave/pkg/template"
	"github.com/zhilyaev/helmwave/pkg/yml"
	helm "helm.sh/helm/v3/pkg/cli"
	"time"
)

type Config struct {
	Version  string
	Helm     *helm.EnvSettings
	Tags     cli.StringSlice
	Tpl      template.Tpl
	Yml      yml.Config
	PlanPath string
	Logger   *Log
	Parallel bool
	Kubedog  *KubedogConfig
}

type KubedogConfig struct {
	Enabled        bool
	StatusInterval time.Duration
	Timeout        time.Duration
}
