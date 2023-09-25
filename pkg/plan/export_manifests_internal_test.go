package plan

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"testing"

	"github.com/helmwave/go-fsimpl/filefs"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	helmRelease "helm.sh/helm/v3/pkg/release"
)

type ExportManifestsTestSuite struct {
	suite.Suite
}

func TestExportManifestsTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ExportManifestsTestSuite))
}

func (ts *ExportManifestsTestSuite) TestEmptyReleases() {
	p := New()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: ts.T().TempDir()})
	p.NewBody()

	err := p.exportManifests(context.Background(), baseFS.(ExportFS))

	ts.Require().NoError(err)
}

func (ts *ExportManifestsTestSuite) TestMultipleReleases() {
	p := New()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: ts.T().TempDir()})

	rel1 := &MockReleaseConfig{}
	u1 := uniqname.UniqName("redis1@defaultblabla")
	rel1.On("Logger").Return(log.WithField("test", ts.T().Name()))
	rel1.On("ChartDepsUpd").Return(nil)
	rel1.On("DryRun").Return()
	rel1.On("Sync").Return(&helmRelease.Release{}, nil)
	rel1.On("HooksDisabled").Return(false)
	rel1.On("Uniq").Return(u1)

	rel2 := &MockReleaseConfig{}
	u2 := uniqname.UniqName("redis2@defaultblabla")
	rel2.On("Logger").Return(log.WithField("test", ts.T().Name()))
	rel2.On("ChartDepsUpd").Return(nil)
	rel2.On("DryRun").Return()
	rel2.On("Sync").Return(&helmRelease.Release{}, nil)
	rel2.On("HooksDisabled").Return(false)
	rel2.On("Uniq").Return(u2)

	p.SetReleases(rel1, rel2)

	err := p.exportManifests(context.Background(), baseFS.(ExportFS))

	ts.Require().NoError(err)
	ts.Require().Len(p.manifests, 2)
	ts.Require().Contains(p.manifests, u1)
	ts.Require().Contains(p.manifests, u2)
	ts.Require().Equal(p.manifests[u1], "")
	ts.Require().Equal(p.manifests[u2], "")

	rel1.AssertExpectations(ts.T())
	rel2.AssertExpectations(ts.T())
}

func (ts *ExportManifestsTestSuite) TestChartDepsUpdError() {
	p := New()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: ts.T().TempDir()})

	rel := &MockReleaseConfig{}
	uniq := uniqname.UniqName("redis1@defaultblabla")
	errExpected := errors.New(ts.T().Name())
	rel.On("Logger").Return(log.WithField("test", ts.T().Name()))
	rel.On("ChartDepsUpd").Return(errExpected)
	rel.On("DryRun").Return()
	rel.On("Sync").Return(&helmRelease.Release{}, nil)
	rel.On("HooksDisabled").Return(false)
	rel.On("Uniq").Return(uniq)

	p.SetReleases(rel)

	err := p.exportManifests(context.Background(), baseFS.(ExportFS))

	ts.Require().NoError(err)
	ts.Require().Len(p.manifests, 1)
	ts.Require().Contains(p.manifests, uniq)
	ts.Require().Equal(p.manifests[uniq], "")

	rel.AssertExpectations(ts.T())
}

func (ts *ExportManifestsTestSuite) TestSyncError() {
	p := New()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: ts.T().TempDir()})

	rel := &MockReleaseConfig{}
	errExpected := errors.New(ts.T().Name())
	rel.On("Logger").Return(log.WithField("test", ts.T().Name()))
	rel.On("ChartDepsUpd").Return(nil)
	rel.On("DryRun").Return()
	rel.On("Sync").Return(&helmRelease.Release{}, errExpected)

	p.SetReleases(rel)

	err := p.exportManifests(context.Background(), baseFS.(ExportFS))

	ts.Require().ErrorIs(err, errExpected)

	rel.AssertExpectations(ts.T())
}

func (ts *ExportManifestsTestSuite) TestDisabledHooks() {
	p := New()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: ts.T().TempDir()})

	rel := &MockReleaseConfig{}
	uniq := uniqname.UniqName("redis1@defaultblabla")
	rel.On("Logger").Return(log.WithField("test", ts.T().Name()))
	rel.On("ChartDepsUpd").Return(nil)
	rel.On("DryRun").Return()
	rel.On("Sync").Return(&helmRelease.Release{
		Manifest: ts.T().Name(),
		Hooks:    []*helmRelease.Hook{{}},
	}, nil)
	rel.On("HooksDisabled").Return(true)
	rel.On("Uniq").Return(uniq)

	p.SetReleases(rel)

	err := p.exportManifests(context.Background(), baseFS.(ExportFS))

	ts.Require().NoError(err)
	ts.Require().Len(p.manifests, 1)
	ts.Require().Contains(p.manifests, uniq)
	ts.Require().Equal(p.manifests[uniq], ts.T().Name())

	rel.AssertExpectations(ts.T())
}

func (ts *ExportManifestsTestSuite) TestEnabledHooks() {
	p := New()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: ts.T().TempDir()})

	rel := &MockReleaseConfig{}
	uniq := uniqname.UniqName("redis1@defaultblabla")
	rel.On("Logger").Return(log.WithField("test", ts.T().Name()))
	rel.On("ChartDepsUpd").Return(nil)
	rel.On("DryRun").Return()
	rel.On("Sync").Return(&helmRelease.Release{
		Manifest: ts.T().Name(),
		Hooks: []*helmRelease.Hook{
			{
				Path:     ts.T().Name(),
				Manifest: ts.T().Name(),
			},
		},
	}, nil)
	rel.On("HooksDisabled").Return(false)
	rel.On("Uniq").Return(uniq)

	p.SetReleases(rel)

	err := p.exportManifests(context.Background(), baseFS.(ExportFS))

	ts.Require().NoError(err)
	ts.Require().Len(p.manifests, 1)
	ts.Require().Contains(p.manifests, uniq)
	ts.Require().Equal(p.manifests[uniq], fmt.Sprintf("%[1]s---\n# Source: %[1]s\n%[1]s\n", ts.T().Name()))

	rel.AssertExpectations(ts.T())
}
