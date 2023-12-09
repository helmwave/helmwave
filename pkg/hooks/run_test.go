package hooks_test

import (
	"context"
	"errors"
	"os/exec"
	"testing"

	"github.com/helmwave/helmwave/pkg/hooks"
	"github.com/helmwave/helmwave/tests"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var _ hooks.Hook = (*MockHook)(nil)

type MockHook struct {
	mock.Mock
}

func (m *MockHook) Run(_ context.Context) error {
	return m.Called().Error(0)
}

func (m *MockHook) Log() *log.Entry {
	return m.Called().Get(0).(*log.Entry)
}

type LifecycleRunTestSuite struct {
	suite.Suite

	ctx context.Context
}

func TestLifecycleRunTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(LifecycleRunTestSuite))
}

func (ts *LifecycleRunTestSuite) SetupTest() {
	ts.ctx = tests.GetContext(ts.T())
}

func (ts *LifecycleRunTestSuite) TestRunNoError() {
	hook := &MockHook{}
	hook.On("Run").Return(nil)

	lifecycle := hooks.Lifecycle{PreBuild: []hooks.Hook{hook}}

	err := lifecycle.RunPreBuild(ts.ctx)

	ts.Require().NoError(err)
	hook.AssertExpectations(ts.T())
}

func (ts *LifecycleRunTestSuite) TestRunError() {
	hook := &MockHook{}
	errExpected := errors.New("test 123")
	hook.On("Run").Return(errExpected)
	hook.On("Log").Return(log.WithField("test", ts.T().Name()))

	lifecycle := hooks.Lifecycle{PreBuild: []hooks.Hook{hook}}

	err := lifecycle.RunPreBuild(ts.ctx)

	ts.Require().ErrorIs(err, errExpected)
	hook.AssertExpectations(ts.T())
}

func (ts *LifecycleRunTestSuite) TestEmpty() {
	lifecycle := hooks.Lifecycle{}

	ts.Require().NoError(lifecycle.RunPreBuild(ts.ctx))
	ts.Require().NoError(lifecycle.RunPostBuild(ts.ctx))
	ts.Require().NoError(lifecycle.RunPreUp(ts.ctx))
	ts.Require().NoError(lifecycle.RunPostUp(ts.ctx))
	ts.Require().NoError(lifecycle.RunPreDown(ts.ctx))
	ts.Require().NoError(lifecycle.RunPostDown(ts.ctx))
	ts.Require().NoError(lifecycle.RunPreRollback(ts.ctx))
	ts.Require().NoError(lifecycle.RunPostRollback(ts.ctx))
}

type HookRunTestSuite struct {
	suite.Suite

	ctx context.Context
}

func TestHookRunTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(HookRunTestSuite))
}

func (ts *HookRunTestSuite) SetupTest() {
	ts.ctx = tests.GetContext(ts.T())
}

func (ts *HookRunTestSuite) TestRunCanceledContext() {
	hook := hooks.NewHook()
	hook.Cmd = "id"

	lifecycle := hooks.Lifecycle{PreBuild: []hooks.Hook{hook}}
	ctx, cancel := context.WithCancel(ts.ctx)
	cancel()

	err := lifecycle.RunPreBuild(ctx)

	ts.Require().ErrorIs(err, context.Canceled)
}

func (ts *HookRunTestSuite) TestRunNoError() {
	hook := hooks.NewHook()
	hook.Cmd = "id"

	lifecycle := hooks.Lifecycle{PreBuild: []hooks.Hook{hook}}

	err := lifecycle.RunPreBuild(ts.ctx)

	ts.Require().NoError(err)
}

func (ts *HookRunTestSuite) TestRunWrongCommand() {
	hook := hooks.NewHook()
	hook.Cmd = "id 123"

	lifecycle := hooks.Lifecycle{PreBuild: []hooks.Hook{hook}}

	err := lifecycle.RunPreBuild(ts.ctx)

	ts.Require().ErrorIs(err, exec.ErrNotFound)
}
