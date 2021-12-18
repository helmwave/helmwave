//go:build ignore || integration

package action

import (
	"bytes"
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

func (s *DiffPlanTestSuite) SetupTest() {
	log.StandardLogger().SetOutput(&buf)
}

func (s *DiffPlanTestSuite) TearDownTest() {
	log.StandardLogger().SetOutput(os.Stderr)
}

func (ts *DiffPlanTestSuite) TestRun() {
	s1 := &Build{
		plandir: ts.T().TempDir(),
		tags:    cli.StringSlice{},
		yml: &Yml{
			filepath.Join(tests.Root, "02_helmwave.yml"),
			filepath.Join(tests.Root, "02_helmwave.yml"),
		},
		diff:     &Diff{},
		diffMode: diffModeLive,
	}

	s2 := &Build{
		plandir: ts.T().TempDir(),
		tags:    cli.StringSlice{},
		yml: &Yml{
			filepath.Join(tests.Root, "03_helmwave.yml"),
			filepath.Join(tests.Root, "03_helmwave.yml"),
		},
		diff:     &Diff{},
		diffMode: diffModeLive,
	}

	d := DiffLocalPlan{diff: s1.diff, plandir1: s1.plandir, plandir2: s2.plandir}

	ts.Require().ErrorIs(d.Run(), os.ErrNotExist)
	ts.Require().NoError(s1.Run())
	ts.Require().ErrorIs(d.Run(), os.ErrNotExist)
	ts.Require().NoError(s2.Run())

	buf.Reset()
	ts.Require().NoError(d.Run())

	output := buf.String()
	buf.Reset()

	ts.Require().Contains(output, "nginx, Deployment (apps) has been added")
	ts.Require().Contains(output, "memcached-a-redis, Secret (v1) has been removed")
}

//nolint:paralleltest // we capture output for global logger
func TestDiffPlanTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(DiffPlanTestSuite))
}
