package release

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/helmwave/helmwave/pkg/parallel"
)

func (rel *config) BuildStoreFiles(ctx context.Context, dir, templater string) error {
	files := rel.StoreFiles

	wg := parallel.NewWaitGroup()
	wg.Add(len(files))

	// we need to keep rendered values in memory to use them in other values
	renderedMap := newRenderedValuesFiles()
	// we need to keep track of values that we need to delete (e.g. non-existent files)
	toDeleteMap := make(map[*ValuesReference]bool)

	// just in case of dependency cycle or long http requests
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	l := rel.Logger()

	for i := range files {
		go func(v *ValuesReference) {
			defer wg.Done()

			l := l.WithField("store file", v)

			err := v.SetViaRelease(ctx, rel, dir, templater, renderedMap)
			switch {
			case !v.Strict && errors.Is(ErrValuesNotExist, err):
				l.WithError(err).Warn("skipping store file...")
				toDeleteMap[v] = true
			case err != nil:
				l.WithError(err).Error("failed to build store files")

				wg.ErrChan() <- err
			}
		}(&files[i])
	}

	err := wg.WaitWithContext(ctx)
	if err != nil {
		return err
	}

	for i := len(files) - 1; i >= 0; i-- {
		if toDeleteMap[&files[i]] {
			files = slices.Delete(files, i, i+1)
		}
	}
	rel.StoreFiles = files

	return nil
}
