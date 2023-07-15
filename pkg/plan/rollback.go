package plan

import (
	"context"

	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
)

// Rollback rollbacks helm release.
func (p *Plan) Rollback(ctx context.Context, version int) error {
	// Run hooks
	err := p.body.Lifecycle.RunPreRollback(ctx)
	if err != nil {
		return err
	}

	defer func() {
		err := p.body.Lifecycle.RunPostRollback(ctx)
		if err != nil {
			log.Errorf("got an error from postrollback hooks: %v", err)
		}
	}()

	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Releases))

	for i := range p.body.Releases {
		go func(wg *parallel.WaitGroup, rel release.Config) {
			defer wg.Done()
			err := rel.Rollback(ctx, version)
			if err != nil {
				rel.Logger().WithError(err).Error("❌ rollback")
				wg.ErrChan() <- err
			} else {
				rel.Logger().Info("✅ rollback!")
			}
		}(wg, p.body.Releases[i])
	}

	err = wg.Wait()
	if err != nil {
		return err
	}

	return nil
}
