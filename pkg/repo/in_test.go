package repo_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/stretchr/testify/suite"
)

type InTestSuite struct {
	suite.Suite
}

func TestInTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(InTestSuite))
}

func (ts *InTestSuite) TestIndexOfName() {
	rep := repo.NewConfig()
	rep.Entry.Name = ts.T().Name()

	idx, found := repo.IndexOfName([]repo.Config{rep, rep, rep}, ts.T().Name())

	ts.Require().True(found)
	ts.Require().Equal(0, idx)
}

func (ts *InTestSuite) TestIndexOfNameNotFound() {
	rep := repo.NewConfig()
	rep.Entry.Name = ts.T().Name()

	_, found := repo.IndexOfName([]repo.Config{rep}, "")

	ts.Require().False(found)
}
