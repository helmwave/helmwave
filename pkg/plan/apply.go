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
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/werf/kubedog/pkg/kube"
	"github.com/werf/kubedog/pkg/tracker"
	"github.com/werf/kubedog/pkg/trackers/rollout/multitrack"
	helmRepo "helm.sh/helm/v3/pkg/repo"
	"k8s.io/client-go/kubernetes"
)

// ErrDeploy is returned when deploy is failed for whatever reason.
var ErrDeploy = errors.New("deploy failed")

// Apply syncs repositories and releases.
func (p *Plan) Apply() (err error) {
	log.Info("🗄 Sync repositories...")
	err = SyncRepositories(p.body.Repositories)
	if err != nil {
		return err
	}

	if len(p.body.Releases) == 0 {
		return nil
	}

	log.Info("🛥 Sync releases...")

	return p.syncReleases()
}

// ApplyWithKubedog runs kubedog in goroutine and syncs repositories and releases.
func (p *Plan) ApplyWithKubedog(kubedogConfig *kubedog.Config) (err error) {
	log.Info("🗄 Sync repositories...")
	err = SyncRepositories(p.body.Repositories)
	if err != nil {
		return err
	}

	if len(p.body.Releases) == 0 {
		return nil
	}

	log.Info("🛥 Sync releases...")

	return p.syncReleasesKubedog(kubedogConfig)
}

// SyncRepositories initializes helm repository.yaml file with flock and installs provided repositories.
func SyncRepositories(repositories repoConfigs) error {
	log.Trace("🗄 helm repository.yaml: ", helper.Helm.RepositoryConfig)

	// Create if not exits
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
	lockCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	// We need to unlock in deferred mode in case of any other errors returned
	defer fileLock.Unlock() //nolint:errcheck // TODO: add error checking
	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err != nil && !locked {
		return fmt.Errorf("failed to get lock %s: %w", fileLock.Path(), err)
	}

	f, err := helmRepo.LoadFile(helper.Helm.RepositoryConfig)
	if err != nil {
		return fmt.Errorf("failed to load helm repositories file: %w", err)
	}

	// We cannot parallel repositories installation as helm manages single repositories.yaml.
	// To prevent data race we need either make helm use futex or not parallel at all
	for i := range repositories {
		err := repositories[i].Install(helper.Helm, f)
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

func (p *Plan) syncReleases() (err error) {
	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Releases))

	fails := make(map[release.Config]error)

	mu := &sync.Mutex{}

	for i := range p.body.Releases {
		p.body.Releases[i].HandleDependencies(p.body.Releases)
		go func(wg *parallel.WaitGroup, rel release.Config, mu *sync.Mutex) {
			defer wg.Done()
			l := log.WithField("release", rel.Uniq())
			l.Info("🛥 deploying... ")
			_, err = rel.Sync()
			if err != nil {
				l.WithError(err).Error("❌")

				rel.NotifyFailed()

				mu.Lock()
				fails[rel] = err
				mu.Unlock()

				wg.ErrChan() <- err
			} else {
				rel.NotifySuccess()
				l.Info("✅")
			}
		}(wg, p.body.Releases[i], mu)
	}

	if err := wg.Wait(); err != nil {
		return err
	}

	return p.ApplyReport(fails)
}

// ApplyReport renders table report for failed releases.
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

func (p *Plan) syncReleasesKubedog(kubedogConfig *kubedog.Config) error {
	mapSpecs, err := p.kubedogSpecs()
	if err != nil {
		return err
	}

	wg := parallel.NewWaitGroup()
	// wg.Add(len(p.body.Releases))
	wg.Add(1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// KubeInit
	err = helper.KubeInit()
	if err != nil {
		return err
	}

	err = runMultiracks(ctx, mapSpecs, kubedogConfig, wg)
	if err != nil {
		return err
	}

	go func(wg *parallel.WaitGroup, cancel context.CancelFunc) {
		defer wg.Done()
		defer cancel()
		wg.ErrChan() <- p.syncReleases()
	}(wg, cancel)

	return wg.Wait()
}

func runMultiracks(
	ctx context.Context,
	mapSpecs map[string]*multitrack.MultitrackSpecs,
	kubedogConfig *kubedog.Config,
	wg *parallel.WaitGroup) error {
	opts := multitrack.MultitrackOptions{
		StatusProgressPeriod: kubedogConfig.StatusInterval,
		Options: tracker.Options{
			ParentContext: ctx,
			Timeout:       kubedogConfig.Timeout,
			LogsFromTime:  time.Now(),
		},
	}

	for ns, specs := range mapSpecs {
		// Todo Test it with different namespace
		kube.Context = helper.Helm.KubeContext
		kube.DefaultNamespace = ns

		log.Info("🐶 kubedog for ", ns)

		go func(
			delay time.Duration,
			kubeClient kubernetes.Interface,
			specs multitrack.MultitrackSpecs,
			opts multitrack.MultitrackOptions,
			wg *parallel.WaitGroup,
		) {
			defer wg.Done()
			time.Sleep(delay)
			wg.Add(1)

			wg.ErrChan() <- multitrack.Multitrack(kubeClient, specs, opts)
		}(kubedogConfig.StartDelay, kube.Client, *specs, opts, wg)
	}

	return nil
}

func (p *Plan) kubedogSpecs() (map[string]*multitrack.MultitrackSpecs, error) {
	mapSpecs := make(map[string]*multitrack.MultitrackSpecs)

	for _, rel := range p.body.Releases {
		manifest := kubedog.Parse([]byte(p.manifests[rel.Uniq()]))
		relSpecs, err := kubedog.MakeSpecs(manifest, rel.Namespace())
		if err != nil {
			return nil, err
		}

		log.WithFields(log.Fields{
			"Deployments":  len(relSpecs.Deployments),
			"Jobs":         len(relSpecs.Jobs),
			"DaemonSets":   len(relSpecs.DaemonSets),
			"StatefulSets": len(relSpecs.StatefulSets),
		}).Tracef("%s specs", rel.Uniq())

		nsSpec, found := mapSpecs[rel.Namespace()]
		if found {
			// Merge
			nsSpec.DaemonSets = append(nsSpec.DaemonSets, relSpecs.DaemonSets...)
			nsSpec.Deployments = append(nsSpec.Deployments, relSpecs.Deployments...)
			nsSpec.StatefulSets = append(nsSpec.StatefulSets, relSpecs.StatefulSets...)
			nsSpec.Jobs = append(nsSpec.Jobs, relSpecs.Jobs...)
			mapSpecs[rel.Namespace()] = nsSpec
		} else {
			mapSpecs[rel.Namespace()] = relSpecs
		}
	}

	return mapSpecs, nil
}
