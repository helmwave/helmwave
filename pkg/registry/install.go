package registry

import (
	"fmt"

	"github.com/helmwave/helmwave/pkg/helper"
	"helm.sh/helm/v3/pkg/registry"
)

func (reg *config) Install() error {
	err := helper.HelmRegistryClient.Login(
		reg.Host(),
		registry.LoginOptBasicAuth(reg.Username, reg.Password),
		registry.LoginOptInsecure(reg.Insecure),
	)
	if err != nil {
		return fmt.Errorf("failed to login in helm registry: %w", err)
	}

	return nil
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
