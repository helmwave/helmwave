package plan_test

import (
	"context"
	"testing"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
)

type NewTestSuite struct {
	suite.Suite

	ctx context.Context
}

func TestNewTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(NewTestSuite))
}

func (ts *NewTestSuite) SetupTest() {
	ts.ctx = tests.GetContext(ts.T())
}

func (ts *NewTestSuite) TestNew() {
	dir := "/proc/1/bla"
	p := plan.New(dir)

	ts.Require().NotNil(p)
	ts.Require().False(p.IsExist())
	ts.Require().False(p.IsManifestExist())
}

func (ts *NewTestSuite) TestNewAndImportError() {
	_, err := plan.NewAndImport(ts.ctx, "/proc/1/blabla")

	ts.Require().Error(err)
	ts.Require().ErrorContains(err, "failed to read plan file")
}

func (ts *NewTestSuite) TestLogger() {
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

	ts.Require().NotNil(logger)

	ts.Require().Contains(logger.Data, "releases")
	ts.Require().Equal([]string{uniq.String()}, logger.Data["releases"])

	ts.Require().Contains(logger.Data, "repositories")
	ts.Require().Equal([]string{repoName}, logger.Data["repositories"])

	rel.AssertExpectations(ts.T())
	repo.AssertExpectations(ts.T())
}

func (ts *NewTestSuite) TestJSONSchema() {
	schema := plan.GenSchema()

	ts.Require().NotNil(schema)

	ts.NotNil(schema.Properties.GetPair("repositories"))
	ts.NotNil(schema.Properties.GetPair("registries"))
	ts.NotNil(schema.Properties.GetPair("releases"))
	ts.NotNil(schema.Properties.GetPair("lifecycle"))
}
