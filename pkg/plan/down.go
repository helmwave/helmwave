package plan

import (
	"context"

	"github.com/helmwave/helmwave/pkg/hooks"
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
)

// Down destroys all releases that exist in a plan.
func (p *Plan) Down(ctx context.Context) error {
	if len(p.body.Hooks.PreDown) != 0 {
		log.Info("🩼 Running pre-down hooks...")
		hooks.Run(p.body.Hooks.PreDown)
	}

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

	err := wg.Wait()
	if err != nil {
		return err
	}

	if len(p.body.Hooks.PostDown) != 0 {
		log.Info("🩼 Running post-down hooks...")
		hooks.Run(p.body.Hooks.PostDown)
	}

	return nil
}
