package repo

import (
	"fmt"
	"github.com/helmwave/helmwave/pkg/helper"
	"helm.sh/helm/v3/pkg/registry"

	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
)

func (rep *config) Install(settings *helm.EnvSettings, f *repo.File) error {
	if rep.OCI {
		return rep.installOCI()
	}

	return rep.install(settings, f)
}

func (rep *config) install(settings *helm.EnvSettings, f *repo.File) error {
	if !rep.Force && f.Has(rep.Name()) {
		existing := f.Get(rep.Name())
		if rep.Entry != *existing {
			// The input coming in for the name is different from what is already
			// configured. Return an error.
			return fmt.Errorf("❌ repository name (%q) already exists, please specify a different name", rep.Name())
		}

		// The add is idempotent so do nothing
		rep.Logger().Info("❎ repository already exists with the same configuration, skipping")

		return nil
	}

	chartRepo, err := repo.NewChartRepository(&rep.Entry, getter.All(settings))
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

	f.Update(&rep.Entry)

	return nil
}

func (rep *config) installOCI() error {
	return helper.HelmRegistryClient.Login(
		rep.Name(),
		registry.LoginOptBasicAuth(rep.Username, rep.Password),
		registry.LoginOptInsecure(rep.InsecureSkipTLSverify))
}
