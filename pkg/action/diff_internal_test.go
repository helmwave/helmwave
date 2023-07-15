//go:build ignore || integration

package action

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type DiffTestSuite struct {
	suite.Suite
}

//nolintlint:paralleltest // uses helm repository.yaml flock
func TestDiffTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(DiffTestSuite))
}

func (ts *DiffTestSuite) TestCmd() {
	s := &Diff{}
	cmd := s.Cmd()

	ts.Require().NotNil(cmd)
	ts.Require().NotEmpty(cmd.Name)
}
