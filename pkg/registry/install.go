package registry

import (
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
		return NewLoginError(err)
	}

	return nil
}
