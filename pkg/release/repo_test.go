package release

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type RepoTestSuite struct {
	suite.Suite
}

func (s *RepoTestSuite) TestRepoWithSlash() {
	const bitnami = "bitnami"
	r := &config{ChartF: Chart{
		Name: bitnami + "/redis",
	}}

	s.Require().Equal(bitnami, r.Repo())
}

func (s *RepoTestSuite) TestRepoWithoutSlash() {
	r := &config{ChartF: Chart{
		Name: "api",
	}}

	s.Require().Equal("api", r.Repo())
}

func TestRepoTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RepoTestSuite))
}
