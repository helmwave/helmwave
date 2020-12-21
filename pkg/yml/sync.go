package yml

import (
	"github.com/zhilyaev/helmwave/pkg/release"
	"github.com/zhilyaev/helmwave/pkg/repo"
	helm "helm.sh/helm/v3/pkg/cli"
)

func (c *Body) SyncRepos(settings *helm.EnvSettings) error {
	return repo.Sync(c.Repositories, settings)
}

func (c *Body) SyncReleases(manifestPath string, async bool) error {
	return release.Sync(c.Releases, manifestPath, async)
}

func (c *Body) Sync(manifestPath string, async bool, settings *helm.EnvSettings) (err error) {
	err = c.SyncRepos(settings)
	if err != nil {
		return err
	}

	return c.SyncReleases(manifestPath, async)
}
