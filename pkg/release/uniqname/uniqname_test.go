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
		"my@test@Context",
		"my-release@test-1@Context-1",
	}

	for _, d := range data {
		_, err := uniqname.NewFromString(d)
		s.NoError(err)
	}
}

func (s *ValidateTestSuite) TestBad() {
	data := []string{
		"my-release",
		"my",
		"my@-",
		"my@ ",
		"@Name",
		"",
		"@-",
		"-@-",
		"my-release@test-1@Context-1@blabla",
	}

	for _, d := range data {
		_, err := uniqname.NewFromString(d)
		var e *uniqname.ValidationError
		s.ErrorAs(err, &e)
		s.Equal(d, e.Uniq)
	}
}

func TestValidateTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ValidateTestSuite))
}
