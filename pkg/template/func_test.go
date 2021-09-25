//go:build ignore || unit

package template

import (
	"context"
	"testing"

	"github.com/Masterminds/sprig/v3"
	"github.com/hairyhenderson/gomplate/v3/funcs"
	"github.com/stretchr/testify/suite"
)

type FuncTestSuite struct {
	suite.Suite
}

func (s *FuncTestSuite) TestFuncMap() {
	sprigFuncs := sprig.FuncMap()

	fm := FuncMap()

	for key := range sprigFuncs {
		if alias, ok := sprigAliases[key]; ok {
			key = alias
		}

		s.Contains(fm, key)
	}

	for key := range customFuncs {
		s.Contains(fm, key)
	}
}

func (s *FuncTestSuite) TestEnabledGomplate() {
	SetConfig(&Config{
		Gomplate: GomplateConfig{
			Enabled: true,
		},
	})
	fm := FuncMap()

	for key := range funcs.CreateDataFuncs(context.Background(), nil) {
		s.Contains(fm, key)
	}
}

func (s *FuncTestSuite) TestDisabledGomplate() {
	SetConfig(&Config{
		Gomplate: GomplateConfig{
			Enabled: false,
		},
	})
	fm := FuncMap()

	for key := range funcs.CreateDataFuncs(context.Background(), nil) {
		s.NotContains(fm, key)
	}
}

func TestFuncTestSuite(t *testing.T) {
	//t.Parallel()
	suite.Run(t, new(FuncTestSuite))
}
