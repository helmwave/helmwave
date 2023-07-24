package plan

import (
	"context"
	"errors"
	"fmt"
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
func (p *Plan) Rollback(ctx context.Context, version int, dog *kubedog.Config) error {
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

	if dog.Enabled {
		log.Warn("🐶 kubedog is enabled")
		kubedog.FixLog(dog.LogWidth)
		err = p.rollbackReleasesKubedog(ctx, version, dog)
	} else {
		err = p.rollbackReleases(ctx, version)
	}

	if err != nil {
		return err
	}

	return nil
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
		log.WithError(err).Warn("kubedog has error while watching resources.")
	}

	return nil
}

func (p *Plan) kubedogRollbackSpecs(
	version int,
	kubedogConfig *kubedog.Config,
) (multitrack.MultitrackSpecs, string, error) {
	foundContexts := make(map[string]bool)
	var kubecontext string
	specs := multitrack.MultitrackSpecs{}

	for _, rel := range p.body.Releases {
		kubecontext = rel.KubeContext()
		foundContexts[kubecontext] = true

		l := rel.Logger()
		if !rel.HelmWait() {
			l.Error("wait flag is disabled so kubedog can't correctly track this release")
		}

		r, err := rel.Get(version)
		if err != nil {
			return specs, "", fmt.Errorf("cannot get old manifests for kubedog: %w", err)
		}

		manifest := kubedog.Parse([]byte(r.Manifest))
		spec, err := kubedog.MakeSpecs(manifest, rel.Namespace(), kubedogConfig.TrackGeneric)
		if err != nil {
			return specs, "", fmt.Errorf("kubedog can't parse resources: %w", err)
		}

		l.WithFields(log.Fields{
			"Deployments":  len(spec.Deployments),
			"Jobs":         len(spec.Jobs),
			"DaemonSets":   len(spec.DaemonSets),
			"StatefulSets": len(spec.StatefulSets),
			"Canaries":     len(spec.Canaries),
			"Generics":     len(spec.Generics),
			"release":      rel.Uniq(),
		}).Trace("kubedog track resources")

		specs.Jobs = append(specs.Jobs, spec.Jobs...)
		specs.Deployments = append(specs.Deployments, spec.Deployments...)
		specs.DaemonSets = append(specs.DaemonSets, spec.DaemonSets...)
		specs.StatefulSets = append(specs.StatefulSets, spec.StatefulSets...)
		specs.Canaries = append(specs.Canaries, spec.Canaries...)
		specs.Generics = append(specs.Generics, spec.Generics...)
	}

	if len(foundContexts) > 1 {
		return specs, "", fmt.Errorf("kubedog can't work with releases in multiple kubecontexts")
	}

	return specs, kubecontext, nil
}
