//go:build integration

package action

import (
	"context"
	"github.com/helmwave/helmwave/pkg/templater"
	"path/filepath"
	"strings"
	"testing"

	"github.com/helmwave/helmwave/pkg/kubedog"
	"github.com/helmwave/helmwave/pkg/templater/sprig"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
	"github.com/urfave/cli/v2"
)

type UpTestSuite struct {
	suite.Suite

	ctx context.Context
}

//nolint:paralleltest // can't parallel because of setenv
func TestUpTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(UpTestSuite))
}

func (ts *UpTestSuite) SetupTest() {
	ts.ctx = tests.GetContext(ts.T())
}

func (ts *UpTestSuite) TestCmd() {
	s := &Up{}
	cmd := s.Cmd()

	ts.Require().NotNil(cmd)
	ts.Require().NotEmpty(cmd.Name)
}

func (ts *UpTestSuite) TestAutoBuild() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		tpl:       filepath.Join(tests.Root, "01_helmwave.yml.tpl"),
		file:      filepath.Join(tmpDir, "02_helmwave.yml"),
		templater: &sprig.Templater{},
	}

	u := &Up{
		build: &Build{
			plandir: tmpDir,
			tags:    cli.StringSlice{},
			autoYml: true,
			yml:     y,
		},
		dog:       &kubedog.Config{},
		autoBuild: true,
	}

	value := strings.ToLower(strings.ReplaceAll(ts.T().Name(), "/", ""))
	ts.T().Setenv("NAMESPACE", value)

	ts.Require().NoError(u.Run(ts.ctx))
}

func (ts *UpTestSuite) TestPrometheusMonitors() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		tpl:       filepath.Join(tests.Root, "20_helmwave.yml"),
		file:      filepath.Join(tests.Root, "20_helmwave.yml"),
		templater: templater.Default,
	}

	s := &Up{
		build: &Build{
			plandir: tmpDir,
			tags:    cli.StringSlice{},
			autoYml: false,
			yml:     y,
		},
		dog:       &kubedog.Config{},
		autoBuild: true,
	}

	ts.Require().NoError(s.Run(ts.ctx))
}
