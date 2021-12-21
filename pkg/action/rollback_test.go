package action_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/action"
	"github.com/stretchr/testify/suite"
)

type RollbackTestSuite struct {
	suite.Suite
}

func (ts *RollbackTestSuite) TestImplementsAction() {
	ts.Require().Implements((*action.Action)(nil), &action.Rollback{})
}

func TestRollbackTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RollbackTestSuite))
}
