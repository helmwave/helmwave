package plan

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/gofrs/flock"
	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/kubedog"
	"github.com/helmwave/helmwave/pkg/parallel"
	regi "github.com/helmwave/helmwave/pkg/registry"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/dependency"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/werf/kubedog/pkg/kube"
	"github.com/werf/kubedog/pkg/tracker"
	"github.com/werf/kubedog/pkg/trackers/rollout/multitrack"
	helmRepo "helm.sh/helm/v3/pkg/repo"
)

// Up syncs repositories and releases.
func (p *Plan) Up(ctx context.Context, dog *kubedog.Config) error {
	// Run hooks
	err := p.body.Lifecycle.RunPreUp(ctx)
	if err != nil {
		return err
	}

	defer func() {
		err := p.body.Lifecycle.RunPostUp(ctx)
		if err != nil {
			log.Errorf("got an error from postup hooks: %v", err)
		}
	}()

	log.Info("🗄 sync repositories...")
	err = SyncRepositories(ctx, p.body.Repositories)
	if err != nil {
		return err
	}

	log.Info("🗄 sync registries...")
	err = p.syncRegistries(ctx)
	if err != nil {
		return err
	}

	if len(p.body.Releases) == 0 {
		return nil
	}

	log.Info("🛥 sync releases...")

	if dog.Enabled {
		log.Warn("🐶 kubedog is enabled")
		kubedog.FixLog(dog.LogWidth)
		err = p.syncReleasesKubedog(ctx, dog)
	} else {
		err = p.syncReleases(ctx)
	}

	if err != nil {
		return err
	}

	return nil
}

func (p *Plan) syncRegistries(ctx context.Context) (err error) {
	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Registries))

	for i := range p.body.Registries {
		go func(wg *parallel.WaitGroup, reg regi.Config) {
			defer wg.Done()
			err := reg.Install()
			if err != nil {
				wg.ErrChan() <- err
			}
		}(wg, p.body.Registries[i])
	}

	if err := wg.WaitWithContext(ctx); err != nil {
		return err
	}

	return err
}

// SyncRepositories initializes helm repository.yaml file with flock and installs provided repositories.
func SyncRepositories(ctx context.Context, repositories repo.Configs) error {
	log.Trace("🗄 helm repository.yaml: ", helper.Helm.RepositoryConfig)

	// Create if not exists
	if !helper.IsExists(helper.Helm.RepositoryConfig) {
		f, err := helper.CreateFile(helper.Helm.RepositoryConfig)
		if err != nil {
			return err
		}
		if err := f.Close(); err != nil {
			return fmt.Errorf("failed to close fresh helm repository.yaml: %w", err)
		}
	}

	// we need to get a flock first
	lockPath := helper.Helm.RepositoryConfig + ".lock"
	fileLock := flock.New(lockPath)
	lockCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	// We need to unlock in deferred mode in case of any other errors returned
	defer func(fileLock *flock.Flock) {
		err := fileLock.Unlock()
		if err != nil {
			log.Errorf("failed to release flock %s: %v", fileLock.Path(), err)
		}
	}(fileLock)

	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err != nil && !locked {
		return fmt.Errorf("failed to get lock %s: %w", fileLock.Path(), err)
	}

	f, err := helmRepo.LoadFile(helper.Helm.RepositoryConfig)
	if err != nil {
		return fmt.Errorf("failed to load helm repositories file: %w", err)
	}

	// We can't parallel repositories installation as helm manages single repositories.yaml.
	// To prevent data race, we need to either make helm use futex or not parallel at all
	for i := range repositories {
		err := repositories[i].Install(ctx, helper.Helm, f)
		if err != nil {
			return fmt.Errorf("failed to install %s repository: %w", repositories[i].Name(), err)
		}
	}

	err = f.WriteFile(helper.Helm.RepositoryConfig, os.FileMode(0o644))
	if err != nil {
		return fmt.Errorf("failed to write repositories file: %w", err)
	}

	// If we haven't met any errors yet unlock the repository file. Deferred unlock will exit quickly after this.
	if err := fileLock.Unlock(); err != nil {
		return fmt.Errorf("failed to unlock %s: %w", fileLock.Path(), err)
	}

	return nil
}

