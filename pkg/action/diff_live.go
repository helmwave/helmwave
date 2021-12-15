package action

import (
	"os"

	"github.com/helmwave/helmwave/pkg/plan"
)

type DiffLive struct {
	diff    *Diff
	plandir string
}

func (d *DiffLive) Run() error {
	p := plan.New(d.plandir)
	if err := p.Import(); err != nil {
		return err
	}
	if ok := p.IsManifestExist(); !ok {
		return os.ErrNotExist
	}

	p.DiffLive(d.diff.ShowSecret, d.diff.Wide)

	return nil
}
