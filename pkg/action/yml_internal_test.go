//go:build ignore || integration

package action

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
)

type YmlTestSuite struct {
	suite.Suite
}

func (ts *YmlTestSuite) TestImplementsAction() {
	ts.Require().Implements((*Action)(nil), &Yml{})
}

func (ts *YmlTestSuite) TestRenderEnv() {
	tmpDir := ts.T().TempDir()
	y := &Yml{
		tpl:       filepath.Join(tests.Root, "01_helmwave.yml.tpl"),
		file:      filepath.Join(tmpDir, "01_helmwave.yml"),
		templater: template.TemplaterSprig,
	}

	value := "test01"
	ts.T().Setenv("NAMESPACE", value)
	ts.T().Setenv("PROJECT_NAME", value)

	ts.Require().NoError(y.Run(context.Background()))

	b, err := plan.NewBody(context.Background(), y.file)
	ts.Require().NoError(err)

	ts.Require().Equal(value, b.Project)
	ts.Require().Len(b.Releases, 1)
	ts.Require().Equal(value, b.Releases[0].Namespace())
}

//nolintlint:paralleltest // can't parallel because of setenv
func TestYmlTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(YmlTestSuite))
}
