package action

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type RollbackTestSuite struct {
	suite.Suite
}

func (ts *RollbackTestSuite) TestImplementsAction() {
	ts.Require().Implements((*Action)(nil), &Rollback{})
}

func TestRollbackTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RollbackTestSuite))
}
