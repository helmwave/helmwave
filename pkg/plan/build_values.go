package plan

import (
	"context"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
)

func (p *Plan) buildValues(ctx context.Context) error {
	log.Info("🔨 Building values...")
	if err := p.ValidateValuesBuild(); err != nil {
		return err
	}

	ch := p.Graph().Run()
	wg := parallel.NewWaitGroup()
	limiter := p.ParallelLimiter(ctx)
	wg.Add(limiter)

	for range limiter {
		go func() {
			for n := range ch {
				wg.ErrChan() <- p.buildReleaseValues(ctx, n.Data)
			}
			wg.Done()
		}()
	}

	return wg.Wait()
}

func (p *Plan) buildReleaseValues(ctx context.Context, rel release.Config) error {
	err := rel.BuildValues(ctx, p.tmpDir, p.templater)
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
	}

	return nil
}
