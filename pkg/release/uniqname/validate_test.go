//go:build ignore || unit

package uniqname

import (
	"testing"

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
		s.Require().True(UniqName(d).Validate())
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
		s.Require().False(UniqName(d).Validate())
	}
}

func TestValidateTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ValidateTestSuite))
}
