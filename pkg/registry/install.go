package registry

import (
	"github.com/helmwave/helmwave/pkg/helper"
	"helm.sh/helm/v3/pkg/registry"
)

func (reg *config) Install() error {
	return helper.HelmRegistryClient.Login( //nolint:wrapcheck
		reg.Host(),
		registry.LoginOptBasicAuth(reg.Username, reg.Password),
		registry.LoginOptInsecure(reg.Insecure),
	)
}

// IndexOfHost searches registry in slice of registries by host. Returns offset.
func IndexOfHost(a []Config, host string) (i int, found bool) {
	for i, r := range a {
		if host == r.Host() {
			return i, true
		}
	}

	return i, false
}
