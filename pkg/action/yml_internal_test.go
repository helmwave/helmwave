//go:build integration

package action

import (
	"context"
	"testing"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
)

type YmlTestSuite struct {
	suite.Suite
}

//nolint:paralleltest // can't parallel because of setenv
func TestYmlTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(YmlTestSuite))
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
		templater: template.TemplaterSprig,
	}
	createGenericFS(&y.srcFS, tests.Root, "01_helmwave.yml.tpl")
	createGenericFS(&y.destFS, tmpDir, "01_helmwave.yml")

	value := "test01"
	ts.T().Setenv("NAMESPACE", value)
	ts.T().Setenv("PROJECT_NAME", value)

	ts.Require().NoError(y.Run(context.Background()))

	b, err := plan.NewBody(context.Background(), y.destFS, true)
	ts.Require().NoError(err)

	ts.Require().Equal(value, b.Project)
	ts.Require().Len(b.Releases, 1)
	ts.Require().Equal(value, b.Releases[0].Namespace())
}
