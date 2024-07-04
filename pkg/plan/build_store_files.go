package plan

import (
	"context"
	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
)

func (p *Plan) buildStoreFiles(ctx context.Context) error {
	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Releases))

	for _, rel := range p.body.Releases {
		go func() {
			defer wg.Done()

			wg.ErrChan() <- p.buildReleaseStoreFiles(ctx, rel)
		}()
	}

	return wg.Wait()
}

func (p *Plan) buildReleaseStoreFiles(ctx context.Context, rel release.Config) error {
	err := rel.BuildStoreFiles(ctx, p.tmpDir, p.templater)
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
			log.WithField("release", rel.Uniq()).WithField("values", vals).Infof("✅ found %d store files count", len(vals))
		}
	}

	return nil
}
