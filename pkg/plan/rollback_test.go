package plan_test

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/kubedog"
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
)

type RollbackTestSuite struct {
	suite.Suite

	ctx context.Context
}

func TestRollbackTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RollbackTestSuite))
}

func (ts *RollbackTestSuite) SetupTest() {
	ts.ctx = tests.GetContext(ts.T())
}

func (ts *RollbackTestSuite) TestRollback() {
	tmpDir := ts.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))

	mockedRelease := plan.NewMockReleaseConfig(ts.T())
	mockedRelease.On("Rollback").Return(nil)

	p.SetReleases(mockedRelease)

	err := p.Rollback(ts.ctx, -1, &kubedog.Config{Enabled: false})
	ts.Require().NoError(err)

	mockedRelease.AssertExpectations(ts.T())
}

func (ts *RollbackTestSuite) TestRollbackError() {
	tmpDir := ts.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))

	mockedRelease := plan.NewMockReleaseConfig(ts.T())
	e := errors.New(ts.T().Name())
	mockedRelease.On("Rollback").Return(e)

	p.SetReleases(mockedRelease)

	err := p.Rollback(ts.ctx, -1, &kubedog.Config{Enabled: false})
	ts.Require().ErrorIs(err, e)

	mockedRelease.AssertExpectations(ts.T())
}

func (ts *RollbackTestSuite) TestRollbackNoReleases() {
	tmpDir := ts.T().TempDir()
	p := plan.New(filepath.Join(tmpDir, plan.Dir))
	p.NewBody()

	err := p.Rollback(ts.ctx, -1, &kubedog.Config{Enabled: false})
	ts.Require().NoError(err)
}
