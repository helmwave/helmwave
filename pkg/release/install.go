package release

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/storage/driver"
)

func (rel *Config) isInstalled() bool {
	client := action.NewHistory(rel.Cfg())
	client.Max = 1
	_, err := client.Run(rel.Name)
	switch {
	case errors.Is(err, driver.ErrReleaseNotFound):
		return false
	case err == nil:
		return true
	default:
		log.Fatal(err)
		return false
	}
}
