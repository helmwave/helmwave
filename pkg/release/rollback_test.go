//go:build integration

package release_test

import (
	"context"
	"testing"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
)

type RollbackTestSuite struct {
	suite.Suite

	ctx context.Context
}

func TestRollbackTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RollbackTestSuite))
}

func (ts *RollbackTestSuite) SetupSuite() {
	ts.ctx = tests.GetContext(ts.T())
}

func (ts *RollbackTestSuite) TestNonExistingRollback() {
	rel := release.NewConfig()
	rel.NameF = "blabla"
	rel.NamespaceF = "blabla"

	err := rel.Rollback(ts.ctx, 1)

	ts.Require().ErrorContains(err, "failed to rollback release blabla@blabla:")
}
