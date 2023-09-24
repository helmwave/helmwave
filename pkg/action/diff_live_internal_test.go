//go:build integration

package action

import (
	"context"
	"os"
	"testing"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
	"github.com/urfave/cli/v2"
)

type DiffLiveTestSuite struct {
	suite.Suite
}

//nolint:paralleltest // uses helm repository.yaml flock
func TestDiffLiveTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(DiffLiveTestSuite))
}

func (ts *DiffLiveTestSuite) TestCmd() {
	s := &DiffLive{}
	cmd := s.Cmd()

	ts.Require().NotNil(cmd)
	ts.Require().NotEmpty(cmd.Name)
}

func (ts *DiffLiveTestSuite) TestRun() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		templater: template.TemplaterSprig,
	}

	s := &Build{
		tags:     cli.StringSlice{},
		yml:      y,
		diff:     &Diff{},
		diffMode: DiffModeLive,
	}
	createGenericFS(&s.yml.srcFS, tests.Root, "02_helmwave.yml")
	createGenericFS(&s.yml.destFS, tests.Root, "02_helmwave.yml")
	createGenericFS(&s.planFS, tmpDir)

	d := DiffLive{diff: s.diff, planFS: s.planFS.(plan.PlanImportFS)}
	createGenericFS(&d.planFS)

	ts.Require().ErrorIs(d.Run(context.Background()), os.ErrNotExist)
	ts.Require().NoError(s.Run(context.Background()))
	ts.Require().NoError(d.Run(context.Background()))
}
