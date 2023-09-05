package repo_test

import (
	"context"
	"testing"

	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/stretchr/testify/suite"
	helm "helm.sh/helm/v3/pkg/cli"
	helmRepo "helm.sh/helm/v3/pkg/repo"
)

type InstallTestSuite struct {
	suite.Suite
}

func TestInstallTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(InstallTestSuite))
}

func (ts *InstallTestSuite) TestInstallNonExisting() {
	rep := repo.NewConfig()
	settings := &helm.EnvSettings{}
	f := &helmRepo.File{}

	err := rep.Install(context.Background(), settings, f)

	ts.Require().NoError(err)
	ts.Require().Contains(f.Repositories, &rep.Entry)
}

func (ts *InstallTestSuite) TestInstallExistingSame() {
	rep := repo.NewConfig()
	settings := &helm.EnvSettings{}
	f := &helmRepo.File{
		Repositories: []*helmRepo.Entry{&rep.Entry},
	}

	err := rep.Install(context.Background(), settings, f)

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

	err := rep1.Install(context.Background(), settings, f)

	ts.Require().ErrorIs(err, repo.DuplicateError{})
	ts.Require().Contains(f.Repositories, &rep2.Entry)
	ts.Require().NotContains(f.Repositories, &rep1.Entry)
}
