package release_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/release/uniqname"

	"github.com/stretchr/testify/suite"
)

type ValidateTestSuite struct {
	suite.Suite
}

func (s *ValidateTestSuite) TestEmptyName() {
	rel := release.NewConfig()
	rel.NameF = ""

	s.Require().ErrorIs(rel.Validate(), release.ErrNameEmpty)
}

func (s *ValidateTestSuite) TestInvalidNamespace() {
	rel := release.NewConfig()
	rel.NamespaceF = "///"

	var e *release.InvalidNamespaceError
	s.Require().ErrorAs(rel.Validate(), &e)
	s.Equal(rel.NamespaceF, e.Namespace)
}

func (s *ValidateTestSuite) TestInvalidUniq() {
	rel := release.NewConfig()
	rel.NameF = "bla@1@2@3"

	var e *uniqname.ValidationError
	s.Require().ErrorAs(rel.Validate(), &e)
	s.Equal(rel.Uniq().String(), e.Uniq)
}

func TestValidateTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ValidateTestSuite))
}
