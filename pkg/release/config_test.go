package release_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (s *ConfigTestSuite) TestConfigUniq() {
	r := release.NewConfig()
	r.NameF = "redis"
	r.NamespaceF = "test"

	s.Require().NoError(r.Uniq().Validate())
}

func (s *ConfigTestSuite) TestDependsOn() {
	r := release.NewConfig()

	r.NamespaceF = "testns"
	r.DependsOnF = []string{"bla", "blabla@testns", "blablabla@testtestns"}

	expected := []uniqname.UniqName{
		uniqname.UniqName("bla@testns"),
		uniqname.UniqName("blabla@testns"),
		uniqname.UniqName("blablabla@testtestns"),
	}
	s.Require().ElementsMatch(r.DependsOn(), expected)
}

func TestConfigTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ConfigTestSuite))
}
