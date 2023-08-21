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

	s.Require().ErrorIs(release.ErrNameEmpty, rel.Validate())
}

func (s *ValidateTestSuite) TestInvalidNamespace() {
	rel := release.NewConfig()
	rel.NamespaceF = "///"

	s.Require().ErrorIs(release.InvalidNamespaceError{}, rel.Validate())
}

func (s *ValidateTestSuite) TestInvalidUniq() {
	rel := release.NewConfig()
	rel.NameF = "bla@bla"

	s.Require().ErrorIs(uniqname.ValidationError{}, rel.Validate())
}

func TestValidateTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ValidateTestSuite))
}
