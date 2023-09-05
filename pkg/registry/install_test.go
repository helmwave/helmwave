package registry_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/registry"
	"github.com/stretchr/testify/suite"
)

type InstallTestSuite struct {
	suite.Suite
}

func TestInstallTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(InstallTestSuite))
}

func (ts *InstallTestSuite) TestInstallPublic() {
	reg := registry.NewConfig()

	err := reg.Install()

	ts.Require().NoError(err)
}

func (ts *InstallTestSuite) TestInstallPrivateError() {
	reg := registry.NewConfig()
	reg.Username = ts.T().Name()

	err := reg.Install()

	ts.Require().ErrorIs(err, registry.LoginError{})
}
