package plan

import (
	"io/fs"

	"github.com/helmwave/go-fsimpl"
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
)

func (p *Plan) buildValues(srcFS fs.StatFS, destFS fsimpl.WriteableFS) error {
	if err := p.ValidateValuesBuild(); err != nil {
		return err
	}

	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Releases))

	for _, rel := range p.body.Releases {
		go func(wg *parallel.WaitGroup, rel release.Config) {
			defer wg.Done()
			err := rel.BuildValues(srcFS, destFS, p.tmpDir, p.templater)
			if err != nil {
				log.Errorf("❌ %s values: %v", rel.Uniq(), err)
				wg.ErrChan() <- err
			} else {
				var vals []string
				for i := range rel.Values() {
					vals = append(vals, rel.Values()[i].Dst)
				}

				if len(vals) == 0 {
					rel.Logger().Info("🔨 no values provided")
				} else {
					log.WithField("release", rel.Uniq()).WithField("values", vals).Infof("✅ found %d values count", len(vals))
				}
			}
		}(wg, rel)
	}

	return wg.Wait()
}
