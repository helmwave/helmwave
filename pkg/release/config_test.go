package release_test

import (
	"slices"
	"testing"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func TestConfigTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ConfigTestSuite))
}

func (s *ConfigTestSuite) TestConfigUniq() {
	r := release.NewConfig()
	r.NameF = "redis"
	r.NamespaceF = "test"
	r.KubeContextF = "ctx"

	s.Require().NoError(r.Uniq().Validate())
}

func (s *ConfigTestSuite) TestConfigUniqTags() {
	r := release.NewConfig()

	r.BuildAfterUnmarshal()

	s.Require().True(slices.Contains(r.TagsF, r.Uniq().String()))
}

func (s *ConfigTestSuite) TestConfigInvalidUniq() {
	r := release.NewConfig()
	r.NameF = "redis"
	r.NamespaceF = ""

	s.Require().Error(r.Uniq().Validate())
}

func (s *ConfigTestSuite) TestDependsOn() {
	r := release.NewConfig()

	r.NamespaceF = "testns"
	r.KubeContextF = "testctx"
	r.DependsOnF = []*release.DependsOnReference{
		{Name: "bla"},
		{Name: "blabla@testns"},
		{Name: "blablabla@testtestns"},
		{Name: "---=-=-==-@kk;'[["},
	}

	r.BuildAfterUnmarshal(r)

	expected := []*release.DependsOnReference{
		{Name: "bla@testns@testctx"},
		{Name: "blabla@testns@testctx"},
		{Name: "blablabla@testtestns@testctx"},
	}
	s.Require().ElementsMatch(r.DependsOn(), expected)
}

func (s *ConfigTestSuite) TestDryRun() {
	rel := release.NewConfig()

	s.Require().False(rel.IsDryRun())
	rel.DryRun(true)
	s.Require().True(rel.IsDryRun())
}
