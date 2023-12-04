package hooks_test

import (
	"context"
	"errors"
	"os/exec"
	"testing"

	"github.com/helmwave/helmwave/pkg/hooks"
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
}

func TestLifecycleRunTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(LifecycleRunTestSuite))
}

func (s *LifecycleRunTestSuite) TestRunNoError() {
	hook := &MockHook{}
	hook.On("Run").Return(nil)

	lifecycle := hooks.Lifecycle{PreBuild: []hooks.Hook{hook}}
	ctx := context.Background()

	err := lifecycle.RunPreBuild(ctx)

	s.Require().NoError(err)
	hook.AssertExpectations(s.T())
}

func (s *LifecycleRunTestSuite) TestRunError() {
	hook := &MockHook{}
	errExpected := errors.New("test 123")
	hook.On("Run").Return(errExpected)
	hook.On("Log").Return(log.WithField("test", s.T().Name()))

	lifecycle := hooks.Lifecycle{PreBuild: []hooks.Hook{hook}}
	ctx := context.Background()

	err := lifecycle.RunPreBuild(ctx)

	s.Require().ErrorIs(err, errExpected)
	hook.AssertExpectations(s.T())
}

func (s *LifecycleRunTestSuite) TestEmpty() {
	lifecycle := hooks.Lifecycle{}
	ctx := context.Background()

	s.Require().NoError(lifecycle.RunPreBuild(ctx))
	s.Require().NoError(lifecycle.RunPostBuild(ctx))
	s.Require().NoError(lifecycle.RunPreUp(ctx))
	s.Require().NoError(lifecycle.RunPostUp(ctx))
	s.Require().NoError(lifecycle.RunPreDown(ctx))
	s.Require().NoError(lifecycle.RunPostDown(ctx))
	s.Require().NoError(lifecycle.RunPreRollback(ctx))
	s.Require().NoError(lifecycle.RunPostRollback(ctx))
}

type HookRunTestSuite struct {
	suite.Suite
}

func TestHookRunTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(HookRunTestSuite))
}

func (s *HookRunTestSuite) TestRunCanceledContext() {
	hook := hooks.NewHook()
	hook.Cmd = "id"

	lifecycle := hooks.Lifecycle{PreBuild: []hooks.Hook{hook}}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := lifecycle.RunPreBuild(ctx)

	s.Require().ErrorIs(err, context.Canceled)
}

func (s *HookRunTestSuite) TestRunNoError() {
	hook := hooks.NewHook()
	hook.Cmd = "id"

	lifecycle := hooks.Lifecycle{PreBuild: []hooks.Hook{hook}}
	ctx := context.Background()

	err := lifecycle.RunPreBuild(ctx)

	s.Require().NoError(err)
}

func (s *HookRunTestSuite) TestRunWrongCommand() {
	hook := hooks.NewHook()
	hook.Cmd = "id 123"

	lifecycle := hooks.Lifecycle{PreBuild: []hooks.Hook{hook}}
	ctx := context.Background()

	err := lifecycle.RunPreBuild(ctx)

	s.Require().ErrorIs(err, exec.ErrNotFound)
}
