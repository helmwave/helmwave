package plan

import (
	"context"
	"errors"
	"time"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/kubedog"
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
	"github.com/werf/kubedog/pkg/kube"
	"github.com/werf/kubedog/pkg/tracker"
	"github.com/werf/kubedog/pkg/trackers/rollout/multitrack"
)

// Rollback rollbacks helm release.
func (p *Plan) Rollback(ctx context.Context, version int, dog *kubedog.Config) (err error) {
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
		log.Warn("🐶 kubedog is enabled")
		kubedog.FixLog(ctx, dog.LogWidth)
		err = p.rollbackReleasesKubedog(ctx, version, dog)
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

func (p *Plan) rollbackReleasesKubedog(ctx context.Context, version int, kubedogConfig *kubedog.Config) error {
	ctxCancel, cancel := context.WithCancel(ctx)
	defer cancel() // Don't forget!

	specs, kubecontext, err := p.kubedogRollbackSpecs(version, kubedogConfig)
	if err != nil {
		return err
	}

	err = helper.KubeInit(kubecontext)
	if err != nil {
		return err
	}

	opts := multitrack.MultitrackOptions{
		DynamicClient:        kube.DynamicClient,
		DiscoveryClient:      kube.CachedDiscoveryClient,
		Mapper:               kube.Mapper,
		StatusProgressPeriod: kubedogConfig.StatusInterval,
		Options: tracker.Options{
			ParentContext: ctxCancel,
			Timeout:       kubedogConfig.Timeout,
			LogsFromTime:  time.Now(),
		},
	}

	// Run kubedog
	dogroup := parallel.NewWaitGroup()
	dogroup.Add(1)
	go func() {
		defer dogroup.Done()
		log.Trace("Multitrack is starting...")
		dogroup.ErrChan() <- multitrack.Multitrack(kube.Client, specs, opts)
	}()

	// Run helm
	time.Sleep(kubedogConfig.StartDelay)
	err = p.rollbackReleases(ctx, version)
	if err != nil {
		cancel()

		return err
	}

	// Allow kubedog to catch release installed
	time.Sleep(kubedogConfig.StatusInterval)
	cancel() // stop kubedog

	err = dogroup.WaitWithContext(ctx)
	if err != nil && !errors.Is(err, context.Canceled) {
		// Ignore kubedog error
		log.WithError(err).Warn("kubedog caught error while watching resources.")
	}

	return nil
}

func (p *Plan) kubedogRollbackSpecs(
	version int,
	kubedogConfig *kubedog.Config,
) (multitrack.MultitrackSpecs, string, error) {
	return p.kubedogSpecs(kubedogConfig, func(rel release.Config) (string, error) {
		return p.kubedogRollbackManifest(version, rel)
	})
}

func (p *Plan) kubedogRollbackManifest(version int, rel release.Config) (string, error) {
	r, err := rel.Get(version)
	if err != nil {
		return "", err
	}

	return r.Manifest, nil
}
