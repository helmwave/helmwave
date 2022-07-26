package plan_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/plan"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	helmRelease "helm.sh/helm/v3/pkg/release"
)

type ApplyTestSuite struct {
	suite.Suite
}

func (s *ApplyTestSuite) TestApplyBadRepoInstallation() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))

	repoName := "blablanami"

	mockedRepo := &plan.MockRepoConfig{}
	mockedRepo.On("Name").Return(repoName)
	e := errors.New(s.T().Name())
	mockedRepo.On("Install").Return(e)

	p.SetRepositories(mockedRepo)

	err := p.Apply()
	s.Require().ErrorIs(err, e)

	mockedRepo.AssertExpectations(s.T())
}

func (s *ApplyTestSuite) TestApplyNoReleases() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))

	mockedRepo := &plan.MockRepoConfig{}
	mockedRepo.On("Install").Return(nil)

	p.SetRepositories(mockedRepo)

	err := p.Apply()
	s.Require().NoError(err)

	mockedRepo.AssertExpectations(s.T())
}

func (s *ApplyTestSuite) TestApplyFailedRelease() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))

	mockedRelease := &plan.MockReleaseConfig{}
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("Uniq").Return()
	mockedRelease.On("DependsOn").Return([]string{})
	mockedRelease.On("Logger").Return(log.WithField("test", s.T().Name()))
	e := errors.New(s.T().Name())
	mockedRelease.On("Sync").Return(&helmRelease.Release{}, e)
	mockedRelease.On("AllowFailure").Return(false)

	p.SetReleases(mockedRelease)

	err := p.Apply()
	s.Require().ErrorIs(err, e)

	mockedRelease.AssertExpectations(s.T())
}

func (s *ApplyTestSuite) TestApply() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))

	mockedRelease := &plan.MockReleaseConfig{}
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("Uniq").Return()
	mockedRelease.On("Sync").Return(&helmRelease.Release{}, nil)
	mockedRelease.On("DependsOn").Return([]string{})
	mockedRelease.On("Logger").Return(log.WithField("test", s.T().Name()))

	mockedRepo := &plan.MockRepoConfig{}
	mockedRepo.On("Install").Return(nil)

	p.SetRepositories(mockedRepo)
	p.SetReleases(mockedRelease)

	err := p.Apply()
	s.Require().NoError(err)

	mockedRepo.AssertExpectations(s.T())
	mockedRelease.AssertExpectations(s.T())
}

//nolint:paralleltest // cannot parallel because of flock timeout
func TestApplyTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(ApplyTestSuite))
}
