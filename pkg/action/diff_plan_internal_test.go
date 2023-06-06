//go:build ignore || integration

package action

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/tests"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/urfave/cli/v2"
)

var buf bytes.Buffer

type DiffPlanTestSuite struct {
	suite.Suite
}

func (ts *DiffPlanTestSuite) SetupTest() {
	log.StandardLogger().SetOutput(&buf)
}

func (ts *DiffPlanTestSuite) TearDownTest() {
	log.StandardLogger().SetOutput(os.Stderr)
}

func (ts *DiffPlanTestSuite) TestImplementsAction() {
	ts.Require().Implements((*Action)(nil), &DiffLocalPlan{})
}

func (ts *DiffPlanTestSuite) TestRun() {
	s1 := &Build{
		plandir: ts.T().TempDir(),
		tags:    cli.StringSlice{},
		yml: &Yml{
			tpl:       filepath.Join(tests.Root, "02_helmwave.yml"),
			file:      filepath.Join(tests.Root, "02_helmwave.yml"),
			templater: "sprig",
		},
		diff:     &Diff{},
		diffMode: DiffModeLive,
	}

	s2 := &Build{
		plandir: ts.T().TempDir(),
		tags:    cli.StringSlice{},
		yml: &Yml{
			tpl:       filepath.Join(tests.Root, "03_helmwave.yml"),
			file:      filepath.Join(tests.Root, "03_helmwave.yml"),
			templater: "sprig",
		},
		diff:     &Diff{},
		diffMode: DiffModeLive,
	}

	d := DiffLocalPlan{diff: s1.diff, plandir1: s1.plandir, plandir2: s2.plandir}

	ts.Require().ErrorIs(d.Run(context.Background()), os.ErrNotExist)
	ts.Require().NoError(s1.Run(context.Background()))
	ts.Require().ErrorIs(d.Run(context.Background()), os.ErrNotExist)
	ts.Require().NoError(s2.Run(context.Background()))

	buf.Reset()
	ts.Require().NoError(d.Run(context.Background()))

	output := buf.String()
	buf.Reset()

	ts.Require().Contains(output, "nginx, Deployment (apps) has been added")
	ts.Require().Contains(output, "memcached-a-redis, Secret (v1) has been removed")
}

//nolintlint:paralleltest // we capture output for global logger and uses helm repository.yaml flock
func TestDiffPlanTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(DiffPlanTestSuite))
}
