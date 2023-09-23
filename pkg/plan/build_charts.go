package plan

import (
	"io/fs"

	"github.com/helmwave/go-fsimpl"
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
)

func (p *Plan) buildCharts(srcFS fs.StatFS, destFS fsimpl.WriteableFS) error {
	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Releases))

	for _, rel := range p.body.Releases {
		go func(wg *parallel.WaitGroup, rel release.Config) {
			defer wg.Done()
			err := rel.DownloadChart(srcFS, destFS, p.tmpDir)
			if err != nil {
				wg.ErrChan() <- err
			}
		}(wg, rel)
	}

	return wg.Wait()
}
