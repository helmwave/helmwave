package plan

import (
	regi "github.com/helmwave/helmwave/pkg/registry"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/repo"
)

func (p *Plan) NewBody() *planBody {
	p.body = &planBody{}

	return p.body
}

func (p *Plan) SetReleases(r ...*MockReleaseConfig) {
	if p.body == nil {
		p.NewBody()
	}
	c := make(release.Configs, len(r))
	for i := range r {
		c[i] = r[i]
	}
	p.body.Releases = c
}

func (p *Plan) SetRepositories(r ...*MockRepositoryConfig) {
	if p.body == nil {
		p.NewBody()
	}
	c := make(repo.Configs, len(r))
	for i := range r {
		c[i] = r[i]
	}
	p.body.Repositories = c
}

func (p *Plan) SetRegistries(r ...*MockRegistryConfig) {
	if p.body == nil {
		p.NewBody()
	}
	c := make(regi.Configs, len(r))
	for i := range r {
		c[i] = r[i]
	}
	p.body.Registries = c
}
