package release

import (
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/storage/driver"
)

func (rel *Config) isInstalled() bool {
	client := action.NewHistory(rel.cfg)
	client.Max = 1
	_, err := client.Run(rel.Name)
	switch err {
	case driver.ErrReleaseNotFound:
		return false
	case nil:
		return true
	default:
		log.Fatal(err)
		return false
	}
}
