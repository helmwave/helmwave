package plan

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
)

func (p *Plan) Status(names []string) error {
	return status(p.body.Releases, names)
}

func status(all []*release.Config, names []string) error {
	r := all

	if len(names) > 0 {
		sort.Strings(names)
		r = make([]*release.Config, 0, len(all))

		for _, rel := range all {
			if helper.Contains(string(rel.Uniq()), names) {
				r = append(r, rel)
			}
		}
	}

	for _, rel := range r {
		s, err := rel.Status()
		if err != nil {
			log.Errorf("Failed to get s of %s: %v", rel.Uniq(), err)
			continue
		}

		labels, _ := json.Marshal(s.Labels)
		values, _ := json.Marshal(s.Config)

		log.WithFields(log.Fields{
			"name":          s.Name,
			"namespace":     s.Namespace,
			"chart":         fmt.Sprintf("%s-%s", s.Chart.Name(), s.Chart.Metadata.Version),
			"last deployed": s.Info.LastDeployed,
			"s":             s.Info.Status,
			"revision":      s.Version,
		}).Infof("General s of %s", rel.Uniq())

		log.WithFields(log.Fields{
			"notes":         s.Info.Notes,
			"labels":        string(labels),
			"chart sources": s.Chart.Metadata.Sources,
			"values":        string(values),
		}).Debugf("Debug s of %s", rel.Uniq())

		log.WithFields(log.Fields{
			"hooks":    s.Hooks,
			"manifest": s.Manifest,
		}).Tracef("Superdebug s of %s", rel.Uniq())
	}

	return nil
}
