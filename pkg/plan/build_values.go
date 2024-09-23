package plan

import (
	"context"

	"github.com/helmwave/helmwave/pkg/fileref"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

func (p *Plan) buildValues(ctx context.Context) error {
	log.Info("🔨 Building values...")
	if err := p.ValidateValuesBuild(); err != nil {
		return err
	}

	limit := p.ParallelLimiter(ctx)
	wg, ctx := errgroup.WithContext(ctx)
	wg.SetLimit(limit)

	for _, rel := range p.body.Releases {
		wg.Go(func() error {
			return p.buildReleaseValues(ctx, rel)
		})
	}
	//nolint:wrapcheck
	return wg.Wait()
}

func (p *Plan) buildReleaseValues(ctx context.Context, rel release.Config) error {
	err := rel.BuildValues(ctx, p.tmpDir, p.templater)
	if err != nil {
		log.Errorf("❌ %s values: %v", rel.Uniq(), err)

		return err
	} else {
		vals := helper.SlicesMap(rel.Values(), func(v fileref.Config) string {
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
