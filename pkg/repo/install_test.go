package repo_test

import (
	"context"
	"testing"

	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
	helm "helm.sh/helm/v3/pkg/cli"
	helmRepo "helm.sh/helm/v3/pkg/repo"
)

type InstallTestSuite struct {
	suite.Suite

	ctx context.Context
}

func TestInstallTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(InstallTestSuite))
}

func (ts *InstallTestSuite) SetupTest() {
	ts.ctx = tests.GetContext(ts.T())
}

func (ts *InstallTestSuite) TestInstallNonExisting() {
	rep := repo.NewConfig()
	settings := &helm.EnvSettings{}
	f := &helmRepo.File{}

	err := rep.Install(ts.ctx, settings, f)

	ts.Require().NoError(err)
	ts.Require().Contains(f.Repositories, &rep.Entry)
}

func (ts *InstallTestSuite) TestInstallExistingSame() {
	rep := repo.NewConfig()
	settings := &helm.EnvSettings{}
	f := &helmRepo.File{
		Repositories: []*helmRepo.Entry{&rep.Entry},
	}

	err := rep.Install(ts.ctx, settings, f)

	ts.Require().NoError(err)
	ts.Require().Contains(f.Repositories, &rep.Entry)
}

func (ts *InstallTestSuite) TestInstallExistingNotSame() {
	rep1 := repo.NewConfig()
	rep2 := repo.NewConfig()
	rep2.Entry.URL = ts.T().Name()
	settings := &helm.EnvSettings{}
	f := &helmRepo.File{
		Repositories: []*helmRepo.Entry{&rep2.Entry},
	}

	err := rep1.Install(ts.ctx, settings, f)

	var e *repo.DuplicateError
	ts.Require().ErrorAs(err, &e)
	ts.Equal(rep1.Name(), e.Name)

	ts.Contains(f.Repositories, &rep2.Entry)
	ts.NotContains(f.Repositories, &rep1.Entry)
}
