//go:build integration

package action

import (
	"testing"

	"github.com/databus23/helm-diff/v3/diff"
	"github.com/stretchr/testify/suite"
)

type DiffTestSuite struct {
	suite.Suite
}

func TestDiffTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(DiffTestSuite))
}

func (ts *DiffTestSuite) TestCmd() {
	s := &Diff{Options: &diff.Options{}}
	cmd := s.Cmd()

	ts.Require().NotNil(cmd)
	ts.Require().NotEmpty(cmd.Name)
}
