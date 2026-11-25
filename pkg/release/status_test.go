//go:build integration

package release_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
)

type StatusTestSuite struct {
	suite.Suite
}

func (s *StatusTestSuite) TestNonExistingStatus() {
	tests.RequireCluster(s.T())
	rel := release.NewConfig()
	rel.NameF = "blabla"
	rel.NamespaceF = "blabla"

	_, err := rel.Status()

	s.Require().ErrorContains(err, "failed to get status of release blabla@blabla:")
}

func TestStatusTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(StatusTestSuite))
}
