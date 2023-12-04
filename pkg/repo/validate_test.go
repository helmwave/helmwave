package repo_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/repo"

	"github.com/stretchr/testify/suite"
)

type ValidateTestSuite struct {
	suite.Suite
}

func TestValidateTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ValidateTestSuite))
}

func (s *ValidateTestSuite) TestEmptyName() {
	rep := repo.NewConfig()
	rep.Entry.Name = ""

	s.Require().ErrorIs(rep.Validate(), repo.ErrNameEmpty)
}

func (s *ValidateTestSuite) TestEmptyURL() {
	rep := repo.NewConfig()
	rep.Entry.URL = ""

	s.Require().ErrorIs(rep.Validate(), repo.ErrURLEmpty)
}

func (s *ValidateTestSuite) TestInvalidURL() {
	rep := repo.NewConfig()
	rep.Entry.URL = "\\asdasd://null"

	var e *repo.InvalidURLError
	s.Require().ErrorAs(rep.Validate(), &e)
	s.Equal(rep.Entry.URL, e.URL)
}

func (s *ValidateTestSuite) TestValid() {
	rep := repo.NewConfig()

	s.Require().NoError(rep.Validate())
}
