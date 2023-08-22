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

func TestConfigTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ConfigTestSuite))
}

func (s *ConfigTestSuite) TestConfigUniq() {
	r := release.NewConfig()
	r.NameF = "redis"
	r.NamespaceF = "test"

	s.Require().NoError(r.Uniq().Validate())
}

func (s *ConfigTestSuite) TestConfigInvalidUniq() {
	r := release.NewConfig()
	r.NameF = "redis"
	r.NamespaceF = ""

	s.Require().ErrorIs(r.Uniq().Validate(), uniqname.ValidationError{})
}

func (s *ConfigTestSuite) TestDependsOn() {
	r := release.NewConfig()

	r.NamespaceF = "testns"
	r.DependsOnF = []*release.DependsOnReference{
		{Name: "bla"},
		{Name: "blabla@testns"},
		{Name: "blablabla@testtestns"},
		{Name: "---=-=-==-@kk;'[["},
	}

	r.BuildAfterUnmarshal(r)

	expected := []*release.DependsOnReference{
		{Name: "bla@testns"},
		{Name: "blabla@testns"},
		{Name: "blablabla@testtestns"},
	}
	s.Require().ElementsMatch(r.DependsOn(), expected)
}

func (s *ConfigTestSuite) TestDryRun() {
	rel := release.NewConfig()

	s.Require().False(rel.IsDryRun())
	rel.DryRun(true)
	s.Require().True(rel.IsDryRun())
}
