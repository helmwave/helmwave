package plan

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
)

func (p *Plan) exportManifests(ctx context.Context, plandirFS ExportFS) error {
	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Releases))

	for _, rel := range p.body.Releases {
		go p.exportReleaseManifest(ctx, wg, rel, plandirFS)
	}

	return wg.Wait()
}

func (p *Plan) exportReleaseManifest(
	ctx context.Context,
	wg *parallel.WaitGroup,
	rel release.Config,
	baseFS ExportFS,
) {
	defer wg.Done()

	l := rel.Logger()

	if err := rel.ChartDepsUpd(baseFS); err != nil {
		l.WithError(err).Warn("❌ can't get dependencies")
	}

	r, err := rel.SyncDryRun(ctx, baseFS)
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
	p.manifests[rel.Uniq()] = document

	l.Info("✅ manifest done")

	m := filepath.Join(Manifest, rel.Uniq().String()+".yml")

	f, err := helper.CreateFile(baseFS, m)
	if err != nil {
		wg.ErrChan() <- err

		return
	}

	_, err = f.Write([]byte(document))
	if err != nil {
		wg.ErrChan() <- fmt.Errorf("failed to write manifest %s: %w", m, err)

		return
	}

	err = f.Close()
	if err != nil {
		wg.ErrChan() <- fmt.Errorf("failed to close manifest %s: %w", m, err)

		return
	}
}
