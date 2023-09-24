package plan

import (
	"io/fs"

	"github.com/helmwave/go-fsimpl"
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/sirupsen/logrus"
)

func (p *Plan) exportValues(srcFS fs.FS, plandirFS fsimpl.WriteableFS) error {
	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Releases))

	for _, rel := range p.body.Releases {
		go func(wg *parallel.WaitGroup, rel release.Config) {
			defer wg.Done()
			err := rel.ExportValues(srcFS, plandirFS, p.templater)
			if err != nil {
				logrus.Errorf("❌ %s values: %v", rel.Uniq(), err)
				wg.ErrChan() <- err
			} else {
				var vals []string
				for i := range rel.Values() {
					vals = append(vals, rel.Values()[i].Src)
				}

				if len(vals) == 0 {
					rel.Logger().Info("🔨 no values provided")
				} else {
					rel.Logger().WithField("values", vals).Infof("✅ found %d values count", len(vals))
				}
			}
		}(wg, rel)
	}

	return wg.Wait()
}
