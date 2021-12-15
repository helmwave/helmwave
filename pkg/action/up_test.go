//go:build ignore || integration

package action

import (
	"strings"
	"testing"

	"github.com/helmwave/helmwave/pkg/kubedog"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
	"github.com/urfave/cli/v2"
)

type UpTestSuite struct {
	suite.Suite
}

func (ts *UpTestSuite) TestAutoBuild() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		tests.Root + "01_helmwave.yml.tpl",
		tmpDir + "02_helmwave.yml",
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
	ts.T().Setenv("PROJECT_NAME", value)
	ts.T().Setenv("NAMESPACE", value)

	ts.Require().NoError(u.Run())
}

//nolint:paralleltest // cannot parallel because of setenv
func TestUpTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(UpTestSuite))
}
