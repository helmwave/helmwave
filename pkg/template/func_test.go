// +build ignore unit

package template

import (
	"testing"

	"github.com/Masterminds/sprig/v3"
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

func TestFuncTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(FuncTestSuite))
}
