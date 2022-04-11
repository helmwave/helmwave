package plan

import (
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
)

// Rollback rollbacks helm release.
func (p *Plan) Rollback(version int) error { //nolint:dupl
	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Releases))

	for i := range p.body.Releases {
		go func(wg *parallel.WaitGroup, rel release.Config) {
			defer wg.Done()
			err := rel.Rollback(version)
			if err != nil {
				log.Errorf("❌ %s: %v", rel.Uniq(), err)
				wg.ErrChan() <- err
			} else {
				log.Infof("✅ %s rollback!", rel.Uniq())
			}
		}(wg, p.body.Releases[i])
	}

	return wg.Wait()
}
