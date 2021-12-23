//go:build ignore || integration

package action

import (
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/plan"
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
		templater: "sprig",
	}

	value := "test01"
	ts.T().Setenv("PROJECT_NAME", value)
	ts.T().Setenv("NAMESPACE", value)

	ts.Require().NoError(y.Run())

	b, err := plan.NewBody(y.file)
	ts.Require().NoError(err)

	ts.Require().Equal(value, b.Project)
	ts.Require().Len(b.Releases, 1)
	ts.Require().Equal(value, b.Releases[0].Namespace())
}

//nolint:paralleltest // cannot parallel because of setenv
func TestYmlTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(YmlTestSuite))
}
