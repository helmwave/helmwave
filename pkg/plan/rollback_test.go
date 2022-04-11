package plan_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/stretchr/testify/suite"
)

type RollbackTestSuite struct {
	suite.Suite
}

func (s *RollbackTestSuite) TestRollback() {
	tmpDir := s.T().TempDir()
	p, err := plan.New(filepath.Join(tmpDir, plan.Dir))
	s.Require().NoError(err)

	mockedRelease := &plan.MockReleaseConfig{}
	mockedRelease.On("Name").Return(s.T().Name())
	mockedRelease.On("Namespace").Return(s.T().Name())
	mockedRelease.On("Uniq").Return()
	mockedRelease.On("Rollback").Return(nil)

	p.SetReleases(mockedRelease)

	err = p.Rollback()
	s.Require().NoError(err)

	mockedRelease.AssertExpectations(s.T())
}

func (s *RollbackTestSuite) TestRollbackError() {
	tmpDir := s.T().TempDir()
	p, err := plan.New(filepath.Join(tmpDir, plan.Dir))
	s.Require().NoError(err)

	mockedRelease := &plan.MockReleaseConfig{}
	e := errors.New(s.T().Name())
	mockedRelease.On("Name").Return(s.T().Name())
	mockedRelease.On("Namespace").Return(s.T().Name())
	mockedRelease.On("Uniq").Return()
	mockedRelease.On("Rollback").Return(e)

	p.SetReleases(mockedRelease)

	err = p.Rollback()
	s.Require().ErrorIs(err, e)

	mockedRelease.AssertExpectations(s.T())
}

func (s *RollbackTestSuite) TestRollbackNoReleases() {
	tmpDir := s.T().TempDir()
	p, err := plan.New(filepath.Join(tmpDir, plan.Dir))
	s.Require().NoError(err)

	p.NewBody()

	err = p.Rollback()
	s.Require().NoError(err)
}

func TestRollbackTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RollbackTestSuite))
}
