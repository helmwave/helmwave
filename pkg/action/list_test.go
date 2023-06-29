package action_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/action"
	"github.com/stretchr/testify/suite"
)

type ListTestSuite struct {
	suite.Suite
}

func TestListTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ListTestSuite))
}

func (ts *ListTestSuite) TestCmd() {
	s := &action.List{}
	cmd := s.Cmd()

	ts.Require().NotNil(cmd)
	ts.Require().NotEmpty(cmd.Name)
}
