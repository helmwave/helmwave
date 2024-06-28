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
		"my@test@context",
		"my-release@test-1@context-1",
		"my-release@test-1@gke_project_asia-southeast1_cluster-name",
	}

	for _, d := range data {
		_, err := uniqname.NewFromString(d)
		s.Require().NoError(err)
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
		"my-release@test-1@Context-1@-blabla",
	}

	for _, d := range data {
		_, err := uniqname.NewFromString(d)
		s.Require().Error(err)
	}
}

func TestValidateTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ValidateTestSuite))
}
