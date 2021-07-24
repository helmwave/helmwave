package plan

import (
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
)

func (p *Plan) buildValues(dir string) error {
	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Releases))

	for _, rel := range p.body.Releases {
		go func(wg *parallel.WaitGroup, rel *release.Config) {
			defer wg.Done()
			err := rel.BuildValues(dir)
			if err != nil {
				log.Fatalf("❌ %s values: %v", rel.Uniq(), err)
			} else {
				log.Infof("✅ %s values count %d", rel.Uniq(), len(rel.Values))
			}
		}(wg, rel)
	}

	return wg.Wait()
}
