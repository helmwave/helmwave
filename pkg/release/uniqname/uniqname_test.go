package uniqname_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/release/uniqname"
	"github.com/stretchr/testify/suite"
)

type ValidateTestSuite struct {
	suite.Suite
}

func (s *ValidateTestSuite) TestGood() {
	data := []string{
		"my@test",
		"my-release@test-1",
	}

	for _, d := range data {
		s.Require().NoError(uniqname.UniqName(d).Validate())
	}
}

func (s *ValidateTestSuite) TestBad() {
	data := []string{
		"my-release",
		"my",
		"my@",
		"my@-",
		"my@ ",
		"@name",
		"@",
		"@-",
		"-@-",
	}

	for _, d := range data {
		s.Require().ErrorIs(uniqname.UniqName(d).Validate(), uniqname.ValidationError{})
	}
}

func TestValidateTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ValidateTestSuite))
}
