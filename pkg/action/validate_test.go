package action

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ValidateTestSuite struct {
	suite.Suite
}

func (ts *ValidateTestSuite) TestImplementsAction() {
	ts.Require().Implements((*Action)(nil), &Validate{})
}

func TestValidateTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ValidateTestSuite))
}
