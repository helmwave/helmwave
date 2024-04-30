package release

import (
	"errors"

	"helm.sh/helm/v3/pkg/storage/driver"
)

func (rel *config) isInstalled() bool {
	client := rel.newHistory()
	client.Max = 1
	_, err := client.Run(rel.Name())
	switch {
	case errors.Is(err, driver.ErrReleaseNotFound):
		return false
	case err == nil:
		return true
	default:
		rel.Logger().WithError(err).Warn("I can't check if release is installed")

		return false
	}
}
