package action_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/action"
	"github.com/stretchr/testify/suite"
)

type RollbackTestSuite struct {
	suite.Suite
}

func TestRollbackTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RollbackTestSuite))
}

func (ts *RollbackTestSuite) TestCmd() {
	s := &action.Rollback{}
	cmd := s.Cmd()

	ts.Require().NotNil(cmd)
	ts.Require().NotEmpty(cmd.Name)
}
