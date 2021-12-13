//go:build ignore || integration

package action

import (
	"os"
	"strings"
	"testing"

	"github.com/helmwave/helmwave/pkg/kubedog"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
	"github.com/urfave/cli/v2"
)

type DownTestSuite struct {
	suite.Suite
}

func (ts *DownTestSuite) TestRun() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		tests.Root + "01_helmwave.yml.tpl",
		tmpDir + "02_helmwave.yml",
	}

	s := &Build{
		plandir: tmpDir,
		tags:    cli.StringSlice{},
		autoYml: true,
		yml:     y,
	}

	value := strings.ToLower(strings.ReplaceAll(ts.T().Name(), "/", ""))
	ts.T().Setenv("PROJECT_NAME", value)
	ts.T().Setenv("NAMESPACE", value)

	d := Down{plandir: s.plandir}
	ts.Require().ErrorIs(d.Run(), os.ErrNotExist, "down should fail before build")
	ts.Require().NoError(s.Run())

	u := &Up{
		build: s,
		dog:   &kubedog.Config{},
	}

	ts.Require().NoError(u.Run())
	ts.Require().NoError(d.Run())
}

func TestDownTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(DownTestSuite))
}
