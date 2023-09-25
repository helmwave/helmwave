package release_test

import (
	"context"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"testing"

	"github.com/helmwave/go-fsimpl"
	"github.com/helmwave/go-fsimpl/filefs"
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

func (ts *ChartTestSuite) SetupSuite() {
	var rs repo.Configs
	str := `
- name: bitnami
  url: https://charts.bitnami.com/bitnami
`
	err := yaml.Unmarshal([]byte(str), &rs)

	ts.Require().NoError(err)
	ts.Require().Len(rs, 1)

	ts.Require().NoError(plan.SyncRepositories(context.Background(), rs))
}

func (ts *ChartTestSuite) TestLocateChartLocal() {
	tmpDir := ts.T().TempDir()

	rel := release.NewConfig()
	rel.ChartF.Name = filepath.Join(tmpDir, "blabla")

	wd, _ := os.Getwd()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: wd})
	c, err := rel.GetChart(baseFS.(fs.StatFS))
	ts.Require().Error(err)
	ts.Require().Contains(err.Error(), "failed to locate chart")
	ts.Require().Nil(c)
}

func (ts *ChartTestSuite) TestLoadChartLocal() {
	tmpDir := ts.T().TempDir()

	rel := release.NewConfig()
	rel.ChartF.Name = tmpDir

	wd, _ := os.Getwd()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: wd})
	c, err := rel.GetChart(baseFS.(fs.StatFS))
	ts.Require().Error(err)
	ts.Require().Contains(err.Error(), "failed to load chart")
	ts.Require().Contains(err.Error(), "Chart.yaml file is missing")
	ts.Require().Nil(c)
}

func (ts *ChartTestSuite) TestUnmarshalYAMLString() {
	var rs release.Chart
	str := "blabla"
	err := yaml.Unmarshal([]byte(str), &rs)

	ts.Require().NoError(err)
	ts.Require().Equal(rs.Name, str)
}

func (ts *ChartTestSuite) TestUnmarshalYAMLMapping() {
	var rs release.Chart
	str := `
name: blabla
version: 1.2.3
`
	err := yaml.Unmarshal([]byte(str), &rs)

	ts.Require().NoError(err)
	ts.Require().Equal(rs.Name, "blabla")
	ts.Require().Equal(rs.Version, "1.2.3")
}

func (ts *ChartTestSuite) TestUnmarshalYAMLInvalid() {
	var rs release.Chart
	str := "[1, 2, 3]"
	err := yaml.Unmarshal([]byte(str), &rs)

	ts.Require().ErrorIs(err, release.ErrUnknownFormat)
}

func (ts *ChartTestSuite) TestIsRemote() {
	c := &release.Chart{Name: "/nonexisting"}
	wd, _ := os.Getwd()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: wd})

	ts.Require().True(c.IsRemote(baseFS.(fs.StatFS)))

	c.Name = ts.T().TempDir()

	ts.Require().False(c.IsRemote(baseFS.(fs.StatFS)))
}

func (ts *ChartTestSuite) TestChartDepsUpdRemote() {
	rel := release.NewConfig()
	rel.SetChartName("bitnami/redis")

	wd, _ := os.Getwd()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: wd})
	err := rel.ChartDepsUpd(baseFS.(fsimpl.CurrentPathFS))

	ts.Require().NoError(err)
}

func (ts *ChartTestSuite) TestSkipChartDepsUpd() {
	rel := release.NewConfig()
	rel.ChartF.Name = ts.T().TempDir()
	rel.ChartF.SkipDependencyUpdate = true

	wd, _ := os.Getwd()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: wd})
	err := rel.ChartDepsUpd(baseFS.(fsimpl.CurrentPathFS))

	ts.Require().NoError(err)
}

func (ts *ChartTestSuite) TestChartDepsUpdInvalid() {
	rel := release.NewConfig()
	rel.ChartF.Name = ts.T().TempDir()

	wd, _ := os.Getwd()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: wd})
	err := rel.ChartDepsUpd(baseFS.(fsimpl.CurrentPathFS))

	ts.Require().ErrorContains(err, "Chart.yaml file is missing")
}

func (ts *ChartTestSuite) TestDownloadChartRemote() {
	rel := release.NewConfig()
	rel.SetChartName("bitnami/redis")

	wd, _ := os.Getwd()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: wd})
	err := rel.DownloadChart(baseFS.(fsimpl.CurrentPathFS), baseFS.(fsimpl.WriteableFS), ts.T().TempDir())

	ts.Require().NoError(err)
}

func (ts *ChartTestSuite) TestDownloadChartLocal() {
	rel := release.NewConfig()

	wd, _ := os.Getwd()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: wd})
	err := rel.DownloadChart(baseFS.(fsimpl.CurrentPathFS), baseFS.(fsimpl.WriteableFS), ts.T().TempDir())

	ts.Require().NoError(err)
}
