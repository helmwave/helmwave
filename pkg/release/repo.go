package release

import (
	"strings"

	"helm.sh/helm/v3/pkg/registry"
)

func (rel *config) Repo() string {
	s := rel.Chart().Name
	if registry.IsOCI(s) {
		s = strings.TrimPrefix(s, registry.OCIScheme+"://")
	}

	return strings.Split(s, "/")[0]
}
