//go:build integration

package release_test

import (
	"context"
	"testing"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
)

type UninstallTestSuite struct {
	suite.Suite

	ctx context.Context
}

func TestUninstallTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UninstallTestSuite))
}

func (ts *UninstallTestSuite) SetupTest() {
	ts.ctx = tests.GetContext(ts.T())
}

func (ts *UninstallTestSuite) TestNonExistingUninstall() {
	rel := release.NewConfig()
	rel.NameF = "blabla"
	rel.NamespaceF = "blabla"

	_, err := rel.Uninstall(ts.ctx)

	ts.Require().NoError(err)
}
