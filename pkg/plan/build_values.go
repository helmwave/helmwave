package plan

import (
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
)

func (p *Plan) buildValues(dir string) error {
	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Releases))

	gomplateConfig := &p.body.Template.Gomplate

	for _, rel := range p.body.Releases {
		go func(wg *parallel.WaitGroup, rel *release.Config) {
			defer wg.Done()
			err := rel.BuildValues(dir, gomplateConfig)
			if err != nil {
				log.Errorf("❌ %s values: %v", rel.Uniq(), err)
				wg.ErrChan() <- err
			} else {
				var vals []string
				for i := range rel.Values {
					vals = append(vals, rel.Values[i].Get())
				}
				log.WithField("values", vals).Infof("✅ %s values count %d", rel.Uniq(), len(rel.Values))
			}
		}(wg, rel)
	}

	return wg.Wait()
}
