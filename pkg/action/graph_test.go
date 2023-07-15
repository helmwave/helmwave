package action_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/action"
	"github.com/stretchr/testify/suite"
)

type GraphTestSuite struct {
	suite.Suite
}

func TestGraphTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(GraphTestSuite))
}

func (ts *GraphTestSuite) TestCmd() {
	s := &action.Graph{}
	cmd := s.Cmd()

	ts.Require().NotNil(cmd)
	ts.Require().NotEmpty(cmd.Name)
}
