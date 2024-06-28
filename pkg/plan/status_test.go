package plan_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/stretchr/testify/suite"
	"helm.sh/helm/v3/pkg/chart"
	helmRelease "helm.sh/helm/v3/pkg/release"
)

type StatusTestSuite struct {
	suite.Suite
}

func (s *StatusTestSuite) TestStatusByName() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))

	mockedRelease := plan.NewMockReleaseConfig(s.T())
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("Uniq").Return()
	mockedRelease.On("KubeContext").Return("")

	r := &helmRelease.Release{
		Info: &helmRelease.Info{},
		Chart: &chart.Chart{
			Metadata: &chart.Metadata{},
		},
	}
	mockedRelease.On("Status").Return(r, nil)

	p.SetReleases(mockedRelease)

	err := p.Status(mockedRelease.Uniq().String())
	s.Require().NoError(err)

	err = p.Status()
	s.Require().NoError(err)

	mockedRelease.AssertExpectations(s.T())
}

// TestStatusFailedRelease tests that Status method should just skip releases that fail Status method.
func (s *StatusTestSuite) TestStatusFailedRelease() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))

	mockedRelease := plan.NewMockReleaseConfig(s.T())
	mockedRelease.On("Status").Return(&helmRelease.Release{}, errors.New(s.T().Name()))

	p.SetReleases(mockedRelease)

	err := p.Status()
	s.Require().NoError(err)

	mockedRelease.AssertExpectations(s.T())
}

func (s *StatusTestSuite) TestStatusNoReleases() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))
	p.NewBody()

	err := p.Status()
	s.Require().NoError(err)
}

func TestStatusTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(StatusTestSuite))
}
