package action_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/action"
	"github.com/stretchr/testify/suite"
)

type StatusTestSuite struct {
	suite.Suite
}

func (ts *StatusTestSuite) TestImplementsAction() {
	ts.Require().Implements((*action.Action)(nil), &action.Status{})
}

func TestStatusTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(StatusTestSuite))
}
