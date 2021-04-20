package yml

import (
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/repo"
)

type Config struct {
	Project      string
	Version      string
	Repositories []*repo.Config
	Releases     []*release.Config
}
