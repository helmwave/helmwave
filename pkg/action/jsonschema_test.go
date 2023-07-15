package action_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/action"
	"github.com/stretchr/testify/suite"
)

type GenSchemaTestSuite struct {
	suite.Suite
}

func TestGenSchemaTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(GenSchemaTestSuite))
}

func (ts *GenSchemaTestSuite) TestCmd() {
	s := &action.GenSchema{}
	cmd := s.Cmd()

	ts.Require().NotNil(cmd)
	ts.Require().NotEmpty(cmd.Name)
}
