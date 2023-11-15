package plan_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/stretchr/testify/suite"
	helmRelease "helm.sh/helm/v3/pkg/release"
)

type DestroyTestSuite struct {
	suite.Suite
}

func (s *DestroyTestSuite) TestDestroy() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))

	mockedRelease := &plan.MockReleaseConfig{}
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("Uniq").Return()
	mockedRelease.On("Uninstall").Return(&helmRelease.UninstallReleaseResponse{}, nil)
	mockedRelease.On("DependsOn").Return([]*release.DependsOnReference{})

	p.SetReleases(mockedRelease)

	err := p.Down(context.Background())
	s.Require().NoError(err)

	mockedRelease.AssertExpectations(s.T())
}

func (s *DestroyTestSuite) TestDestroyFailedRelease() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))

	mockedRelease := &plan.MockReleaseConfig{}
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("Uniq").Return()
	e := errors.New(s.T().Name())
	mockedRelease.On("Uninstall").Return(&helmRelease.UninstallReleaseResponse{}, e)
	mockedRelease.On("DependsOn").Return([]*release.DependsOnReference{})

	p.SetReleases(mockedRelease)

	err := p.Down(context.Background())
	s.Require().ErrorIs(err, e)

	mockedRelease.AssertExpectations(s.T())
}

func (s *DestroyTestSuite) TestDestroyNoReleases() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))
	p.NewBody()

	err := p.Down(context.Background())
	s.Require().NoError(err)
}

func TestDestroyTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(DestroyTestSuite))
}
