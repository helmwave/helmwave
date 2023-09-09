package release_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/suite"
)

type MonitorsTestSuite struct {
	suite.Suite
}

func TestMonitorsTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(MonitorsTestSuite))
}

func (s *MonitorsTestSuite) TestActionJSONSchema() {
	schema := jsonschema.Reflect(release.MonitorActionNone)

	s.Require().NotNil(schema)
	s.Require().Equal("string", schema.Type)
}
