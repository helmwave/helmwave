package plan

import (
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
)

func (p *Plan) Destroy() error {
	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Releases))

	for i := range p.body.Releases {
		go func(wg *parallel.WaitGroup, rel *release.Config) {
			defer wg.Done()
			_, err := rel.Uninstall()
			if err != nil {
				log.Errorf("❌ %s: %v", rel.Uniq(), err)
				wg.ErrChan() <- err
			} else {
				log.Infof("✅ %s uninstalled!", rel.Uniq())
			}
		}(wg, p.body.Releases[i])
	}

	return wg.Wait()
}
