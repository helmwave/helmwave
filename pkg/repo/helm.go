package repo

import (
	"context"
	"fmt"
	"github.com/gofrs/flock"
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func (rep *Config) Sync(settings *helm.EnvSettings) {
	Write(settings.RepositoryConfig, &rep.Entry, settings)
}

// TODO it better later
func Write(repofile string, o *repo.Entry, helm *helm.EnvSettings) {
	//Ensure the file directory exists as it is required for file locking
	err := os.MkdirAll(filepath.Dir(repofile), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}

	// Acquire a file lock for process synchronization
	fileLock := flock.New(strings.Replace(repofile, filepath.Ext(repofile), ".lock", 1))
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer fileLock.Unlock()
	}

	f, err := repo.LoadFile(repofile)
	if err != nil {
		panic(err)
	}

	if f.Has(o.Name) {
		fmt.Printf("❎ %q already exists with the same configuration, skipping\n", o.Name)
	} else {
		chartRepo, err := repo.NewChartRepository(o, getter.All(helm))
		if err != nil {
			panic(err)
		}

		_, err = chartRepo.DownloadIndexFile()
		if err != nil {
			fmt.Printf("⚠️ looks like %v is not a valid chart repository or cannot be reached", o.URL)
		}

		f.Update(o)

		if err := f.WriteFile(repofile, 0644); err != nil {
			panic(err)
		}

		fmt.Printf("✅ %q has been added to your repositories\n", o.Name)
	}
}
