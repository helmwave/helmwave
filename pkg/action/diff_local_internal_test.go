//go:build integration

package action

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/helmwave/helmwave/tests"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/urfave/cli/v2"
)

var buf bytes.Buffer

type DiffLocalTestSuite struct {
	suite.Suite
}

//nolint:paralleltest // we capture output for global logger and uses helm repository.yaml flock
func TestDiffLocalTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(DiffLocalTestSuite))
}

func (ts *DiffLocalTestSuite) SetupTest() {
	log.StandardLogger().SetOutput(&buf)
}

func (ts *DiffLocalTestSuite) TearDownTest() {
	log.StandardLogger().SetOutput(os.Stderr)
}

func (ts *DiffLocalTestSuite) TestCmd() {
	s := &DiffLocal{}
	cmd := s.Cmd()

	ts.Require().NotNil(cmd)
	ts.Require().NotEmpty(cmd.Name)
}

func (ts *DiffLocalTestSuite) TestRun() {
	s1 := &Build{
		tags: cli.StringSlice{},
		yml: &Yml{
			templater: template.TemplaterSprig,
		},
		diff:     &Diff{},
		diffMode: DiffModeLive,
	}
	createGenericFS(&s1.yml.srcFS, tests.Root, "02_helmwave.yml")
	createGenericFS(&s1.yml.destFS, tests.Root, "02_helmwave.yml")
	createGenericFS(&s1.planFS, ts.T().TempDir())

	s2 := &Build{
		tags: cli.StringSlice{},
		yml: &Yml{
			templater: template.TemplaterSprig,
		},
		diff:     &Diff{},
		diffMode: DiffModeLive,
	}
	createGenericFS(&s2.yml.srcFS, tests.Root, "03_helmwave.yml")
	createGenericFS(&s2.yml.destFS, tests.Root, "03_helmwave.yml")
	createGenericFS(&s2.planFS, ts.T().TempDir())

	d := DiffLocal{diff: s1.diff, plan1FS: s1.planFS.(plan.ImportFS), plan2FS: s2.planFS.(plan.ImportFS)}

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
