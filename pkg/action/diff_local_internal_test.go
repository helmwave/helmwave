//go:build integration

package action

import (
	"bytes"
	"context"
	"github.com/databus23/helm-diff/v3/diff"
	"os"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/template"
	"github.com/helmwave/helmwave/tests"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/urfave/cli/v2"
)

var buf bytes.Buffer

type DiffLocalTestSuite struct {
	suite.Suite

	ctx context.Context
}

//nolint:paralleltest // we capture output for global logger and uses helm repository.yaml flock
func TestDiffLocalTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(DiffLocalTestSuite))
}

func (ts *DiffLocalTestSuite) SetupTest() {
	ts.ctx = tests.GetContext(ts.T())

	log.StandardLogger().SetOutput(&buf)
	ts.T().Cleanup(func() {
		log.StandardLogger().SetOutput(os.Stderr)
	})
}

func (ts *DiffLocalTestSuite) TestCmd() {
	s := &DiffLocal{}
	cmd := s.Cmd()

	ts.Require().NotNil(cmd)
	ts.Require().NotEmpty(cmd.Name)
}

func (ts *DiffLocalTestSuite) TestRun() {
	s1 := &Build{
		plandir: ts.T().TempDir(),
		tags:    cli.StringSlice{},
		yml: &Yml{
			tpl:       filepath.Join(tests.Root, "02_helmwave.yml"),
			file:      filepath.Join(tests.Root, "02_helmwave.yml"),
			templater: template.TemplaterSprig,
		},
		diff:     &Diff{Options: &diff.Options{}},
		diffMode: DiffModeLive,
	}

	s2 := &Build{
		plandir: ts.T().TempDir(),
		tags:    cli.StringSlice{},
		yml: &Yml{
			tpl:       filepath.Join(tests.Root, "03_helmwave.yml"),
			file:      filepath.Join(tests.Root, "03_helmwave.yml"),
			templater: template.TemplaterSprig,
		},
		diff:     &Diff{Options: &diff.Options{}},
		diffMode: DiffModeLive,
	}

	d := DiffLocal{diff: s1.diff, plandir1: s1.plandir, plandir2: s2.plandir}

	ts.Require().ErrorIs(d.Run(ts.ctx), os.ErrNotExist)
	ts.Require().NoError(s1.Run(ts.ctx))
	ts.Require().ErrorIs(d.Run(ts.ctx), os.ErrNotExist)
	ts.Require().NoError(s2.Run(ts.ctx))

	buf.Reset()
	ts.Require().NoError(d.Run(ts.ctx))

	output := buf.String()
	buf.Reset()

	ts.Require().Contains(output, "nginx, Deployment (apps) has been added")
	ts.Require().Contains(output, "memcached-a-redis, Secret (v1) has been removed")
}
