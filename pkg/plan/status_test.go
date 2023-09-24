package plan_test

import (
	"errors"
	"testing"

	"github.com/helmwave/helmwave/pkg/plan"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"helm.sh/helm/v3/pkg/chart"
	helmRelease "helm.sh/helm/v3/pkg/release"
)

type StatusTestSuite struct {
	suite.Suite
}

func (s *StatusTestSuite) TestStatusByName() {
	p := plan.New()

	mockedRelease := &plan.MockReleaseConfig{}
	mockedRelease.On("Name").Return("redis")
	mockedRelease.On("Namespace").Return("defaultblabla")
	mockedRelease.On("Uniq").Return()
	mockedRelease.On("Logger").Return(log.WithField("test", s.T().Name()))
	r := &helmRelease.Release{
		Info: &helmRelease.Info{},
		Chart: &chart.Chart{
			Metadata: &chart.Metadata{},
		},
	}
	mockedRelease.On("Status").Return(r, nil)

	p.SetReleases(mockedRelease)

	err := p.Status(string(mockedRelease.Uniq()))
	s.Require().NoError(err)

	err = p.Status()
	s.Require().NoError(err)

	mockedRelease.AssertExpectations(s.T())
}

// TestStatusFailedRelease tests that Status method should just skip releases that fail Status method.
func (s *StatusTestSuite) TestStatusFailedRelease() {
	p := plan.New()

	mockedRelease := &plan.MockReleaseConfig{}
	mockedRelease.On("Status").Return(&helmRelease.Release{}, errors.New(s.T().Name()))
	mockedRelease.On("Logger").Return(log.WithField("test", s.T().Name()))

	p.SetReleases(mockedRelease)

	err := p.Status()
	s.Require().NoError(err)

	mockedRelease.AssertExpectations(s.T())
}

func (s *StatusTestSuite) TestStatusNoReleases() {
	p := plan.New()
	p.NewBody()

	err := p.Status()
	s.Require().NoError(err)
}

func TestStatusTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(StatusTestSuite))
}
