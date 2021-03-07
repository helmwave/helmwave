package yml

import (
	log "github.com/sirupsen/logrus"
	"github.com/wayt/parallel"
	"github.com/werf/kubedog/pkg/kube"
	"github.com/werf/kubedog/pkg/tracker"
	"github.com/werf/kubedog/pkg/trackers/rollout/multitrack"
	"github.com/zhilyaev/helmwave/pkg/kubedog"
	"github.com/zhilyaev/helmwave/pkg/release"
	"github.com/zhilyaev/helmwave/pkg/repo"
	helm "helm.sh/helm/v3/pkg/cli"
	"time"
)

func (c *Config) SyncRepos(settings *helm.EnvSettings) error {
	return repo.Sync(c.Repositories, settings)
}

func (c *Config) SyncReleases(manifestPath string, async bool) error {
	return release.Sync(c.Releases, manifestPath, async)
}

func (c *Config) Sync(manifestPath string, async bool, settings *helm.EnvSettings) (err error) {
	err = c.SyncRepos(settings)
	if err != nil {
		return err
	}

	return c.SyncReleases(manifestPath, async)
}

func (c *Config) SyncFake(manifestPath string, async bool, settings *helm.EnvSettings) error {
	log.Info("🛫 Fake deploy")
	for i := range c.Releases {
		c.Releases[i].Options.DryRun = true
	}
	return c.Sync(manifestPath, async, settings)
}

func (c *Config) SyncWithKubedog(manifestPath string, async bool, settings *helm.EnvSettings, kubedogConfig *kubedog.Config) error {
	err := c.SyncFake(manifestPath, async, settings)
	if err != nil {
		return err
	}
	log.Debug("🛫 Fake deploy has been finished")

	mapSpecs, err := release.MakeMapSpecs(c.Releases, manifestPath)
	if err != nil {
		return err
	}

	opts := multitrack.MultitrackOptions{
		StatusProgressPeriod: kubedogConfig.StatusInterval,
		Options: tracker.Options{
			Timeout:      kubedogConfig.Timeout,
			LogsFromTime: time.Now(),
		},
	}

	goSpecs := &parallel.Group{}
	for ns, specs := range mapSpecs {
		log.Info("🐶 kubedog for ", ns)
		// Needs to testing with several  ns
		err := kube.Init(kube.InitOptions{})
		if err != nil {
			return err
		}
		kube.Context = settings.KubeContext
		kube.DefaultNamespace = ns

		//multitrack.Multitrack(client, *specs, opts)
		go func(delay time.Duration, f interface{}, args ...interface{}) {
			time.Sleep(delay)
			goSpecs.Go(f, args...)
		}(kubedogConfig.StartDelay, multitrack.Multitrack, kube.Kubernetes, *specs, opts)
	}

	g := &parallel.Group{}
	g.Go(c.SyncReleases, manifestPath, async)

	err = g.Wait()
	if err != nil {
		return err
	}

	return goSpecs.Wait()
}
