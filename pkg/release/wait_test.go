package release_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v4/pkg/kube"
)

type WaitTestSuite struct {
	suite.Suite
}

func TestWaitTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(WaitTestSuite))
}

func (s *WaitTestSuite) TestUnmarshalStrategy() {
	for _, in := range []release.WaitStrategy{
		release.WaitStrategyWatcher,
		release.WaitStrategyLegacy,
		release.WaitStrategyHookOnly,
	} {
		var got release.WaitStrategy

		s.Require().NoError(yaml.Unmarshal([]byte(in), &got), in)
		s.Equal(in, got)
		s.Require().NoError(got.Validate())
	}
}

// Booleans parse as strings and must fail validation with a hint pointing at the strategy names.
func (s *WaitTestSuite) TestBooleansRejected() {
	for _, in := range []string{"true", "false"} {
		var got release.WaitStrategy

		s.Require().NoError(yaml.Unmarshal([]byte(in), &got), in)

		err := got.Validate()

		var e *release.InvalidWaitStrategyError
		s.Require().ErrorAs(err, &e, in)
	}
}

func (s *WaitTestSuite) TestValidateUnknown() {
	err := release.WaitStrategy("yes-please").Validate()

	var e *release.InvalidWaitStrategyError
	s.Require().ErrorAs(err, &e)
}

// helm refuses to build a waiter for an empty strategy, so an unset one must not stay empty.
func (s *WaitTestSuite) TestUnsetDefaultsLikeHelm() {
	var unset release.WaitStrategy

	s.Equal(kube.HookOnlyStrategy, unset.Helm())
	s.False(unset.Enabled())
	s.Require().NoError(unset.Validate())
}

func (s *WaitTestSuite) TestEnabled() {
	s.True(release.WaitStrategyWatcher.Enabled())
	s.True(release.WaitStrategyLegacy.Enabled())
	s.False(release.WaitStrategyHookOnly.Enabled())
}

func (s *WaitTestSuite) TestJSONSchema() {
	schema := release.WaitStrategy("").JSONSchema()

	s.Require().NotNil(schema)
	s.Require().Len(schema.Enum, 4)
}
