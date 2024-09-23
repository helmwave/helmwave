package release

import (
	"context"
	"errors"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"github.com/helmwave/helmwave/pkg/fileref"
	"github.com/helmwave/helmwave/pkg/parallel"
)

func (rel *config) BuildValues(ctx context.Context, dir, templater string) error {
	vals := rel.Values()

	wg := parallel.NewWaitGroup()
	wg.Add(len(vals))

	// we need to keep rendered values in memory to use them in other values
	renderedValuesMap := fileref.NewRenderedFiles()
	// we need to keep track of values that we need to delete (e.g. non-existent files)
	toDeleteMap := make(map[*fileref.Config]bool)

	// just in case of dependency cycle or long http requests
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	l := rel.Logger()

	for i := range vals {
		go func(v *fileref.Config) {
			defer wg.Done()

			l := l.WithField("values", v)

			filename := filepath.Join(dir, "values", rel.Uniq().String(), strconv.Itoa(i)+".yml")
			data := struct {
				Release Config
			}{
				Release: rel,
			}

			err := v.Set(ctx, filename, templater, data, renderedValuesMap)
			switch {
			case !v.Strict && errors.Is(err, fileref.ErrValuesNotExist):
				l.WithError(err).Warn("skipping values...")
				toDeleteMap[v] = true
			case err != nil:
				l.WithError(err).Error("failed to build values")

				wg.ErrChan() <- err
			}
		}(&vals[i])
	}

	err := wg.WaitWithContext(ctx)
	if err != nil {
		return err
	}

	for i := len(vals) - 1; i >= 0; i-- {
		if toDeleteMap[&vals[i]] {
			vals = slices.Delete(vals, i, i+1)
		}
	}
	rel.ValuesF = vals

	return nil
}
