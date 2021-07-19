package plan

import (
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
)

func (p *Plan) Rollback() error {
	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Releases))

	for i := range p.body.Releases {
		go func(wg *parallel.WaitGroup, rel *release.Config) {
			defer wg.Done()
			log.Info(rel.Uniq(), " rollback...")
			err := rel.Rollback()
			if err != nil {
				log.Warn(err)
			}
		}(wg, p.body.Releases[i])
	}

	return wg.Wait()
}
