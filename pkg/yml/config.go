package yml

import (
	"github.com/zhilyaev/helmwave/pkg/release"
	"github.com/zhilyaev/helmwave/pkg/repo"
)

type Config struct {
	File string
	Body Body
}

type Body struct {
	Project      string
	Version      string
	Repositories []repo.Config
	Releases     []release.Config
}
