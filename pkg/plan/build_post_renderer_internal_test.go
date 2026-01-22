package plan

import (
	"context"
	"errors"
	"testing"

	"github.com/helmwave/helmwave/pkg/hooks"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
	helmRelease "helm.sh/helm/v3/pkg/release"
)

type BuildPostRendererTestSuite struct {
	suite.Suite

	ctx context.Context
}

func TestBuildPostRendererTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(BuildPostRendererTestSuite))
}

func (ts *BuildPostRendererTestSuite) SetupTest() {
	ts.ctx = tests.GetContext(ts.T())
}

func (ts *BuildPostRendererTestSuite) TestSuccess() {
	p := New(".")

	rel := NewMockReleaseConfig(ts.T())
	uniq, _ := uniqname.NewFromString("redis@default")
	rel.On("ChartDepsUpd").Return(nil)
	rel.On("DryRun").Return()
	rel.On("Sync").Return(&helmRelease.Release{}, nil)
	rel.On("HooksDisabled").Return(false)
	rel.On("Uniq").Return(uniq)
	rel.On("DependsOn").Return([]*release.DependsOnReference{})
	rel.On("Lifecycle").Return(hooks.Lifecycle{})
	rel.On("BuildValues").Return(map[string]string{}, nil)
	rel.On("Values").Return([]release.ValuesReference{})
	rel.On("BuildPostRenderer").Return(nil)

	p.SetReleases(rel)

	err := p.buildManifest(ts.ctx)

	ts.Require().NoError(err)
	rel.AssertCalled(ts.T(), "BuildPostRenderer")
}

func (ts *BuildPostRendererTestSuite) TestBuildError() {
	p := New(".")

	errBuild := errors.New("post-renderer build error")

	rel := NewMockReleaseConfig(ts.T())
	uniq, _ := uniqname.NewFromString("redis@default")
	rel.On("ChartDepsUpd").Return(nil)
	rel.On("Uniq").Return(uniq)
	rel.On("DependsOn").Return([]*release.DependsOnReference{})
	rel.On("AllowFailure").Return(false)
	rel.On("Values").Return([]release.ValuesReference{})
	rel.On("BuildValues").Return(map[string]string{}, nil)
	rel.On("BuildPostRenderer").Return(errBuild)
	rel.On("Lifecycle").Return(hooks.Lifecycle{})

	p.SetReleases(rel)

	err := p.buildManifest(ts.ctx)

	ts.Require().ErrorIs(err, errBuild)
	rel.AssertCalled(ts.T(), "BuildPostRenderer")
}
