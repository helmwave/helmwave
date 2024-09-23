//go:build integration

package action

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/templater/sprig"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
)

type YmlTestSuite struct {
	suite.Suite

	ctx context.Context
}

//nolint:paralleltest // can't parallel because of setenv
func TestYmlTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(YmlTestSuite))
}

func (ts *YmlTestSuite) SetupTest() {
	ts.ctx = tests.GetContext(ts.T())
}

func (ts *YmlTestSuite) TestCmd() {
	s := &Yml{}
	cmd := s.Cmd()

	ts.Require().NotNil(cmd)
	ts.Require().NotEmpty(cmd.Name)
}

func (ts *YmlTestSuite) TestRenderEnv() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		tpl:       filepath.Join(tests.Root, "01_helmwave.yml.tpl"),
		file:      filepath.Join(tmpDir, "01_helmwave.yml"),
		templater: &sprig.Templater{},
	}

	value := "test01"
	ts.T().Setenv("NAMESPACE", value)
	ts.T().Setenv("PROJECT_NAME", value)

	ts.Require().NoError(y.Run(ts.ctx))

	b, err := plan.NewBody(ts.ctx, y.file, true)
	ts.Require().NoError(err)

	ts.Require().Equal(value, b.Project)
	ts.Require().Len(b.Releases, 1)
	ts.Require().Equal(value, b.Releases[0].Namespace())
}
