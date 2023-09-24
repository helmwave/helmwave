package plan_test

import (
	"context"
	"errors"
	"testing"

	"github.com/helmwave/helmwave/pkg/kubedog"
	"github.com/helmwave/helmwave/pkg/plan"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type RollbackTestSuite struct {
	suite.Suite
}

func (s *RollbackTestSuite) TestRollback() {
	p := plan.New()

	mockedRelease := &plan.MockReleaseConfig{}
	mockedRelease.On("Rollback").Return(nil)
	mockedRelease.On("Logger").Return(log.WithField("test", s.T().Name()))

	p.SetReleases(mockedRelease)

	err := p.Rollback(context.Background(), -1, &kubedog.Config{Enabled: false})
	s.Require().NoError(err)

	mockedRelease.AssertExpectations(s.T())
}

func (s *RollbackTestSuite) TestRollbackError() {
	p := plan.New()

	mockedRelease := &plan.MockReleaseConfig{}
	e := errors.New(s.T().Name())
	mockedRelease.On("Rollback").Return(e)
	mockedRelease.On("Logger").Return(log.WithField("test", s.T().Name()))

	p.SetReleases(mockedRelease)

	err := p.Rollback(context.Background(), -1, &kubedog.Config{Enabled: false})
	s.Require().ErrorIs(err, e)

	mockedRelease.AssertExpectations(s.T())
}

func (s *RollbackTestSuite) TestRollbackNoReleases() {
	p := plan.New()
	p.NewBody()

	err := p.Rollback(context.Background(), -1, &kubedog.Config{Enabled: false})
	s.Require().NoError(err)
}

func TestRollbackTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RollbackTestSuite))
}
