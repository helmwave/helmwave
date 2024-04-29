package plan

import (
	"github.com/helmwave/helmwave/pkg/helper"
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

	p.body.Releases = helper.SlicesMap(r, func(r *MockReleaseConfig) release.Config {
		return r
	})
}

func (p *Plan) SetRepositories(r ...*MockRepositoryConfig) {
	if p.body == nil {
		p.NewBody()
	}

	p.body.Repositories = helper.SlicesMap(r, func(r *MockRepositoryConfig) repo.Config {
		return r
	})
}

func (p *Plan) SetRegistries(r ...*MockRegistryConfig) {
	if p.body == nil {
		p.NewBody()
	}

	p.body.Registries = helper.SlicesMap(r, func(r *MockRegistryConfig) regi.Config {
		return r
	})
}
