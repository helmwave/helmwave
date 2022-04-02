package plan_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/stretchr/testify/suite"
	helmRelease "helm.sh/helm/v3/pkg/release"
)

type DestroyTestSuite struct {
	suite.Suite
}

func (s *DestroyTestSuite) TestDestroy() {
	tmpDir := s.T().TempDir()
	p, err := plan.New(filepath.Join(tmpDir, plan.Dir))
	s.Require().NoError(err)

	mockedRelease := &plan.MockReleaseConfig{}
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("Uniq").Return()
	mockedRelease.On("Uninstall").Return(&helmRelease.UninstallReleaseResponse{}, nil)

	p.SetReleases(mockedRelease)

	err = p.Destroy()
	s.Require().NoError(err)

	mockedRelease.AssertExpectations(s.T())
}

func (s *DestroyTestSuite) TestDestroyFailedRelease() {
	tmpDir := s.T().TempDir()
	p, err := plan.New(filepath.Join(tmpDir, plan.Dir))
	s.Require().NoError(err)

	mockedRelease := &plan.MockReleaseConfig{}
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("Uniq").Return()
	e := errors.New(s.T().Name())
	mockedRelease.On("Uninstall").Return(&helmRelease.UninstallReleaseResponse{}, e)

	p.SetReleases(mockedRelease)

	err = p.Destroy()
	s.Require().ErrorIs(err, e)

	mockedRelease.AssertExpectations(s.T())
}

func (s *DestroyTestSuite) TestDestroyNoReleases() {
	tmpDir := s.T().TempDir()
	p, err := plan.New(filepath.Join(tmpDir, plan.Dir))
	s.Require().NoError(err)
	p.NewBody()

	err = p.Destroy()
	s.Require().NoError(err)
}

func TestDestroyTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(DestroyTestSuite))
}
