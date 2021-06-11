package plan

import (
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/repo"
)

type Plan struct {
	body     *planBody
	dir      string
	fullPath string
}

const planfile = "planfile"

type planBody struct {
	Project      string
	Version      string
	Repositories []*repo.Config
	Releases     []*release.Config
}

func New(dir string) *Plan {
	if dir[len(dir)-1:] != "/" {
		dir += "/"
	}

	plan := &Plan{
		dir:      dir,
		fullPath: dir + planfile,
	}

	return plan
}
