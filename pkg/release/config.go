package release

import "helm.sh/helm/v3/pkg/action"

type Config struct {
	Name    string
	Chart   string
	Tags    []string
	Values  []string
	Options action.Upgrade
}
