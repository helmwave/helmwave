package registry_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/registry"
	"github.com/stretchr/testify/suite"
)

type ValidateTestSuite struct {
	suite.Suite
}

func TestValidateTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ValidateTestSuite))
}

func (ts *ValidateTestSuite) TestValidate() {
	cfg := registry.NewConfig()

	err := cfg.Validate()
	ts.Require().ErrorIs(err, registry.ErrNameEmpty)

	cfg.HostF = "123"
	err = cfg.Validate()
	ts.Require().NoError(err)
}
