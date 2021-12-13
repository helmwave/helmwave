//go:build ignore || integration

package action

import (
	"os"
	"testing"

	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
	"github.com/urfave/cli/v2"
)

type DiffLiveTestSuite struct {
	suite.Suite
}

func (ts *DiffLiveTestSuite) TestRun() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		tests.Root + "01_helmwave.yml.tpl",
		tmpDir + "02_helmwave.yml",
	}

	s := &Build{
		plandir:  tmpDir,
		tags:     cli.StringSlice{},
		autoYml:  true,
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
	// t.Parallel()
	suite.Run(t, new(DiffLiveTestSuite))
}
