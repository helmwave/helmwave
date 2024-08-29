package plan

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gofrs/flock"
	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/repo"
	log "github.com/sirupsen/logrus"
	helmRepo "helm.sh/helm/v3/pkg/repo"
)

// SyncRepositories initializes helm repository.yaml file with flock and installs provided repositories.
func SyncRepositories(ctx context.Context, repositories repo.Configs) error {
	log.Trace("🗄 helm repository.yaml: ", helper.Helm.RepositoryConfig)

	// Create if not exists
	if !helper.IsExists(helper.Helm.RepositoryConfig) {
		f, err := helper.CreateFile(helper.Helm.RepositoryConfig)
		if err != nil {
			return err
		}
		if err := f.Close(); err != nil {
			return fmt.Errorf("failed to close fresh helm repository.yaml: %w", err)
		}
	}

	// we need to get a flock first
	lockPath := helper.Helm.RepositoryConfig + ".lock"
	fileLock := flock.New(lockPath)
	lockCtx, cancel := context.WithTimeout(ctx, 20*time.Second)
	defer cancel()

	// We need to unlock in deferred mode in case of any other errors returned
	defer func(fileLock *flock.Flock) {
		err := fileLock.Unlock()
		if err != nil {
			log.Errorf("failed to release flock %s: %v", fileLock.Path(), err)
		}
	}(fileLock)

	locked, err := fileLock.TryLockContext(lockCtx, time.Second)
	if err != nil && !locked {
		return fmt.Errorf("failed to get lock %s: %w", fileLock.Path(), err)
	}

	f, err := helmRepo.LoadFile(helper.Helm.RepositoryConfig)
	if err != nil {
		return fmt.Errorf("failed to load helm repositories file: %w", err)
	}

	// We can't parallel repositories installation as helm manages single repositories.yaml.
	// To prevent data race, we need to either make helm use futex or not parallel at all
	for i := range repositories {
		err := repositories[i].Install(ctx, helper.Helm, f)
		if err != nil {
			return fmt.Errorf("failed to install %s repository: %w", repositories[i].Name(), err)
		}
	}

	err = f.WriteFile(helper.Helm.RepositoryConfig, os.FileMode(0o644))
	if err != nil {
		return fmt.Errorf("failed to write repositories file: %w", err)
	}

	// If we haven't met any errors yet unlock the repository file. Deferred unlock will exit quickly after this.
	if err := fileLock.Unlock(); err != nil {
		return fmt.Errorf("failed to unlock %s: %w", fileLock.Path(), err)
	}

	return nil
}
