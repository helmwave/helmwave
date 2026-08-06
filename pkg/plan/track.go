package plan

import (
	"fmt"

	"github.com/fluxcd/cli-utils/pkg/object"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/tracker"
)

type manifestGetter func(release.Config) (string, error)

func (p *Plan) trackerObjects(cfg *tracker.Config, manifests manifestGetter) (object.ObjMetadataSet, string, error) {
	foundContexts := make(map[string]bool)
	var kubecontext string
	var ids object.ObjMetadataSet

	for _, rel := range p.body.Releases {
		kubecontext = rel.KubeContext()
		foundContexts[kubecontext] = true

		l := rel.Logger()
		if !rel.HelmWait() {
			l.Warn("wait is disabled so tracking may stop before resources become ready")
		}

		m, err := manifests(rel)
		if err != nil {
			return nil, "", fmt.Errorf("cannot get manifests for release: %w", err)
		}

		relIDs := tracker.ObjectsFromManifest(m, rel.Namespace(), cfg.TrackGeneric)
		l.WithFields(map[string]any{
			"resources": len(relIDs),
			"release":   rel.Uniq(),
		}).Trace("resources to track")

		ids = append(ids, relIDs...)
	}

	if len(foundContexts) > 1 {
		return nil, "", ErrMultipleKubecontexts
	}

	return ids, kubecontext, nil
}

func (p *Plan) trackerSyncObjects(cfg *tracker.Config) (object.ObjMetadataSet, string, error) {
	return p.trackerObjects(cfg, p.trackerSyncManifest)
}

func (p *Plan) trackerSyncManifest(rel release.Config) (string, error) {
	return p.manifests[rel.Uniq()], nil
}
