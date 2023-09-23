package plan

import (
	"errors"
	"io/fs"
	"net/url"
	"os"
	"testing"

	"github.com/helmwave/go-fsimpl"
	"github.com/helmwave/go-fsimpl/filefs"
	"github.com/stretchr/testify/suite"
)

type BuildChartsTestSuite struct {
	suite.Suite
}

func TestBuildChartsTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(BuildChartsTestSuite))
}

func (ts *BuildRepositoriesTestSuite) TestEmptyReleases() {
	p := New(".")
	p.NewBody()

	wd, _ := os.Getwd()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: wd})
	err := p.buildCharts(baseFS.(fs.StatFS), baseFS.(fsimpl.WriteableFS))

	ts.Require().NoError(err)
}

func (ts *BuildRepositoriesTestSuite) TestMultipleReleases() {
	p := New(".")

	rel1 := &MockReleaseConfig{}
	rel1.On("DownloadChart").Return(nil)
	rel2 := &MockReleaseConfig{}
	rel2.On("DownloadChart").Return(nil)

	p.SetReleases(rel1, rel2)

	wd, _ := os.Getwd()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: wd})
	err := p.buildCharts(baseFS.(fs.StatFS), baseFS.(fsimpl.WriteableFS))

	ts.Require().NoError(err)
	rel1.AssertExpectations(ts.T())
	rel2.AssertExpectations(ts.T())
}

func (ts *BuildRepositoriesTestSuite) TestError() {
	p := New(".")

	rel1 := &MockReleaseConfig{}
	rel1.On("DownloadChart").Return(nil)
	rel2 := &MockReleaseConfig{}
	errExpected := errors.New(ts.T().Name())
	rel2.On("DownloadChart").Return(errExpected)

	p.SetReleases(rel1, rel2)

	wd, _ := os.Getwd()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: wd})
	err := p.buildCharts(baseFS.(fs.StatFS), baseFS.(fsimpl.WriteableFS))

	ts.Require().ErrorIs(err, errExpected)
	rel1.AssertExpectations(ts.T())
	rel2.AssertExpectations(ts.T())
}
