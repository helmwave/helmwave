package action_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ListTestSuite struct {
	suite.Suite
}

func TestListTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ListTestSuite))
}
