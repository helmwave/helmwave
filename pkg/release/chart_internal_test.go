package release

import (
	"reflect"
	"slices"
	"testing"

	"github.com/stretchr/testify/suite"
	"helm.sh/helm/v4/pkg/action"
	chart "helm.sh/helm/v4/pkg/chart/v2"
)

type ChartInternalTestSuite struct {
	suite.Suite
}

func TestChartInternalTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ChartInternalTestSuite))
}

// TestChartTypeFields checks that all fields of helm upgrade action exist in config structure.
func (ts *ChartInternalTestSuite) TestChartTypeFields() {
	skipFields := []string{
		"Name",
	}

	a := Chart{}
	aa := reflect.ValueOf(a).Type()
	fieldsR := make([]string, aa.NumField())

	b := action.ChartPathOptions{}
	bb := reflect.ValueOf(b).Type()

	for i := range fieldsR {
		f := aa.Field(i)
		fieldsR[i] = f.Name
	}

	for f := range bb.Fields() {
		if !f.IsExported() {
			continue
		}
		if !slices.Contains(skipFields, f.Name) {
			ts.Require().Contains(fieldsR, f.Name)
		}
	}
}

func (ts *ChartInternalTestSuite) TestChartCheckMissingDependency() {
	rel := NewConfig()
	err := rel.chartCheck(&chart.Chart{
		Metadata: &chart.Metadata{
			Dependencies: []*chart.Dependency{
				{
					Name: ts.T().Name(),
				},
			},
		},
	})

	ts.Require().ErrorContains(err, "found in Chart.yaml, but missing in charts/ directory")
}

func (ts *ChartInternalTestSuite) TestChartCheckInvalidType() {
	rel := NewConfig()
	err := rel.chartCheck(&chart.Chart{
		Metadata: &chart.Metadata{
			Type: "library",
		},
	})

	ts.Require().NoError(err)
}

func (ts *ChartInternalTestSuite) TestChartCheckDeprecated() {
	rel := NewConfig()
	err := rel.chartCheck(&chart.Chart{
		Metadata: &chart.Metadata{
			Deprecated: true,
		},
	})

	ts.Require().NoError(err)
}
