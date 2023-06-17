package plan

import (
	"github.com/helmwave/helmwave/pkg/hooks"
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
)

// Rollback rollbacks helm release.
func (p *Plan) Rollback(version int) error {
	if len(p.body.Hooks.PreRollback) != 0 {
		log.Info("🩼 Running pre-rollback hooks...")
		hooks.Run(p.body.Hooks.PreRollback)
	}

	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Releases))

	for i := range p.body.Releases {
		go func(wg *parallel.WaitGroup, rel release.Config) {
			defer wg.Done()
			err := rel.Rollback(version)
			if err != nil {
				rel.Logger().WithError(err).Error("❌ rollback")
				wg.ErrChan() <- err
			} else {
				rel.Logger().Info("✅ rollback!")
			}
		}(wg, p.body.Releases[i])
	}

	err := wg.Wait()
	if err != nil {
		return err
	}

	if len(p.body.Hooks.PostRollback) != 0 {
		log.Info("🩼 Running post-rollback hooks...")
		hooks.Run(p.body.Hooks.PostRollback)
	}

	return nil
}
