package plan

import (
	"errors"
	"testing"

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

	err := p.buildCharts()

	ts.Require().NoError(err)
}

func (ts *BuildRepositoriesTestSuite) TestMultipleReleases() {
	p := New(".")

	rel1 := NewMockReleaseConfig(ts.T())
	rel1.On("DownloadChart").Return(nil)
	rel2 := NewMockReleaseConfig(ts.T())
	rel2.On("DownloadChart").Return(nil)

	p.SetReleases(rel1, rel2)

	err := p.buildCharts()

	ts.Require().NoError(err)
	rel1.AssertExpectations(ts.T())
	rel2.AssertExpectations(ts.T())
}

func (ts *BuildRepositoriesTestSuite) TestError() {
	p := New(".")

	rel1 := NewMockReleaseConfig(ts.T())
	rel1.On("DownloadChart").Return(nil)
	rel2 := NewMockReleaseConfig(ts.T())
	errExpected := errors.New(ts.T().Name())
	rel2.On("DownloadChart").Return(errExpected)

	p.SetReleases(rel1, rel2)

	err := p.buildCharts()

	ts.Require().ErrorIs(err, errExpected)
	rel1.AssertExpectations(ts.T())
	rel2.AssertExpectations(ts.T())
}
