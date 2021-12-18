package helper

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type FileTestSuite struct {
	suite.Suite
}

func (s *FileTestSuite) TestGood() {
	b := Contains("c", []string{
		"a",
		"b",
		"c",
		"d",
		"c",
	})

	s.Require().True(b)
}

func (s *FileTestSuite) TestBad() {
	b := Contains("12", []string{
		"a",
		"b",
		"c",
		"d",
		"c",
	})

	s.Require().False(b)
}

func TestFileTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(FileTestSuite))
}
