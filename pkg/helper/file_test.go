package helper_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/stretchr/testify/suite"
)

type FileTestSuite struct {
	suite.Suite
}

func (s *FileTestSuite) TestGood() {
	b := helper.Contains("c", []string{
		"a",
		"b",
		"c",
		"d",
		"c",
	})

	s.Require().True(b)
}

func (s *FileTestSuite) TestBad() {
	b := helper.Contains("12", []string{
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
