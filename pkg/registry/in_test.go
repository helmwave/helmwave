package registry_test

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/registry"
	"github.com/stretchr/testify/suite"
)

type InTestSuite struct {
	suite.Suite
}

func TestInTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(InTestSuite))
}

func (ts *InTestSuite) TestIndexOfHost() {
	reg := registry.NewConfig()
	reg.HostF = ts.T().Name()

	idx, found := registry.IndexOfHost([]registry.Config{reg, reg, reg}, ts.T().Name())

	ts.Require().True(found)
	ts.Require().Equal(0, idx)
}

func (ts *InTestSuite) TestIndexOfHostNotFound() {
	reg := registry.NewConfig()
	reg.HostF = ts.T().Name()

	_, found := registry.IndexOfHost([]registry.Config{reg}, "")

	ts.Require().False(found)
}
