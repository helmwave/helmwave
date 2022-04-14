package plan_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/plan"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"helm.sh/helm/v3/pkg/chart"
	helmRelease "helm.sh/helm/v3/pkg/release"
)

type ListTestSuite struct {
	suite.Suite
}

func (s *ListTestSuite) TestList() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))

	mockedRelease := &plan.MockReleaseConfig{}
	r := &helmRelease.Release{
		Info: &helmRelease.Info{},
		Chart: &chart.Chart{
			Metadata: &chart.Metadata{},
		},
	}
	mockedRelease.On("List").Return(r, nil)

	p.SetReleases(mockedRelease)

	err := p.List()
	s.Require().NoError(err)

	mockedRelease.AssertExpectations(s.T())
}

// TestListError tests that List method should just skip releases that fail List method.
func (s *ListTestSuite) TestListError() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))

	mockedRelease := &plan.MockReleaseConfig{}
	mockedRelease.On("List").Return(&helmRelease.Release{}, errors.New(s.T().Name()))
	mockedRelease.On("Logger").Return(log.WithField("test", s.T().Name()))

	p.SetReleases(mockedRelease)

	err := p.List()
	s.Require().NoError(err)

	mockedRelease.AssertExpectations(s.T())
}

func (s *ListTestSuite) TestListNoReleases() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))
	p.NewBody()

	err := p.List()
	s.Require().NoError(err)
}

func TestListTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ListTestSuite))
}
