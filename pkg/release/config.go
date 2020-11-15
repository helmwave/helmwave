package release

import "helm.sh/helm/v3/pkg/action"

type Config struct {
	Name    string
	Chart   string
	Tags    []string
	Store   map[string]interface{}
	Values  []string
	Options action.Upgrade
}
