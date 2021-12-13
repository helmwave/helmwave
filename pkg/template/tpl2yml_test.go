//go:build ignore || unit

package template

import (
	"os"
	"path"
	"testing"

	"github.com/helmwave/helmwave/tests"

	"github.com/stretchr/testify/suite"
)

type Tpl2YmlTestSuite struct {
	suite.Suite
}

func (s *Tpl2YmlTestSuite) TestDisabledGomplate() {
	gomplateConfig := &GomplateConfig{
		Enabled: false,
	}

	tpl := path.Join(tests.Root, "09_values.yaml")

	dst, err := os.CreateTemp("", "helmwave")
	s.Require().NoError(err)
	dst.Close()
	defer os.Remove(dst.Name())

	err = Tpl2yml(tpl, dst.Name(), nil, gomplateConfig)
	s.Require().Error(err)
}

func (s *Tpl2YmlTestSuite) TestEnabledGomplate() {
	gomplateConfig := &GomplateConfig{
		Enabled: true,
	}

	tpl := path.Join(tests.Root, "09_values.yaml")

	dst, err := os.CreateTemp("", "helmwave")
	s.Require().NoError(err)
	dst.Close()
	defer os.Remove(dst.Name())

	err = Tpl2yml(tpl, dst.Name(), nil, gomplateConfig)
	s.Require().NoError(err)
}

func TestTpl2YmlTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(Tpl2YmlTestSuite))
}