func (p *planBody) generateDependencyGraph() (*dependency.Graph[uniqname.UniqName, release.Config], error) {
	dependenciesGraph := dependency.NewGraph[uniqname.UniqName, release.Config]()

	for i := range p.Releases {
		rel := p.Releases[i]
		err := dependenciesGraph.NewNode(rel.Uniq(), rel)
		if err != nil {
			return nil, err
		}

		for _, dep := range rel.DependsOn() {
			dependenciesGraph.AddDependency(rel.Uniq(), dep.Uniq())
		}
	}

	err := dependenciesGraph.Build()
	if err != nil {
		return nil, err
	}

	return dependenciesGraph, nil
}

func getParallelLimit(ctx context.Context, releases release.Configs) int {
	parallelLimit, ok := ctx.Value("parallel-limit").(int)
	if !ok {
		parallelLimit = 0
	}
	if parallelLimit == 0 {
		parallelLimit = len(releases)
	}

	return parallelLimit
}

func (p *Plan) syncReleases(ctx context.Context) (err error) {
	dependenciesGraph, err := p.body.generateDependencyGraph()
	if err != nil {
		return err
	}

	parallelLimit := getParallelLimit(ctx, p.body.Releases)

	const msg = "Deploying releases with limited parallelization"
	if parallelLimit == len(p.body.Releases) {
		log.WithField("limit", parallelLimit).Debug(msg)
	} else {
		log.WithField("limit", parallelLimit).Info(msg)
	}

	nodesChan := dependenciesGraph.Run()

	wg := parallel.NewWaitGroup()
	wg.Add(parallelLimit)

	fails := make(map[release.Config]error)

	mu := &sync.Mutex{}

	for i := 0; i < parallelLimit; i++ {
		go p.syncReleasesWorker(ctx, wg, nodesChan, mu, fails)
	}

	if err := wg.WaitWithContext(ctx); err != nil {
		return err
	}

	return p.ApplyReport(fails)
}

func (p *Plan) syncReleasesWorker(
	ctx context.Context,
	wg *parallel.WaitGroup,
	nodesChan <-chan *dependency.Node[release.Config],
	mu *sync.Mutex,
	fails map[release.Config]error,
) {
	for n := range nodesChan {
		p.syncRelease(ctx, wg, n, mu, fails)
	}
	wg.Done()
}

func (p *Plan) syncRelease(
	ctx context.Context,
	wg *parallel.WaitGroup,
	node *dependency.Node[release.Config],
	mu *sync.Mutex,
	fails map[release.Config]error,
) {
	rel := node.Data

	l := rel.Logger()

	l.Info("🛥 deploying... ")

	if _, err := rel.Sync(ctx); err != nil {
		l.WithError(err).Error("❌ failed to deploy")

		if rel.AllowFailure() {
			l.Errorf("release is allowed to fail, markind as succeeded to dependencies")
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
	}
}

// ApplyReport renders a table report for failed releases.
func (p *Plan) ApplyReport(fails map[release.Config]error) error {
	n := len(p.body.Releases)
	k := len(fails)

	log.Infof("Success %d / %d", n-k, n)

	if len(fails) > 0 {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"name", "namespace", "chart", "version", "err"})
		table.SetAutoFormatHeaders(true)
		table.SetBorder(false)

		for r, err := range fails {
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

	return nil
}

func (p *Plan) syncReleasesKubedog(ctx context.Context, kubedogConfig *kubedog.Config) error {
	ctxCancel, cancel := context.WithCancel(ctx)
	defer cancel() // Don't forget!

	specs, kubecontext, err := p.kubedogSyncSpecs(kubedogConfig)
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
	err = p.syncReleases(ctx)
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

func (p *Plan) kubedogSyncSpecs(kubedogConfig *kubedog.Config) (multitrack.MultitrackSpecs, string, error) {
	return p.kubedogSpecs(kubedogConfig, p.kubedogSyncManifest)
}

func (p *Plan) kubedogSyncManifest(rel release.Config) (string, error) {
	return p.manifests[rel.Uniq()], nil
}
