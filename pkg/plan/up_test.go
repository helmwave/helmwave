package plan_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/kubedog"
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
	helmRelease "helm.sh/helm/v3/pkg/release"
)

type ApplyTestSuite struct {
	suite.Suite

	ctx context.Context
}

//nolint:paralleltest // can't parallel because of flock timeout
func TestApplyTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(ApplyTestSuite))
}

func (ts *ApplyTestSuite) SetupTest() {
	ts.ctx = tests.GetContext(ts.T())
}

func (ts *ApplyTestSuite) TestApplyBadRepoInstallation() {
	tmpDir := ts.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))

	repoName := "blablanami"

	mockedRepo := plan.NewMockRepositoryConfig(ts.T())
	mockedRepo.On("Name").Return(repoName)
	e := errors.New(ts.T().Name())
	mockedRepo.On("Install").Return(e)

	p.SetRepositories(mockedRepo)

	err := p.Up(ts.ctx, &kubedog.Config{})
	ts.Require().ErrorIs(err, e)

	mockedRepo.AssertExpectations(ts.T())
}

func (ts *ApplyTestSuite) TestApplyNoReleases() {
	tmpDir := ts.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))

	mockedRepo := plan.NewMockRepositoryConfig(ts.T())
	mockedRepo.On("Install").Return(nil)

	p.SetRepositories(mockedRepo)
	dog := &kubedog.Config{}

	err := p.Up(ts.ctx, dog)
	ts.Require().NoError(err)

	mockedRepo.AssertExpectations(ts.T())
}

func (ts *ApplyTestSuite) TestApplyFailedRelease() {
	tmpDir := ts.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))

	mockedRelease := plan.NewMockReleaseConfig(ts.T())
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("Uniq").Return()
	mockedRelease.On("DependsOn").Return([]*release.DependsOnReference{})
	e := errors.New(ts.T().Name())
	mockedRelease.On("Sync").Return(&helmRelease.Release{}, e)
	mockedRelease.On("AllowFailure").Return(false)
	mockedRelease.On("Monitors").Return([]release.MonitorReference{})
	mockedRelease.On("KubeContext").Return("")

	p.SetReleases(mockedRelease)

	err := p.Up(ts.ctx, &kubedog.Config{})
	ts.Require().ErrorIs(err, e)

	mockedRelease.AssertExpectations(ts.T())
}

func (ts *ApplyTestSuite) TestApply() {
	tmpDir := ts.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))

	mockedRelease := plan.NewMockReleaseConfig(ts.T())
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("Uniq").Return()
	mockedRelease.On("Sync").Return(&helmRelease.Release{}, nil)
	mockedRelease.On("DependsOn").Return([]*release.DependsOnReference{})
	mockedRelease.On("Monitors").Return([]release.MonitorReference{})
	mockedRelease.On("KubeContext").Return("")

	mockedRepo := plan.NewMockRepositoryConfig(ts.T())
	mockedRepo.On("Install").Return(nil)

	p.SetRepositories(mockedRepo)
	p.SetReleases(mockedRelease)

	err := p.Up(ts.ctx, &kubedog.Config{})
	ts.Require().NoError(err)

	mockedRepo.AssertExpectations(ts.T())
	mockedRelease.AssertExpectations(ts.T())
}
