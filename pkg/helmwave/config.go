package helmwave

import (
	"github.com/helmwave/helmwave/pkg/kubedog"
	helm "helm.sh/helm/v3/pkg/cli"
)

var Version = "dev"

type Config struct {
	Version string
	Helm    *helm.EnvSettings
	Logger  *Log
	Kubedog *kubedog.Config
}
