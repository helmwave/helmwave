//go:build integration

package action

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/databus23/helm-diff/v3/diff"
	"github.com/helmwave/helmwave/pkg/templater/sprig"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
	"github.com/urfave/cli/v2"
)

type DiffLiveTestSuite struct {
	suite.Suite

	ctx context.Context
}

//nolint:paralleltest // uses helm repository.yaml flock
func TestDiffLiveTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(DiffLiveTestSuite))
}

func (ts *DiffLiveTestSuite) SetupTest() {
	ts.ctx = tests.GetContext(ts.T())
}

func (ts *DiffLiveTestSuite) TestCmd() {
	s := &DiffLive{diff: &Diff{Options: &diff.Options{}}}
	cmd := s.Cmd()

	ts.Require().NotNil(cmd)
	ts.Require().NotEmpty(cmd.Name)
}

func (ts *DiffLiveTestSuite) TestRun() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		tpl:       filepath.Join(tests.Root, "02_helmwave.yml"),
		file:      filepath.Join(tests.Root, "02_helmwave.yml"),
		templater: &sprig.Templater{},
	}

	s := &Build{
		plandir:  tmpDir,
		tags:     cli.StringSlice{},
		yml:      y,
		diff:     &Diff{Options: &diff.Options{}},
		diffMode: DiffModeLive,
	}

	d := DiffLive{diff: s.diff, plandir: s.plandir}

	ts.Require().ErrorIs(d.Run(ts.ctx), os.ErrNotExist)
	ts.Require().NoError(s.Run(ts.ctx))
	ts.Require().NoError(d.Run(ts.ctx))
}
