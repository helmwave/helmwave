package plan

import (
	"context"

	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
)

// Down destroys all releases that exist in a plan.
func (p *Plan) Down(ctx context.Context) (err error) {
	// Run hooks
	err = p.body.Lifecycle.RunPreDown(ctx)
	if err != nil {
		return
	}

	defer func() {
		lifecycleErr := p.body.Lifecycle.RunPostDown(ctx)
		if lifecycleErr != nil {
			log.Errorf("got an error from postdown hooks: %v", lifecycleErr)
			if err == nil {
				err = lifecycleErr
			}
		}
	}()

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

	err = wg.Wait()

	return
}
