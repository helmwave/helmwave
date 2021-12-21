//go:build ignore || integration

package action

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
	"github.com/urfave/cli/v2"
)

type DiffLiveTestSuite struct {
	suite.Suite
}

func (ts *DiffLiveTestSuite) TestImplementsAction() {
	ts.Require().Implements((*Action)(nil), &DiffLive{})
}

func (ts *DiffLiveTestSuite) TestRun() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		filepath.Join(tests.Root, "02_helmwave.yml"),
		filepath.Join(tests.Root, "02_helmwave.yml"),
	}

	s := &Build{
		plandir:  tmpDir,
		tags:     cli.StringSlice{},
		yml:      y,
		diff:     &Diff{},
		diffMode: diffModeLive,
	}

	d := DiffLive{diff: s.diff, plandir: s.plandir}

	ts.Require().ErrorIs(d.Run(), os.ErrNotExist)
	ts.Require().NoError(s.Run())
	ts.Require().NoError(d.Run())
}

func TestDiffLiveTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(DiffLiveTestSuite))
}
