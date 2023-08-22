package release_test

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

type ChartTestSuite struct {
	suite.Suite
}

func TestChartTestSuite(t *testing.T) {
	suite.Run(t, new(ChartTestSuite))
}

func (s *ChartTestSuite) SetupSuite() {
	var rs repo.Configs
	str := `
- name: bitnami
  url: https://charts.bitnami.com/bitnami
`
	err := yaml.Unmarshal([]byte(str), &rs)

	s.Require().NoError(err)
	s.Require().Len(rs, 1)

	s.Require().NoError(plan.SyncRepositories(context.Background(), rs))
}

func (s *ChartTestSuite) TestLocateChartLocal() {
	tmpDir := s.T().TempDir()

	rel := release.NewConfig()
	rel.ChartF.Name = filepath.Join(tmpDir, "blabla")

	c, err := rel.GetChart()
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "failed to locate chart")
	s.Require().Nil(c)
}

func (s *ChartTestSuite) TestLoadChartLocal() {
	tmpDir := s.T().TempDir()

	rel := release.NewConfig()
	rel.ChartF.Name = tmpDir

	c, err := rel.GetChart()
	s.Require().Error(err)
	s.Require().Contains(err.Error(), "failed to load chart")
	s.Require().Contains(err.Error(), "Chart.yaml file is missing")
	s.Require().Nil(c)
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

	s.Require().ErrorIs(err, release.ErrUnknownFormat)
}
