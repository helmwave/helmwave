package plan

import (
	"context"
	"errors"
	"github.com/olekukonko/tablewriter"
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

var ErrDeploy = errors.New("deploy failed")

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
	if err != nil && !locked {
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

	err = f.WriteFile(settings.RepositoryConfig, os.FileMode(0o644))
	if err != nil {
		return err
	}

	// Unlock
	return fileLock.Unlock()
}

func (p *Plan) syncReleases() (err error) {
	wg := parallel.NewWaitGroup()
	wg.Add(len(p.body.Releases))

	fails := make(map[*release.Config]error)

	for i := range p.body.Releases {
		go func(wg *parallel.WaitGroup, rel *release.Config) {
			defer wg.Done()
			log.Info(rel.Uniq(), " deploying...")
			_, err = rel.Sync()
			if err != nil {
				log.Errorf("❌ %s: %v", rel.Uniq(), err)

				rel.NotifyFailed()
				fails[rel] = err
			} else {
				rel.NotifySuccess()
				log.Infof("✅ %s", rel.Uniq())
			}
		}(wg, p.body.Releases[i])
	}

	if err = wg.Wait(); err != nil {
		return err
	}

	return p.ApplyReport(fails)
}

func (p *Plan) ApplyReport(fails map[*release.Config]error) error {
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
				r.Name,
				r.Namespace,
				r.Chart.Name,
				r.Chart.Version,
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

		return ErrDeploy
	}

	return nil
}
