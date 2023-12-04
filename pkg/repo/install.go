package repo

import (
	"context"
	"fmt"

	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
)

func (c *config) Install(_ context.Context, settings *helm.EnvSettings, f *repo.File) error {
	if !c.Force && f.Has(c.Name()) {
		existing := f.Get(c.Name())
		if c.Entry != *existing {
			// The input coming in for the name is different from what is already
			// configured. Return an error.
			return NewDuplicateError(c.Name())
		}

		// The add is idempotent so do nothing
		c.Logger().Info("❎ repository already exists with the same configuration, skipping")

		return nil
	}

	chartRepo, err := repo.NewChartRepository(&c.Entry, getter.All(settings))
	if err != nil {
		return fmt.Errorf("failed to install repository %q: %w", c.Name(), err)
	}

	chartRepo.CachePath = settings.RepositoryCache

	// Hang tight while we grab the latest from your chart repositories...
	c.Logger().Debugf("Download IndexFile for %q", chartRepo.Config.Name)
	_, err = chartRepo.DownloadIndexFile()
	if err != nil {
		c.Logger().WithError(err).Warnf("⚠️ looks like %q is not a valid chart repository or can't be reached", c.URL())
	}

	f.Update(&c.Entry)

	return nil
}
