package plan

import (
	"fmt"

	"github.com/helmwave/helmwave/pkg/kubedog"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
	"github.com/werf/kubedog/pkg/trackers/rollout/multitrack"
)

type manifestGetter func(release.Config) (string, error)

func (p *Plan) kubedogSpecs(
	kubedogConfig *kubedog.Config,
	manifests manifestGetter,
) (multitrack.MultitrackSpecs, string, error) {
	foundContexts := make(map[string]bool)
	var kubecontext string
	specs := multitrack.MultitrackSpecs{}

	for _, rel := range p.body.Releases {
		kubecontext = rel.KubeContext()
		foundContexts[kubecontext] = true

		l := rel.Logger()
		if !rel.HelmWait() {
			l.Error("wait flag is disabled so kubedog can't correctly track this release")
		}

		m, err := manifests(rel)
		if err != nil {
			return specs, "", fmt.Errorf("cannot get manifests for release: %w", err)
		}

		manifest := kubedog.Parse([]byte(m))
		spec, err := kubedog.MakeSpecs(manifest, rel.Namespace(), kubedogConfig.TrackGeneric)
		if err != nil {
			return specs, "", fmt.Errorf("kubedog can't parse resources: %w", err)
		}

		l.WithFields(log.Fields{
			"Deployments":  len(spec.Deployments),
			"Jobs":         len(spec.Jobs),
			"DaemonSets":   len(spec.DaemonSets),
			"StatefulSets": len(spec.StatefulSets),
			"Canaries":     len(spec.Canaries),
			"Generics":     len(spec.Generics),
			"release":      rel.Uniq(),
		}).Trace("kubedog track resources")

		specs.Jobs = append(specs.Jobs, spec.Jobs...)
		specs.Deployments = append(specs.Deployments, spec.Deployments...)
		specs.DaemonSets = append(specs.DaemonSets, spec.DaemonSets...)
		specs.StatefulSets = append(specs.StatefulSets, spec.StatefulSets...)
		specs.Canaries = append(specs.Canaries, spec.Canaries...)
		specs.Generics = append(specs.Generics, spec.Generics...)
	}

	if len(foundContexts) > 1 {
		return specs, "", fmt.Errorf("kubedog can't work with releases in multiple kubecontexts")
	}

	return specs, kubecontext, nil
}
