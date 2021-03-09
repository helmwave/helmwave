package yml

import (
	"github.com/zhilyaev/helmwave/pkg/release"
	"github.com/zhilyaev/helmwave/pkg/repo"
)

type Config struct {
	Project            string
	Version            string
	EnableDependencies bool
	Repositories       []*repo.Config
	Releases           []*release.Config
}
