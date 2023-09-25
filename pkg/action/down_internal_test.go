//go:build integration

package action

import (
	"context"
	"os"
	"testing"

	"github.com/helmwave/helmwave/pkg/kubedog"
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
	"github.com/urfave/cli/v2"
)

type DownTestSuite struct {
	suite.Suite
}

//nolint:paralleltest // uses helm repository.yaml flock
func TestDownTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(DownTestSuite))
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
		templater: template.TemplaterSprig,
	}

	s := &Build{
		tags:    cli.StringSlice{},
		autoYml: true,
		yml:     y,
	}
	createGenericFS(&s.yml.srcFS, tests.Root, "02_helmwave.yml")
	createGenericFS(&s.yml.destFS, tests.Root, "02_helmwave.yml")
	createGenericFS(&s.planFS, tmpDir)
	createGenericFS(&s.contextFS, tmpDir)

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
