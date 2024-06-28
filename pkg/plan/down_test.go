package plan_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
	helmRelease "helm.sh/helm/v3/pkg/release"
)

type DestroyTestSuite struct {
	suite.Suite

	ctx context.Context
}

func TestDestroyTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(DestroyTestSuite))
}

func (ts *DestroyTestSuite) SetupTest() {
	ts.ctx = tests.GetContext(ts.T())
}

func (ts *DestroyTestSuite) TestDestroy() {
	tmpDir := ts.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))

	mockedRelease := plan.NewMockReleaseConfig(ts.T())
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("KubeContext").Return("")
	mockedRelease.On("Uniq").Return()
	mockedRelease.On("Uninstall").Return(&helmRelease.UninstallReleaseResponse{}, nil)
	mockedRelease.On("DependsOn").Return([]*release.DependsOnReference{})

	p.SetReleases(mockedRelease)

	err := p.Down(ts.ctx)
	ts.Require().NoError(err)

	mockedRelease.AssertExpectations(ts.T())
}

func (ts *DestroyTestSuite) TestDestroyFailedRelease() {
	tmpDir := ts.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))

	mockedRelease := plan.NewMockReleaseConfig(ts.T())
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("KubeContext").Return("")
	mockedRelease.On("Uniq").Return()
	e := errors.New(ts.T().Name())
	mockedRelease.On("Uninstall").Return(&helmRelease.UninstallReleaseResponse{}, e)
	mockedRelease.On("DependsOn").Return([]*release.DependsOnReference{})

	p.SetReleases(mockedRelease)

	err := p.Down(ts.ctx)
	ts.Require().ErrorIs(err, e)

	mockedRelease.AssertExpectations(ts.T())
}

func (ts *DestroyTestSuite) TestDestroyNoReleases() {
	tmpDir := ts.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))
	p.NewBody()

	err := p.Down(ts.ctx)
	ts.Require().NoError(err)
}
