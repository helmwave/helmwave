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
	"io/ioutil"
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
	for _, v := range c.Releases {
		v.Options.DryRun = true
	}
	return c.Sync(manifestPath, async, settings)
}

func (c *Config) SyncWithKubedog(manifestPath string, async bool, settings *helm.EnvSettings) error {
	err := c.SyncFake(manifestPath, async, settings)
	if err != nil {
		return err
	}

	var specs multitrack.MultitrackSpecs

	for _, rel := range c.Releases {
		// Todo mv to "copy"
		rel.Options.DryRun = false
		src, err := ioutil.ReadFile(manifestPath + rel.UniqName() + ".yml")
		if err != nil {
			return err
		}
		manifest := kubedog.MakeManifest(src)
		relSpecs, err := kubedog.MakeSpecs(manifest)
		log.WithFields(log.Fields{
			"Deployments":  relSpecs.Deployments,
			"Jobs":         relSpecs.Jobs,
			"DaemonSets":   relSpecs.DaemonSets,
			"StatefulSets": relSpecs.StatefulSets,
		}).Debug("Kubedog of ", rel.UniqName())
		if err != nil {
			return err
		}

		// Merge specs
		specs.DaemonSets = append(specs.DaemonSets, relSpecs.DaemonSets...)
		specs.Deployments = append(specs.Deployments, relSpecs.Deployments...)
		specs.StatefulSets = append(specs.StatefulSets, relSpecs.StatefulSets...)
		specs.Jobs = append(specs.Jobs, relSpecs.Jobs...)
	}

	progress, _ := time.ParseDuration("5s")
	timeout, _ := time.ParseDuration("5m")

	err = kube.Init(kube.InitOptions{})
	if err != nil {
		return err
	}

	err = multitrack.Multitrack(kube.Kubernetes,
		specs,
		multitrack.MultitrackOptions{
			StatusProgressPeriod: progress,
			Options: tracker.Options{
				Timeout:      timeout,
				LogsFromTime: time.Now(),
			},
		})
	if err != nil {
		return err
	}

	g := &parallel.Group{}
	log.Debug("🐞 Run sync with kubedog")
	g.Go(c.SyncReleases, manifestPath, async)
	return g.Wait()
}
