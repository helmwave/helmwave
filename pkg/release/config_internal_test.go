package release

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigInternalTestSuite struct {
	suite.Suite
}

func (s *ConfigInternalTestSuite) contains(a []string, b string) bool {
	s.T().Helper()

	for i := range a {
		if a[i] == b {
			return true
		}
	}

	return false
}

// TestConfigHelmTypeFields checks that all fields of helm upgrade action exist in config structure.
func (s *ConfigInternalTestSuite) TestConfigHelmTypeFields() {
	skipFields := []string{
		"ChartPathOptions",
		"Install",
		"Namespace",
		"DryRun",
		"DryRunOption",
		"Description",
		"PostRenderer",
		"DependencyUpdate",
		"Lock",
		"Devel", // we removed that to force everyone specify the version
	}

	r := NewConfig()
	rr := reflect.ValueOf(r).Elem().Type()
	fieldsR := make([]string, rr.NumField())

	c := r.newUpgrade()
	rc := reflect.ValueOf(c).Elem().Type()

	for i := range fieldsR {
		f := rr.Field(i)
		fieldsR[i] = f.Name
	}

	for i := rc.NumField() - 1; i >= 0; i-- {
		f := rc.Field(i)
		if !f.IsExported() {
			continue
		}
		if !s.contains(skipFields, f.Name) {
			s.Require().Contains(fieldsR, f.Name)
		}
	}
}

func TestConfigInternalTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ConfigInternalTestSuite))
}
