package release

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigInternalTestSuite struct {
	suite.Suite
}

// TestConfigHelmTypeFields checks that all fields of helm upgrade action exist in config structure.
func (s *ConfigInternalTestSuite) TestConfigHelmTypeFields() {
	skipFields := []string{
		"ChartPathOptions",
		"Install",
		"Namespace",
		"DryRun",
		"HideSecret",
		"DryRunOption",
		"Description",
		"PostRenderer",
		"DependencyUpdate",
		"Lock",
		"Devel", // we removed that to force everyone specify the version
	}

	r := NewConfig()
	rr := reflect.ValueOf(r).Elem().Type()
	fieldsR := make([]string, 0, rr.NumField())

	c := r.newUpgrade()
	rc := reflect.ValueOf(c).Elem().Type()

	for i := range rr.NumField() {
		f := rr.Field(i)
		if !f.IsExported() {
			continue
		}
		if strings.HasSuffix(f.Name, "F") {
			continue
		}

		fieldsR = append(fieldsR, f.Name)
	}

	for i := range rc.NumField() {
		f := rc.Field(i)
		if !f.IsExported() {
			continue
		}
		if slices.Contains(skipFields, f.Name) {
			continue
		}

		s.Containsf(fieldsR, f.Name, "helm upgrade field %q is not supported by helmwave config", f.Name)
	}
}

func TestConfigInternalTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ConfigInternalTestSuite))
}
