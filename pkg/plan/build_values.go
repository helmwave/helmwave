package plan

import (
	"context"
	"sync"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/dependency"
	log "github.com/sirupsen/logrus"
)

func (p *Plan) buildValues(ctx context.Context) error {
	log.Info("🔨 Building values...")
	if err := p.ValidateValuesBuild(); err != nil {
		return err
	}

	parallelLimit := p.ParallelLimiter(ctx)

	releasesNodesChan := p.Graph().Run()

	releasesWG := parallel.NewWaitGroup()
	releasesWG.Add(parallelLimit)

	releasesFails := make(map[release.Config]error)

	releasesMutex := &sync.Mutex{}

	for range parallelLimit {
		go p.buildReleaseValuesWorker(ctx, releasesWG, releasesNodesChan, releasesMutex, releasesFails)
	}

	if err := releasesWG.WaitWithContext(ctx); err != nil {
		return err
	}

	return p.ApplyReport(releasesFails, nil)
}

//nolint:dupl
func (p *Plan) buildReleaseValuesWorker(
	ctx context.Context,
	wg *parallel.WaitGroup,
	nodesChan <-chan *dependency.Node[release.Config],
	mu *sync.Mutex,
	fails map[release.Config]error,
) {
	for node := range nodesChan {
		rel := node.Data
		err := p.buildReleaseValues(ctx, rel, mu)
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

func (p *Plan) buildReleaseValues(ctx context.Context, rel release.Config, mu *sync.Mutex) error {
	log.Info("🔨 Building release values...")

	templateFuncs := p.templateFuncs(mu)

	renderedValues, err := rel.BuildValues(ctx, p.tmpDir, p.templater, templateFuncs)
	if err != nil {
		log.Errorf("❌ %s values: %v", rel.Uniq(), err)

		return err
	} else {
		vals := helper.SlicesMap(rel.Values(), func(v release.ValuesReference) string {
			return v.Dst
		})

		if len(vals) == 0 {
			rel.Logger().Info("🔨 no values provided")
		} else {
			log.WithField("release", rel.Uniq()).WithField("values", vals).Infof("✅ found %d values count", len(vals))
		}

		mu.Lock()
		p.values[rel.Uniq()] = renderedValues
		mu.Unlock()
	}

	return nil
}
