package plan_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/plan"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type RollbackTestSuite struct {
	suite.Suite
}

func (s *RollbackTestSuite) TestRollback() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))

	mockedRelease := &plan.MockReleaseConfig{}
	mockedRelease.On("Rollback").Return(nil)
	mockedRelease.On("Logger").Return(log.WithField("test", s.T().Name()))

	p.SetReleases(mockedRelease)

	err := p.Rollback(context.Background(), -1)
	s.Require().NoError(err)

	mockedRelease.AssertExpectations(s.T())
}

func (s *RollbackTestSuite) TestRollbackError() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))

	mockedRelease := &plan.MockReleaseConfig{}
	e := errors.New(s.T().Name())
	mockedRelease.On("Rollback").Return(e)
	mockedRelease.On("Logger").Return(log.WithField("test", s.T().Name()))

	p.SetReleases(mockedRelease)

	err := p.Rollback(context.Background(), -1)
	s.Require().ErrorIs(err, e)

	mockedRelease.AssertExpectations(s.T())
}

func (s *RollbackTestSuite) TestRollbackNoReleases() {
	tmpDir := s.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))
	p.NewBody()

	err := p.Rollback(context.Background(), -1)
	s.Require().NoError(err)
}

func TestRollbackTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RollbackTestSuite))
}
