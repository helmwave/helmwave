package plan

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
	helmRelease "helm.sh/helm/v3/pkg/release"
)

type BuildManifestsTestSuite struct {
	suite.Suite

	ctx context.Context
}

func TestBuildManifestsTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(BuildManifestsTestSuite))
}

func (ts *BuildManifestsTestSuite) SetupTest() {
	ts.ctx = tests.GetContext(ts.T())
}

func (ts *BuildManifestsTestSuite) TestEmptyReleases() {
	p := New(".")
	p.NewBody()

	err := p.buildManifest(ts.ctx)

	ts.Require().NoError(err)
}

func (ts *BuildManifestsTestSuite) TestMultipleReleases() {
	p := New(".")

	rel1 := NewMockReleaseConfig(ts.T())
	u1, _ := uniqname.NewFromString("redis1@defaultblabla")
	rel1.On("ChartDepsUpd").Return(nil)
	rel1.On("DryRun").Return()
	rel1.On("Sync").Return(&helmRelease.Release{}, nil)
	rel1.On("HooksDisabled").Return(false)
	rel1.On("Uniq").Return(u1)

	rel2 := NewMockReleaseConfig(ts.T())
	u2, _ := uniqname.NewFromString("redis2@defaultblabla")
	rel2.On("ChartDepsUpd").Return(nil)
	rel2.On("DryRun").Return()
	rel2.On("Sync").Return(&helmRelease.Release{}, nil)
	rel2.On("HooksDisabled").Return(false)
	rel2.On("Uniq").Return(u2)

	p.SetReleases(rel1, rel2)

	err := p.buildManifest(ts.ctx)

	ts.Require().NoError(err)
	ts.Require().Len(p.manifests, 2)
	ts.Require().Contains(p.manifests, u1)
	ts.Require().Contains(p.manifests, u2)
	ts.Require().Equal("", p.manifests[u1])
	ts.Require().Equal("", p.manifests[u2])

	rel1.AssertExpectations(ts.T())
	rel2.AssertExpectations(ts.T())
}

func (ts *BuildManifestsTestSuite) TestChartDepsUpdError() {
	p := New(".")

	rel := NewMockReleaseConfig(ts.T())
	uniq, _ := uniqname.NewFromString("redis1@defaultblabla")
	errExpected := errors.New(ts.T().Name())
	rel.On("ChartDepsUpd").Return(errExpected)
	rel.On("DryRun").Return()
	rel.On("Sync").Return(&helmRelease.Release{}, nil)
	rel.On("HooksDisabled").Return(false)
	rel.On("Uniq").Return(uniq)

	p.SetReleases(rel)

	err := p.buildManifest(ts.ctx)

	ts.Require().NoError(err)
	ts.Require().Len(p.manifests, 1)
	ts.Require().Contains(p.manifests, uniq)
	ts.Require().Equal("", p.manifests[uniq])

	rel.AssertExpectations(ts.T())
}

func (ts *BuildManifestsTestSuite) TestSyncError() {
	p := New(".")

	rel := NewMockReleaseConfig(ts.T())
	errExpected := errors.New(ts.T().Name())
	rel.On("ChartDepsUpd").Return(nil)
	rel.On("DryRun").Return()
	rel.On("Sync").Return(&helmRelease.Release{}, errExpected)

	p.SetReleases(rel)

	err := p.buildManifest(ts.ctx)

	ts.Require().ErrorIs(err, errExpected)

	rel.AssertExpectations(ts.T())
}

func (ts *BuildManifestsTestSuite) TestDisabledHooks() {
	p := New(".")

	rel := NewMockReleaseConfig(ts.T())
	uniq, _ := uniqname.NewFromString("redis1@defaultblabla")
	rel.On("ChartDepsUpd").Return(nil)
	rel.On("DryRun").Return()
	rel.On("Sync").Return(&helmRelease.Release{
		Manifest: ts.T().Name(),
		Hooks:    []*helmRelease.Hook{{}},
	}, nil)
	rel.On("HooksDisabled").Return(true)
	rel.On("Uniq").Return(uniq)

	p.SetReleases(rel)

	err := p.buildManifest(ts.ctx)

	ts.Require().NoError(err)
	ts.Require().Len(p.manifests, 1)
	ts.Require().Contains(p.manifests, uniq)
	ts.Require().Equal(p.manifests[uniq], ts.T().Name())

	rel.AssertExpectations(ts.T())
}

func (ts *BuildManifestsTestSuite) TestEnabledHooks() {
	p := New(".")

	rel := NewMockReleaseConfig(ts.T())
	uniq, _ := uniqname.NewFromString("redis1@defaultblabla")
	rel.On("ChartDepsUpd").Return(nil)
	rel.On("DryRun").Return()
	rel.On("Sync").Return(&helmRelease.Release{
		Manifest: ts.T().Name(),
		Hooks: []*helmRelease.Hook{
			{
				Path:     ts.T().Name(),
				Manifest: ts.T().Name(),
			},
		},
	}, nil)
	rel.On("HooksDisabled").Return(false)
	rel.On("Uniq").Return(uniq)

	p.SetReleases(rel)

	err := p.buildManifest(ts.ctx)

	ts.Require().NoError(err)
	ts.Require().Len(p.manifests, 1)
	ts.Require().Contains(p.manifests, uniq)
	ts.Require().Equal(p.manifests[uniq], fmt.Sprintf("%[1]s---\n# Source: %[1]s\n%[1]s\n", ts.T().Name()))

	rel.AssertExpectations(ts.T())
}
