package release

import (
	"context"
	"path/filepath"
	"strconv"

	"github.com/helmwave/helmwave/pkg/templater"
)

func (rel *config) BuildValues(ctx context.Context, dir string, templater templater.Templater) error {
	vals := rel.Values()

	for i, _ := range vals {
		vals[i].Dst = filepath.Join(dir, "values", strconv.Itoa(i))
		vals[i].Run(ctx)
	}

	return nil

}
