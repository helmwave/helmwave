package release

import (
	"context"
	"github.com/helmwave/helmwave/pkg/fileref"
	"path/filepath"
	"slices"
	"strconv"
	"time"

	"github.com/helmwave/helmwave/pkg/parallel"
)

// SetValuesDst sets the path inside plandir for the given value file
func (rel *config) SetValuesDst(dir string, v *fileref.Config) *fileref.Config {
	i := strconv.Itoa(slices.Index(rel.Values(), *v))
	v.Dst = filepath.Join(dir, "values", rel.Uniq().String(), i+".yml")

	return v
}

func (rel *config) BuildValues(ctx context.Context, dir string) error {
	vals := rel.Values()

	wg := parallel.NewWaitGroup()
	wg.Add(len(vals))

	// just in case of dependency cycle or long http requests
	ctx, cancel := context.WithTimeout(ctx, time.Minute)
	defer cancel()

	l := rel.Logger()

	for i := range vals {
		go func(v *fileref.Config) {
			defer wg.Done()

			// Save values inside plandir
			rel.SetValuesDst(dir, v)

			l := l.WithField("values", v)
			err := v.Render(ctx, rel, l)
			if err != nil {
				l.WithError(err).Error("failed to build values")
				wg.ErrChan() <- err
			}

		}(&vals[i])
	}

	err := wg.WaitWithContext(ctx)
	if err != nil {
		return err
	}

	//for i := len(vals) - 1; i >= 0; i-- {
	//	if toDeleteMap[&vals[i]] {
	//		vals = slices.Delete(vals, i, i+1)
	//	}
	//}

	rel.ValuesF = vals

	return nil
}
