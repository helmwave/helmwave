package plan_test

import (
	"context"
	"testing"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/stretchr/testify/suite"
)

type NewTestSuite struct {
	suite.Suite
}

func TestNewTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(NewTestSuite))
}

func (s *NewTestSuite) TestNew() {
	dir := "/proc/1/bla"
	p := plan.New(dir)

	s.Require().NotNil(p)
	s.Require().False(p.IsExist())
	s.Require().False(p.IsManifestExist())
}

func (s *NewTestSuite) TestNewAndImportError() {
	_, err := plan.NewAndImport(context.Background(), "/proc/1/blabla")

	s.Require().Error(err)
	s.Require().ErrorContains(err, "failed to read plan file")
}

func (s *NewTestSuite) TestLogger() {
	p := plan.New(".")
	body := p.NewBody()

	rel := &plan.MockReleaseConfig{}
	uniq := uniqname.UniqName("blabla@namespace")
	rel.On("Uniq").Return(uniq)

	repo := &plan.MockRepositoryConfig{}
	repoName := "blarepo"
	repo.On("Name").Return(repoName)

	body.Releases = append(body.Releases, rel)
	body.Repositories = append(body.Repositories, repo)

	logger := p.Logger()

	s.Require().NotNil(logger)

	s.Require().Contains(logger.Data, "releases")
	s.Require().Equal([]string{uniq.String()}, logger.Data["releases"])

	s.Require().Contains(logger.Data, "repositories")
	s.Require().Equal([]string{repoName}, logger.Data["repositories"])

	rel.AssertExpectations(s.T())
	repo.AssertExpectations(s.T())
}

func (s *NewTestSuite) TestJSONSchema() {
	schema := plan.GenSchema()

	s.Require().NotNil(schema)

	keys := schema.Properties.Keys()
	s.Require().Contains(keys, "repositories")
	s.Require().Contains(keys, "registries")
	s.Require().Contains(keys, "releases")
	s.Require().Contains(keys, "lifecycle")
}
