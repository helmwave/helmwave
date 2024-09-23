//go:build integration

package action

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/kubedog"
	"github.com/helmwave/helmwave/pkg/templater/sprig"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
	"github.com/urfave/cli/v2"
)

type DownTestSuite struct {
	suite.Suite

	ctx context.Context
}

//nolint:paralleltest // uses helm repository.yaml flock
func TestDownTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(DownTestSuite))
}

func (ts *DownTestSuite) SetupTest() {
	ts.ctx = tests.GetContext(ts.T())
}

func (ts *DownTestSuite) TestCmd() {
	s := &Down{}
	cmd := s.Cmd()

	ts.Require().NotNil(cmd)
	ts.Require().NotEmpty(cmd.Name)
}

func (ts *DownTestSuite) TestRun() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		tpl:       filepath.Join(tests.Root, "02_helmwave.yml"),
		file:      filepath.Join(tests.Root, "02_helmwave.yml"),
		templater: &sprig.Templater{},
	}

	s := &Build{
		plandir: tmpDir,
		tags:    cli.StringSlice{},
		autoYml: true,
		yml:     y,
	}

	d := Down{
		build: s,
	}
	ts.Require().ErrorIs(d.Run(ts.ctx), os.ErrNotExist, "down should fail before build")
	ts.Require().NoError(s.Run(ts.ctx))

	u := &Up{
		build: s,
		dog:   &kubedog.Config{},
	}

	ts.Require().NoError(u.Run(ts.ctx))
	ts.Require().NoError(d.Run(ts.ctx))
}

func (ts *DownTestSuite) TestIdempotency() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		tpl:       filepath.Join(tests.Root, "02_helmwave.yml"),
		file:      filepath.Join(tests.Root, "02_helmwave.yml"),
		templater: &sprig.Templater{},
	}

	s := &Build{
		plandir: tmpDir,
		tags:    cli.StringSlice{},
		autoYml: true,
		yml:     y,
	}

	u := &Up{
		build: s,
		dog:   &kubedog.Config{},
	}
	d := Down{
		build: s,
	}

	ctx, cancel := context.WithCancel(ts.ctx)
	ts.T().Cleanup(cancel)

	ts.Require().NoError(s.Run(ctx))
	ts.Require().NoError(u.Run(ctx))
	ts.Require().NoError(d.Run(ctx))
	ts.Require().NoError(d.Run(ctx))
}
