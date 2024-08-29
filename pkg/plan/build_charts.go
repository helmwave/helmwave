package plan

import (
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
)

func (p *Plan) buildCharts() error {
	log.Info("🔨 Building charts...")

	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Releases))

	for _, rel := range p.body.Releases {
		go func(wg *parallel.WaitGroup, rel release.Config) {
			defer wg.Done()
			err := rel.DownloadChart(p.tmpDir)
			if err != nil {
				wg.ErrChan() <- err
			}
		}(wg, rel)
	}

	return wg.Wait()
}
