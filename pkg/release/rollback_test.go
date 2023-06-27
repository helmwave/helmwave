//go:build ignore || integration

package release_test

import (
	"context"
	"testing"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/stretchr/testify/suite"
)

type RollbackTestSuite struct {
	suite.Suite
}

func (s *RollbackTestSuite) TestNonExistingRollback() {
	rel := release.NewConfig()
	rel.NameF = "blabla"
	rel.NamespaceF = "blabla"

	err := rel.Rollback(context.Background(), 1)

	s.Require().ErrorContains(err, "failed to rollback release blabla@blabla:")
}

func TestRollbackTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RollbackTestSuite))
}
