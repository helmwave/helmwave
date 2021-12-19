package action

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ListTestSuite struct {
	suite.Suite
}

func (ts *ListTestSuite) TestImplementsAction() {
	ts.Require().Implements((*Action)(nil), &List{})
}

func TestListTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ListTestSuite))
}
