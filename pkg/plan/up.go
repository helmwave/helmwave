package plan

import (
	"context"
	"errors"
	"maps"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/helmwave/helmwave/pkg/monitor"
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/dependency"
	"github.com/helmwave/helmwave/pkg/tracker"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
)

// Up syncs repositories and releases.
func (p *Plan) Up(ctx context.Context, dog *tracker.Config) (err error) {
	// Run hooks
	err = p.body.Lifecycle.RunPreUp(ctx)
	if err != nil {
		return
	}

	defer func() {
		lifecycleErr := p.body.Lifecycle.RunPostUp(ctx)
		if lifecycleErr != nil {
			log.Errorf("got an error from postup hooks: %v", lifecycleErr)
			if err == nil {
				err = lifecycleErr
			}
		}
	}()

	log.Info("🗄 sync repositories...")
	err = p.syncRepositories(ctx)
	if err != nil {
		return
	}

	log.Info("🗄 sync registries...")
	err = p.syncRegistries(ctx)
	if err != nil {
		return
	}

	if len(p.body.Releases) == 0 {
		return
	}

	log.Info("🛥 sync releases...")

	if dog.Enabled {
		log.Warn("🐶 live resource tracking is enabled")
		if err = tracker.SilenceKlog(); err != nil {
			return
		}
		err = p.syncReleasesTracked(ctx, dog)
	} else {
		err = p.syncReleases(ctx)
	}

	return
}

func (p *planBody) generateMonitorsLockMap() map[string]*parallel.WaitGroup {
	res := make(map[string]*parallel.WaitGroup)

	for _, rel := range p.Releases {
		allMons := rel.Monitors()
		for i := range allMons {
			mon := allMons[i]
			if _, ok := res[mon.Name]; !ok {
				res[mon.Name] = parallel.NewWaitGroup()
			}

			res[mon.Name].Add(1)
		}
	}

	return res
}

func (p *Plan) syncReleases(ctx context.Context) (err error) {
	parallelLimit := p.ParallelLimiter(ctx)

	monitorsLockMap := p.body.generateMonitorsLockMap()
	monitorsCtx, monitorsCancel := context.WithCancel(ctx)
	defer monitorsCancel()

	releasesNodesChan := p.Graph().Run()

	releasesWG := parallel.NewWaitGroup()
	releasesWG.Add(parallelLimit)

	monitorsWG := parallel.NewWaitGroup()
	monitorsWG.Add(len(p.body.Monitors))

	releasesFails := make(map[release.Config]error)
	monitorsFails := make(map[monitor.Config]error)

	releasesMutex := &sync.Mutex{}

	for range parallelLimit {
		go p.syncReleasesWorker(ctx, releasesWG, releasesNodesChan, releasesMutex, releasesFails, monitorsLockMap)
	}

	for _, mon := range p.body.Monitors {
		go p.monitorsWorker(monitorsCtx, monitorsWG, mon, monitorsFails, monitorsLockMap)
	}

	if err := releasesWG.WaitWithContext(ctx); err != nil {
		return err
	}

	if err := monitorsWG.WaitWithContext(monitorsCtx); err != nil {
		log.WithError(err).Error("monitors failed, need to take actions")
		p.runMonitorsActions(ctx, monitorsFails)
	}

	return p.ApplyReport(releasesFails, monitorsFails)
}

func (p *Plan) runMonitorsActions(
	ctx context.Context,
	fails map[monitor.Config]error,
) {
	mons := slices.Collect(maps.Keys(fails))

	for _, rel := range p.body.Releases {
		rel.NotifyMonitorsFailed(ctx, mons...)
	}
}

func (p *Plan) syncReleasesWorker(
	ctx context.Context,
	wg *parallel.WaitGroup,
	nodesChan <-chan *dependency.Node[release.Config],
	mu *sync.Mutex,
	fails map[release.Config]error,
	monitorsLockMap map[string]*parallel.WaitGroup,
) {
	for n := range nodesChan {
		p.syncRelease(ctx, wg, n, mu, fails, monitorsLockMap)
	}
	wg.Done()
}

