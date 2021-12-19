package action

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type StatusTestSuite struct {
	suite.Suite
}

func (ts *StatusTestSuite) TestImplementsAction() {
	ts.Require().Implements((*Action)(nil), &Status{})
}

func TestStatusTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(StatusTestSuite))
}
