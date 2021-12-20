package release_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/stretchr/testify/suite"
)

type RepoTestSuite struct {
	suite.Suite
}

func (s *RepoTestSuite) TestRepoWithSlash() {
	const bitnami = "bitnami"
	r := release.NewConfig()
	r.ChartF = release.Chart{
		Name: bitnami + "/redis",
	}

	s.Require().Equal(bitnami, r.Repo())
}

func (s *RepoTestSuite) TestRepoWithoutSlash() {
	r := release.NewConfig()
	r.ChartF = release.Chart{
		Name: "api",
	}

	s.Require().Equal("api", r.Repo())
}

func TestRepoTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RepoTestSuite))
}
