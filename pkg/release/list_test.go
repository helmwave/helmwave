//go:build integration

package release_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
)

type ListTestSuite struct {
	suite.Suite
}

func (s *ListTestSuite) TestNonExistingList() {
	tests.RequireCluster(s.T())
	rel := release.NewConfig()
	rel.NameF = "blabla"
	rel.NamespaceF = "blabla"

	_, err := rel.List()

	s.Require().ErrorIs(err, release.ErrNotFound)
}

func TestListTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ListTestSuite))
}
