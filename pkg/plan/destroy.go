package plan

import (
	"context"

	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
)

// Destroy destroys all releases that exist in plan.
func (p *Plan) Destroy(ctx context.Context) error {
	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Releases))

	for i := range p.body.Releases {
		go func(ctx context.Context, wg *parallel.WaitGroup, rel release.Config) {
			defer wg.Done()
			_, err := rel.Uninstall(ctx)
			if err != nil {
				log.Errorf("❌ %s: %v", rel.Uniq(), err)
				wg.ErrChan() <- err
			} else {
				log.Infof("✅ %s uninstalled!", rel.Uniq())
			}
		}(ctx, wg, p.body.Releases[i])
	}

	return wg.Wait()
}
