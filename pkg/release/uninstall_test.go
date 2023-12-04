//go:build integration

package release_test

import (
	"context"
	"testing"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/stretchr/testify/suite"
)

type UninstallTestSuite struct {
	suite.Suite
}

func (s *UninstallTestSuite) TestNonExistingUninstall() {
	rel := release.NewConfig()
	rel.NameF = "blabla"
	rel.NamespaceF = "blabla"

	_, err := rel.Uninstall(context.Background())

	s.Require().NoError(err)
}

func TestUninstallTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(UninstallTestSuite))
}
