package release

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
	"helm.sh/helm/v3/pkg/action"
)

type ChartInternalTestSuite struct {
	suite.Suite
}

func (s *ChartInternalTestSuite) contains(a []string, b string) bool {
	s.T().Helper()

	for i := range a {
		if a[i] == b {
			return true
		}
	}

	return false
}

// TestChartTypeFields checks that all fields of helm upgrade action exist in config structure.
func (s *ChartInternalTestSuite) TestChartTypeFields() {
	skipFields := []string{
		"Name",
		"name",
	}

	a := Chart{}
	aa := reflect.ValueOf(a).Elem().Type()
	fieldsR := make([]string, aa.NumField())

	b := action.ChartPathOptions{}
	bb := reflect.ValueOf(b).Elem().Type()

	for i := range fieldsR {
		f := aa.Field(i)
		fieldsR[i] = f.Name
	}

	for i := bb.NumField() - 1; i >= 0; i-- {
		f := bb.Field(i)
		if !s.contains(skipFields, f.Name) {
			s.Require().Contains(fieldsR, f.Name)
		}
	}
}

func TestChartInternalTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ChartInternalTestSuite))
}
