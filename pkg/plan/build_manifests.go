package plan

import (
	"context"
	"fmt"
	"sync"

	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
)

func (p *Plan) buildManifest(ctx context.Context) error {
	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Releases))

	mu := &sync.Mutex{}

	for _, rel := range p.body.Releases {
		go p.buildReleaseManifest(ctx, wg, rel, mu)
	}

	return wg.Wait()
}

func (p *Plan) buildReleaseManifest(ctx context.Context, wg *parallel.WaitGroup, rel release.Config, mu *sync.Mutex) {
	defer wg.Done()

	l := rel.Logger()

	if err := rel.ChartDepsUpd(); err != nil {
		l.WithError(err).Warn("❌ can't get dependencies")
	}

	r, err := rel.SyncDryRun(ctx, true)
	if err != nil || r == nil {
		l.Errorf("❌ can't get manifests: %v", err)
		wg.ErrChan() <- err

		return
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
}
