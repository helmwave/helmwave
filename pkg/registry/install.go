package registry

import (
	"fmt"

	"github.com/helmwave/helmwave/pkg/helper"
	"helm.sh/helm/v3/pkg/registry"
)

func (c *config) Install() error {
	// Allow public OCI registry #410.
	if c.Username == "" {
		c.Logger().Debugln("Public OCI chart. Skipping helm login.")
		return nil
	}

	err := helper.HelmRegistryClient.Login(
		c.Host(),
		registry.LoginOptBasicAuth(c.Username, c.Password),
		registry.LoginOptInsecure(c.Insecure),
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
