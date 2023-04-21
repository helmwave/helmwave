//go:build ignore || integration

package action

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/kubedog"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
	"github.com/urfave/cli/v2"
)

type DownTestSuite struct {
	suite.Suite
}

func (ts *DownTestSuite) TestImplementsAction() {
	ts.Require().Implements((*Action)(nil), &Down{})
}

func (ts *DownTestSuite) TestRun() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		tpl:       filepath.Join(tests.Root, "02_helmwave.yml"),
		file:      filepath.Join(tests.Root, "02_helmwave.yml"),
		templater: "sprig",
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
	ts.Require().ErrorIs(d.Run(context.Background()), os.ErrNotExist, "down should fail before build")
	ts.Require().NoError(s.Run(context.Background()))

	u := &Up{
		build: s,
		dog:   &kubedog.Config{},
	}

	ts.Require().NoError(u.Run(context.Background()))
	ts.Require().NoError(d.Run(context.Background()))
}

// nolintlint:paralleltest // uses helm repository.yaml flock
func TestDownTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(DownTestSuite))
}
