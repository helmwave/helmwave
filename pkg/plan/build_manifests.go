package plan

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/dependency"
	log "github.com/sirupsen/logrus"
)

func (p *Plan) buildManifest(ctx context.Context) error {
	log.Info("🔨 Building manifests...")

	parallelLimit := p.ParallelLimiter(ctx)

	releasesNodesChan := p.Graph().Run()

	releasesWG := parallel.NewWaitGroup()
	releasesWG.Add(parallelLimit)

	releasesFails := make(map[release.Config]error)

	releasesMutex := &sync.Mutex{}

	for range parallelLimit {
		go p.buildReleaseManifestWorker(ctx, releasesWG, releasesNodesChan, releasesMutex, releasesFails)
	}

	if err := releasesWG.WaitWithContext(ctx); err != nil {
		return err
	}

	return p.ApplyReport(releasesFails, nil)
}

//nolint:dupl
func (p *Plan) buildReleaseManifestWorker(
	ctx context.Context,
	wg *parallel.WaitGroup,
	nodesChan <-chan *dependency.Node[release.Config],
	mu *sync.Mutex,
	fails map[release.Config]error,
) {
	for node := range nodesChan {
		rel := node.Data
		err := p.buildReleaseManifest(ctx, rel, mu)
		if err != nil {
			if rel.AllowFailure() {
				rel.Logger().Errorf("release is allowed to fail, marked as succeeded to dependencies")
				node.SetSucceeded()
			} else {
				node.SetFailed()
			}

			mu.Lock()
			fails[rel] = err
			mu.Unlock()

			wg.ErrChan() <- err
		} else {
			node.SetSucceeded()
		}
	}
	wg.Done()
}

func (p *Plan) buildReleaseManifest(ctx context.Context, rel release.Config, mu *sync.Mutex) (err error) {
	l := rel.Logger()

	if err := rel.ChartDepsUpd(); err != nil {
		l.WithError(err).Warn("❌ can't get dependencies")
	}

	lifecycle := rel.Lifecycle()
	err = lifecycle.RunPreBuild(ctx)
	if err != nil {
		return err
	}
	defer func() {
		lifecycleErr := lifecycle.RunPostBuild(ctx)
		if lifecycleErr != nil && err == nil {
			err = lifecycleErr
		}
	}()

	err = p.buildReleaseValues(ctx, rel, mu)
	if err != nil {
		return err
	}

	r, err := rel.SyncDryRun(ctx, false)
	if err != nil || r == nil {
		l.Errorf("❌ can't get manifests: %v", err)

		return err
	}

	var hm strings.Builder
	if !rel.HooksDisabled() {
		for _, h := range r.Hooks {
			hm.WriteString(fmt.Sprintf("---\n# Source: %s\n%s\n", h.Path, h.Manifest))
		}
	}

	document := r.Manifest
	if len(r.Hooks) > 0 {
		document += hm.String()
	}

	l.Trace(document)

	mu.Lock()
	p.manifests[rel.Uniq()] = document
	mu.Unlock()

	l.Info("✅  manifest done")

	return nil
}
