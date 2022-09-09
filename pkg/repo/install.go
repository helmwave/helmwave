package repo

import (
	"context"
	"fmt"

	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
)

func (rep *config) rep2Entry() {
	rep.entry.Name = rep.Name()
	rep.entry.URL = rep.URL()
	rep.entry.Username = rep.Username
	rep.entry.Password = rep.Password
	rep.entry.CertFile = rep.CertFile
	rep.entry.KeyFile = rep.KeyFile
	rep.entry.CAFile = rep.CAFile
	rep.entry.InsecureSkipTLSverify = rep.InsecureSkipTLSverify
	rep.entry.PassCredentialsAll = rep.PassCredentialsAll
}

func (rep *config) Install(ctx context.Context, settings *helm.EnvSettings, f *repo.File) error {
	rep.rep2Entry()

	if !rep.Force && f.Has(rep.Name()) {
		existing := f.Get(rep.Name())
		if rep.entry != existing {
			// The input coming in for the name is different from what is already
			// configured. Return an error.
			return fmt.Errorf(
				"❌ repository name (%q) already exists with different Repository, cannot overwrite it without force",
				rep.Name(),
			)
		}

		// The add is idempotent so do nothing
		rep.Logger().Info("❎ repository already exists with the same configuration, skipping")

		return nil
	}

	chartRepo, err := repo.NewChartRepository(rep.entry, getter.All(settings))
	if err != nil {
		return fmt.Errorf("failed to install repository %q: %w", rep.Name(), err)
	}

	chartRepo.CachePath = settings.RepositoryCache

	// Hang tight while we grab the latest from your chart repositories...
	rep.Logger().Debugf("Download IndexFile for %q", chartRepo.Config.Name)
	_, err = chartRepo.DownloadIndexFile()
	if err != nil {
		rep.Logger().WithError(err).Warnf("⚠️ looks like %q is not a valid chart repository or cannot be reached", rep.URL())
	}

	f.Update(rep.entry)

	return nil
}
