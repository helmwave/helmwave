package helmwave

import (
	"github.com/helmwave/helmwave/pkg/kubedog"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/urfave/cli/v2"
	helm "helm.sh/helm/v3/pkg/cli"
)

var Version = "dev"

type Config struct {
	Version  string
	Helm     *helm.EnvSettings
	Tags     cli.StringSlice
	Tpl      template.Tpl
	Yml      Yml
	Plandir  string
	Logger   *Log
	Kubedog  *kubedog.Config
	Features *Features
}

type Features struct {
	Kubedog      bool
	Parallel     bool
	OverPlan     bool
	MatchAllTags bool
	PlanDeps     bool
	DependsOn    bool
	ReTpl        bool
}

type Yml struct {
	Project      string
	Version      string
	Repositories []*repo.Config
	Releases     []*release.Config
}
