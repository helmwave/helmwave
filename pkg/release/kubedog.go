package release

import (
	"io/ioutil"

	"github.com/helmwave/helmwave/pkg/kubedog"
	log "github.com/sirupsen/logrus"
	"github.com/werf/kubedog/pkg/trackers/rollout/multitrack"
)

func MakeMapSpecs(releases []*Config, manifestPath string) (map[string]*multitrack.MultitrackSpecs, error) {
	mapSpecs := make(map[string]*multitrack.MultitrackSpecs)

	for _, rel := range releases {
		// Todo mv to "copy"
		rel.Options.DryRun = false
		src, err := ioutil.ReadFile(manifestPath + rel.UniqName() + ".yml")
		if err != nil {
			return nil, err
		}

		manifest := kubedog.MakeManifest(src)
		relSpecs, err := kubedog.MakeSpecs(manifest, rel.Options.Namespace)
		if err != nil {
			return nil, err
		}

		log.WithFields(log.Fields{
			"Deployments":  len(relSpecs.Deployments),
			"Jobs":         len(relSpecs.Jobs),
			"DaemonSets":   len(relSpecs.DaemonSets),
			"StatefulSets": len(relSpecs.StatefulSets),
		}).Trace("Specs count of ", rel.UniqName())

		nsSpec, found := mapSpecs[rel.Options.Namespace]
		if found {
			// Merge
			nsSpec.DaemonSets = append(nsSpec.DaemonSets, relSpecs.DaemonSets...)
			nsSpec.Deployments = append(nsSpec.Deployments, relSpecs.Deployments...)
			nsSpec.StatefulSets = append(nsSpec.StatefulSets, relSpecs.StatefulSets...)
			nsSpec.Jobs = append(nsSpec.Jobs, relSpecs.Jobs...)
			mapSpecs[rel.Options.Namespace] = nsSpec
		} else {
			mapSpecs[rel.Options.Namespace] = relSpecs
		}
	}

	return mapSpecs, nil
}
