package plan

import (
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
)

// Rollback rollbacks helm release.
func (p *Plan) Rollback(version int) error {
	// Run hooks
	p.body.Lifecycle.PreRolling()
	defer p.body.Lifecycle.PostRolling()

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

	return nil
}
