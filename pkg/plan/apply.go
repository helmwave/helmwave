package plan

import (
	"context"
	"os"
	"time"

	"github.com/gofrs/flock"
	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/parallel"
	"github.com/helmwave/helmwave/pkg/release"
	log "github.com/sirupsen/logrus"
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/repo"
)

func (p *Plan) Apply() (err error) {
	if len(p.body.Releases) == 0 {
		return release.ErrEmpty
	}

	log.Info("🗄 Sync repositories...")
	err = p.syncRepositories()
	if err != nil {
		return err
	}

	log.Info("🛥 Sync releases...")
	err = p.syncReleases()
	if err != nil {
		return err
	}

	return nil
}

func (p *Plan) syncRepositories() (err error) {
	settings := helm.New()
	log.Trace("helm repository.yaml: ", settings.RepositoryConfig)

	f := &repo.File{}
	// Create if not exits
	if !helper.IsExists(settings.RepositoryConfig) {
		f = repo.NewFile()
	} else {
		f, err = repo.LoadFile(settings.RepositoryConfig)
		if err != nil {
			return err
		}
	}

	// Flock
	lockPath := settings.RepositoryConfig + ".lock"
	fileLock := flock.New(lockPath)
	lockCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err == nil && locked {
		defer fileLock.Unlock()
	} else if err != nil {
		return err
	}

	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Repositories))

	for i := range p.body.Repositories {
		go func(wg *parallel.WaitGroup, i int) {
			defer wg.Done()
			err := p.body.Repositories[i].Install(settings, f)
			if err != nil {
				log.Fatal(err)
			}
		}(wg, i)
	}

	err = wg.Wait()
	if err != nil {
		return err
	}

	return f.WriteFile(settings.RepositoryConfig, os.FileMode(0o644))
}

func (p *Plan) syncReleases() (err error) {
	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Releases))

	for i := range p.body.Releases {
		go func(wg *parallel.WaitGroup, rel *release.Config) {
			defer wg.Done()
			log.Info(rel.Uniq(), " deploying...")
			_, err = rel.Sync()
			if err != nil {
				log.Fatal(err)
			}
		}(wg, p.body.Releases[i])
	}

	return wg.Wait()
}