func (p *Plan) syncRelease(
	ctx context.Context,
	wg *parallel.WaitGroup,
	node *dependency.Node[release.Config],
	mu *sync.Mutex,
	fails map[release.Config]error,
	monitorsLockMap map[string]*parallel.WaitGroup,
) {
	rel := node.Data

	l := rel.Logger()

	l.Info("🛥 deploying... ")

	if _, err := rel.Sync(ctx, true); err != nil {
		l.WithError(err).Error("❌ failed to deploy")

		if rel.AllowFailure() {
			l.Errorf("release is allowed to fail, marked as succeeded to dependencies")
			node.SetSucceeded()
		} else {
			node.SetFailed()
		}

		mu.Lock()
		fails[rel] = err
		mu.Unlock()

		wg.ErrChan() <- err
	} else {
		node.SetSucceeded()
		l.Info("✅")

		allMons := rel.Monitors()
		for i := range allMons {
			mon := allMons[i]
			m := monitorsLockMap[mon.Name]
			if m != nil {
				m.Done()
			}
		}
	}
}

func (p *Plan) monitorsWorker(
	ctx context.Context,
	wg *parallel.WaitGroup,
	mon monitor.Config,
	fails map[monitor.Config]error,
	monitorsLockMap map[string]*parallel.WaitGroup,
) {
	defer wg.Done()

	l := mon.Logger()

	lock := monitorsLockMap[mon.Name()]
	if lock == nil {
		l.Error("BUG: monitor lock is empty, skipping monitor")

		return
	}
	err := lock.WaitWithContext(ctx)
	if err != nil {
		l.WithError(err).Error("❌ monitor canceled")
		fails[mon] = err
		wg.ErrChan() <- err
	}

	err = mon.Run(ctx)
	if err != nil {
		l.WithError(err).Error("❌ monitor failed")
		fails[mon] = err
		wg.ErrChan() <- err
	} else {
		l.Info("✅")
	}
}

// ApplyReport renders a table report for failed releases.
func (p *Plan) ApplyReport(
	releasesFails map[release.Config]error,
	monitorsFails map[monitor.Config]error,
) error {
	nReleases := len(p.body.Releases)
	kReleases := len(releasesFails)
	nMonitors := len(p.body.Monitors)
	kMonitors := len(monitorsFails)

	log.Infof("Releases Success %d / %d", nReleases-kReleases, nReleases)
	log.Infof("Monitors Success %d / %d", nMonitors-kMonitors, nMonitors)

	if len(releasesFails) > 0 {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{fieldName, fieldNamespace, fieldChart, "version", "error"})
		table.SetAutoFormatHeaders(true)
		table.SetBorder(false)

		for r, err := range releasesFails {
			row := []string{
				r.Name(),
				r.Namespace(),
				r.Chart().Name,
				r.Chart().Version,
				err.Error(),
			}

			table.Rich(row, []tablewriter.Colors{
				{},
				{},
				{},
				{},
				FailStatusColor,
			})
		}

		table.Render()

		return ErrDeploy
	}

	if len(monitorsFails) > 0 {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"name", "error"})
		table.SetAutoFormatHeaders(true)
		table.SetBorder(false)

		for r, err := range monitorsFails {
			row := []string{
				r.Name(),
				err.Error(),
			}

			table.Rich(row, []tablewriter.Colors{
				{},
				FailStatusColor,
			})
		}

		table.Render()

		return ErrDeploy
	}

	return nil
}

func (p *Plan) syncReleasesTracked(ctx context.Context, cfg *tracker.Config) error {
	ctxCancel, cancel := context.WithCancel(ctx)
	defer cancel() // Don't forget!

	ids, kubecontext, err := p.trackerSyncObjects(cfg)
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
	err = p.syncReleases(ctx)
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
