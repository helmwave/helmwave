package release

import (
	"encoding/json"
	"fmt"
	"github.com/helmwave/helmwave/pkg/helper"
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	"sort"
)

func (rel *Config) Status() (*release.Release, error) {
	cfg, err := helper.ActionCfg(rel.Options.Namespace, helm.New())
	if err != nil {
		return nil, err
	}

	client := action.NewStatus(cfg)
	client.ShowDescription = true

	return client.Run(rel.Name)
}

func Status(allReleases []*Config, releasesNames []string) error {
	r := allReleases

	if len(releasesNames) > 0 {
		sort.Strings(releasesNames)
		r = make([]*Config, 0, len(allReleases))

		for _, rel := range allReleases {
			if helper.Contains(rel.UniqName(), releasesNames) {
				r = append(r, rel)
			}
		}
	}

	for _, rel := range r {
		status, err := rel.Status()
		if err != nil {
			log.Errorf("Failed to get status of %s: %v", rel.UniqName(), err)
			continue
		}

		labels, _ := json.Marshal(status.Labels)
		values, _ := json.Marshal(status.Config)

		log.WithFields(log.Fields{
			"name":          status.Name,
			"namespace":     status.Namespace,
			"chart":         fmt.Sprintf("%s-%s", status.Chart.Name(), status.Chart.Metadata.Version),
			"last deployed": status.Info.LastDeployed,
			"status":        status.Info.Status,
			"revision":      status.Version,
		}).Infof("General status of %s", rel.UniqName())

		log.WithFields(log.Fields{
			"notes":         status.Info.Notes,
			"labels":        string(labels),
			"chart sources": status.Chart.Metadata.Sources,
			"values":        string(values),
		}).Debugf("Debug status of %s", rel.UniqName())

		log.WithFields(log.Fields{
			"hooks":    status.Hooks,
			"manifest": status.Manifest,
		}).Tracef("Superdebug status of %s", rel.UniqName())
	}

	return nil
}
