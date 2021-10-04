//go:build ignore || unit

package template

import (
	"github.com/helmwave/helmwave/tests"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/suite"
)

type Tpl2YmlTestSuite struct {
	suite.Suite
}

func (s *Tpl2YmlTestSuite) TestDisabledGomplate() {
	SetConfig(&Config{
		Gomplate: GomplateConfig{
			Enabled: false,
		},
	})

	tpl := path.Join(tests.Root, "09_values.yaml")

	dst, err := os.CreateTemp("", "helmwave")
	s.Require().NoError(err)
	dst.Close()
	defer os.Remove(dst.Name())

	err = Tpl2yml(tpl, dst.Name(), nil)
	s.Require().Error(err)
}

func (s *Tpl2YmlTestSuite) TestEnabledGomplate() {
	SetConfig(&Config{
		Gomplate: GomplateConfig{
			Enabled: true,
		},
	})

	tpl := path.Join(tests.Root, "09_values.yaml")

	dst, err := os.CreateTemp("", "helmwave")
	s.Require().NoError(err)
	dst.Close()
	defer os.Remove(dst.Name())

	err = Tpl2yml(tpl, dst.Name(), nil)
	s.Require().NoError(err)
}

func TestTpl2YmlTestSuite(t *testing.T) {
	//t.Parallel()
	suite.Run(t, new(Tpl2YmlTestSuite))
}
