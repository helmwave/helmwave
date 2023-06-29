package action_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type RollbackTestSuite struct {
	suite.Suite
}

func TestRollbackTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(RollbackTestSuite))
}
