package repo

import (
	"context"
	"github.com/gofrs/flock"
	log "github.com/sirupsen/logrus"
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (rep *Config) Install(settings *helm.EnvSettings) error {
	return Write(settings.RepositoryConfig, &rep.Entry, settings)
}

// TODO it better later
func Write(repofile string, o *repo.Entry, helm *helm.EnvSettings) error {
	//Ensure the file directory exists as it is required for file locking
	err := os.MkdirAll(filepath.Dir(repofile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}

	// Acquire a file lock for process synchronization
	fileLock := flock.New(strings.Replace(repofile, filepath.Ext(repofile), ".lock", 1))
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer func(fileLock *flock.Flock) {
			err := fileLock.Unlock()
			if err != nil {
				log.Errorf("Failed to release flock: %v", err.Error())
			}
		}(fileLock)
	}

	f, err := repo.LoadFile(repofile)
	if err != nil {
		return err
	}

	if f.Has(o.Name) {
		log.Infof("❎ %q already exists with the same configuration, skipping\n", o.Name)
	} else {
		chartRepo, err := repo.NewChartRepository(o, getter.All(helm))
		if err != nil {
			return err
		}

		_, err = chartRepo.DownloadIndexFile()
		if err != nil {
			log.Warnf("⚠️ looks like %v is not a valid chart repository or cannot be reached", o.URL)
		}

		f.Update(o)

		if err := f.WriteFile(repofile, 0644); err != nil {
			return err
		}

		log.Infof("✅ %q has been added to your repositories\n", o.Name)
	}

	return nil
}
