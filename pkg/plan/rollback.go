package plan

import (
	"context"
	"errors"
	"time"

	"github.com/fluxcd/cli-utils/pkg/object"
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/tracker"
	log "github.com/sirupsen/logrus"
)

// Rollback rollbacks helm release.
func (p *Plan) Rollback(ctx context.Context, version int, dog *tracker.Config) (err error) {
	// Run hooks
	err = p.body.Lifecycle.RunPreRollback(ctx)
	if err != nil {
		return
	}

	defer func() {
		lifecycleErr := p.body.Lifecycle.RunPostRollback(ctx)
		if lifecycleErr != nil {
			log.Errorf("got an error from postrollback hooks: %v", lifecycleErr)
			if err == nil {
				err = lifecycleErr
			}
		}
	}()

	if dog.Enabled {
		log.Warn("🐶 live resource tracking is enabled")
		if err = tracker.SilenceKlog(); err != nil {
			return
		}
		err = p.rollbackReleasesTracked(ctx, version, dog)
	} else {
		err = p.rollbackReleases(ctx, version)
	}

	return
}

func (p *Plan) rollbackReleases(ctx context.Context, version int) error {
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

	return wg.Wait()
}

func (p *Plan) rollbackReleasesTracked(ctx context.Context, version int, cfg *tracker.Config) error {
	ctxCancel, cancel := context.WithCancel(ctx)
	defer cancel() // Don't forget!

	ids, kubecontext, err := p.trackerRollbackObjects(version, cfg)
	if err != nil {
		return err
	}

	// Run the tracker
	dogroup := parallel.NewWaitGroup()
	dogroup.Add(1)
	go func() {
		defer dogroup.Done()
		log.Trace("tracker is starting...")
		dogroup.ErrChan() <- tracker.Track(ctxCancel, cfg, ids, kubecontext)
	}()

	// Run helm
	time.Sleep(cfg.StartDelay)
	err = p.rollbackReleases(ctx, version)
	if err != nil {
		cancel()

		return err
	}

	// Allow the tracker to catch the final states
	time.Sleep(cfg.StatusInterval)
	cancel() // stop the tracker

	err = dogroup.WaitWithContext(ctx)
	if err != nil && !errors.Is(err, context.Canceled) {
		// Tracking is advisory: report and move on
		log.WithError(err).Warn("tracker caught an error while watching resources")
	}

	return nil
}

func (p *Plan) trackerRollbackObjects(version int, cfg *tracker.Config) (object.ObjMetadataSet, string, error) {
	return p.trackerObjects(cfg, func(rel release.Config) (string, error) {
		return p.trackerRollbackManifest(version, rel)
	})
}

func (p *Plan) trackerRollbackManifest(version int, rel release.Config) (string, error) {
	r, err := rel.Get(version)
	if err != nil {
		return "", err
	}

	return r.Manifest, nil
}
