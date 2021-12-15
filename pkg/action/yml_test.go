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

func (s *YmlTestSuite) TestRenderEnv() {
	tmpDir := s.T().TempDir()
	y := &Yml{
		filepath.Join(tests.Root, "01_helmwave.yml.tpl"),
		filepath.Join(tmpDir, "01_helmwave.yml"),
	}

	value := "test01"
	s.T().Setenv("PROJECT_NAME", value)
	s.T().Setenv("NAMESPACE", value)

	s.Require().NoError(y.Run())

	b, err := plan.NewBody(y.file)
	s.Require().NoError(err)

	s.Require().Equal(value, b.Project)
	s.Require().Len(b.Releases, 1)
	s.Require().Equal(value, b.Releases[0].Namespace)
}

//nolint:paralleltest // cannot parallel because of setenv
func TestYmlTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(YmlTestSuite))
}
