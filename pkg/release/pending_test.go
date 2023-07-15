package release_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/stretchr/testify/suite"
)

type PendingTestSuite struct {
	suite.Suite
}

func TestPendingTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(PendingTestSuite))
}

func (s *PendingTestSuite) TestConfigsJSONSchema() {
	schema := release.PendingStrategy("").JSONSchema()

	s.Require().NotNil(schema)
	s.Require().Equal("string", schema.Type)
}
