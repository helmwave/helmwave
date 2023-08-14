package release_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
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

func (s *ChartTestSuite) TestDryRun() {
	rel := release.NewConfig()

	s.Require().False(rel.IsDryRun())
	rel.DryRun(true)
	s.Require().True(rel.IsDryRun())
}

func (s *ChartTestSuite) TestUnmarshalYAMLString() {
	var rs release.Chart
	str := "blabla"
	err := yaml.Unmarshal([]byte(str), &rs)

	s.Require().NoError(err)
	s.Require().Equal(rs.Name, str)
}

func (s *ChartTestSuite) TestUnmarshalYAMLMapping() {
	var rs release.Chart
	str := `
name: blabla
version: 1.2.3
`
	err := yaml.Unmarshal([]byte(str), &rs)

	s.Require().NoError(err)
	s.Require().Equal(rs.Name, "blabla")
	s.Require().Equal(rs.Version, "1.2.3")
}

func (s *ChartTestSuite) TestUnmarshalYAMLInvalid() {
	var rs release.Chart
	str := "[1, 2, 3]"
	err := yaml.Unmarshal([]byte(str), &rs)

	s.Require().ErrorContains(err, "unknown format")
}

func TestConfigTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ConfigTestSuite))
}
