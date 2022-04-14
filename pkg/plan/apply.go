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
	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
	"github.com/werf/kubedog/pkg/kube"
	"github.com/werf/kubedog/pkg/tracker"
	"github.com/werf/kubedog/pkg/trackers/rollout/multitrack"
	helmRepo "helm.sh/helm/v3/pkg/repo"
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

	log.Info("🗄 Sync registries...")
	err = p.syncRegistries()
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

	err = p.syncRegistries()
	if err != nil {
		return err
	}

	if len(p.body.Releases) == 0 {
		return nil
	}

	log.Info("🛥 Sync releases...")

	return p.syncReleasesKubedog(kubedogConfig)
}

func (p *Plan) syncRegistries() (err error) {
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

	if err := wg.Wait(); err != nil {
		return err
	}

	return err
}

// SyncRepositories initializes helm repository.yaml file with flock and installs provided repositories.
func SyncRepositories(repositories repo.Configs) error {
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

func (p *Plan) syncReleasesKubedog(kubedogConfig *kubedog.Config) (err error) {
	err = helper.KubeInit()
	if err != nil {
		return err
	}
	// kube.Context = helper.Helm.KubeContext
	// kube.DefaultNamespace = helper.Helm.Namespace()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Dont forget!

	opts := multitrack.MultitrackOptions{
		StatusProgressPeriod: kubedogConfig.StatusInterval,
		Options: tracker.Options{
			ParentContext: ctx,
			Timeout:       kubedogConfig.Timeout,
			LogsFromTime:  time.Now(),
		},
	}

	specs := p.kubedogSpecs()
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
	err = p.syncReleases()
	if err != nil {
		return err
	}

	// ? Kubedog soft exit
	time.Sleep(kubedogConfig.StatusInterval)
	err = dogroup.Wait()
	if err != nil {
		// Ignore kubedog error
		log.WithError(err).Warn("kubedog has error while watching resources.")
	}

	return nil
}

func (p *Plan) kubedogSpecs() (s multitrack.MultitrackSpecs) {
	for _, rel := range p.body.Releases {
		manifest := kubedog.Parse([]byte(p.manifests[rel.Uniq()]))
		spec, err := kubedog.MakeSpecs(manifest, rel.Namespace())
		if err != nil {
			log.WithError(err).Fatal("kubedog can't parse resources")
		}

		log.WithFields(log.Fields{
			"Deployments":  len(spec.Deployments),
			"Jobs":         len(spec.Jobs),
			"DaemonSets":   len(spec.DaemonSets),
			"StatefulSets": len(spec.StatefulSets),
			"Canaries":     len(spec.Canaries),
			"release":      rel.Uniq(),
		}).Trace("kubedog track resources")

		s.Jobs = append(s.Jobs, spec.Jobs...)
		s.Deployments = append(s.Deployments, spec.Deployments...)
		s.DaemonSets = append(s.DaemonSets, spec.DaemonSets...)
		s.StatefulSets = append(s.StatefulSets, spec.StatefulSets...)
		s.Canaries = append(s.Canaries, spec.Canaries...)
	}

	return s
}
