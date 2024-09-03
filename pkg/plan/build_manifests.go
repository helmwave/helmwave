package plan

import (
	"context"
	"fmt"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
)

func (p *Plan) buildManifest(ctx context.Context) error {
	log.Info("🔨 Building manifests...")

	wg, ctx := errgroup.WithContext(ctx)
	wg.SetLimit(p.ParallelLimiter(ctx))

	mu := &sync.Mutex{}

	for _, rel := range p.body.Releases {
		wg.Go(
			func() error {
				return p.buildReleaseManifest(ctx, rel, mu)
			})
	}

	//nolint:wrapcheck
	return wg.Wait()
}

func (p *Plan) buildReleaseManifest(ctx context.Context, rel release.Config, mu *sync.Mutex) error {
	l := rel.Logger()

	if err := rel.ChartDepsUpd(); err != nil {
		l.WithError(err).Warn("❌ can't get dependencies")
	}

	r, err := rel.SyncDryRun(ctx, true)
	if err != nil || r == nil {
		l.Errorf("❌ can't get manifests: %v", err)

		return err
	}

	hm := ""
	if !rel.HooksDisabled() {
		for _, h := range r.Hooks {
			hm += fmt.Sprintf("---\n# Source: %s\n%s\n", h.Path, h.Manifest)
		}
	}

	document := r.Manifest
	if len(r.Hooks) > 0 {
		document += hm
	}

	l.Trace(document)

	mu.Lock()
	p.manifests[rel.Uniq()] = document
	mu.Unlock()

	l.Info("✅  manifest done")

	return nil
}
