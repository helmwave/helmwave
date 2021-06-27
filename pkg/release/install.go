package release

import (
	"github.com/imdario/mergo"
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage/driver"
)

func (rel *Config) install(cfg *action.Configuration, chart *chart.Chart, values map[string]interface{}) (*release.Release, error) {
	client := action.NewInstall(cfg)

	// Merge
	err := mergo.Merge(client, rel.Install)
	if err != nil {
		return nil, err
	}

	return client.Run(chart, values)
}

func (rel *Config) isInstalled (cfg *action.Configuration) bool {
	client := action.NewHistory(cfg)
	client.Max = 1
	_, err := client.Run(rel.ReleaseName)
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

