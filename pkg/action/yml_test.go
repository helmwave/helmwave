//go:build ignore || integration

package action

import (
	"os"
	"testing"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
)

type YmlTestSuite struct {
	suite.Suite
}

func (s *YmlTestSuite) TestRenderEnv() {
	defer clean()

	y := &Yml{
		tests.Root + "01_helmwave.yml.tpl",
		tests.Root + "01_helmwave.yml",
	}
	defer os.Remove(y.file)

	value := "Test01"
	_ = os.Setenv("PROJECT_NAME", value)
	_ = os.Setenv("NAMESPACE", value)

	template.SetConfig(&template.Config{
		Gomplate: template.GomplateConfig{
			Enabled: false,
		},
	})

	err := y.Run()
	s.Require().NoError(err)

	b, err := plan.NewBody(y.file)
	s.Require().NoError(err)

	s.Require().Equal(value, b.Project)
	s.Require().Len(b.Releases, 1)
	s.Require().Equal(value, b.Releases[0].Namespace)
}

func TestYmlTestSuite(t *testing.T) {
	//t.Parallel()
	suite.Run(t, new(YmlTestSuite))
}
