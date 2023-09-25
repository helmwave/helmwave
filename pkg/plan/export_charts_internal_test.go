package plan

import (
	"errors"
	"net/url"
	"testing"

	"github.com/helmwave/go-fsimpl"
	"github.com/helmwave/go-fsimpl/filefs"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/stretchr/testify/suite"
)

type ExportChartsTestSuite struct {
	suite.Suite
}

func TestExportChartsTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ExportChartsTestSuite))
}

func (ts *BuildRepositoriesTestSuite) TestEmptyReleases() {
	p := New()
	p.NewBody()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: ts.T().TempDir()})

	err := p.exportCharts(baseFS.(fsimpl.CurrentPathFS), baseFS.(fsimpl.WriteableFS))

	ts.Require().NoError(err)
}

func (ts *BuildRepositoriesTestSuite) TestMultipleReleases() {
	p := New()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: ts.T().TempDir()})

	rel1 := &MockReleaseConfig{}
	rel1.On("Chart").Return(&release.Chart{})
	rel1.On("Uniq").Return(uniqname.UniqName("rel1@ns1"))
	rel2 := &MockReleaseConfig{}
	rel2.On("Chart").Return(&release.Chart{})
	rel2.On("Uniq").Return(uniqname.UniqName("rel2@ns2"))

	p.SetReleases(rel1, rel2)

	err := p.exportCharts(baseFS.(fsimpl.CurrentPathFS), baseFS.(fsimpl.WriteableFS))

	ts.Require().NoError(err)
	rel1.AssertExpectations(ts.T())
	rel2.AssertExpectations(ts.T())
}

func (ts *BuildRepositoriesTestSuite) TestError() {
	p := New()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: ts.T().TempDir()})

	rel1 := &MockReleaseConfig{}
	rel1.On("Chart").Return(&release.Chart{})
	rel1.On("Uniq").Return(uniqname.UniqName("rel1@ns1"))

	rel2 := &MockReleaseConfig{}
	errExpected := errors.New(ts.T().Name())
	rel2.On("DownloadChart").Return(errExpected)
	rel2.On("Chart").Return(&release.Chart{Name: "bitnami/blalba"})
	rel2.On("Uniq").Return(uniqname.UniqName("rel2@ns2"))

	p.SetReleases(rel1, rel2)

	err := p.exportCharts(baseFS.(fsimpl.CurrentPathFS), baseFS.(fsimpl.WriteableFS))

	ts.Require().ErrorIs(err, errExpected)
	rel1.AssertExpectations(ts.T())
	rel2.AssertExpectations(ts.T())
}
