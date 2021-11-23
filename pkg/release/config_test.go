//go:build ignore || unit

package release

import "testing"

func TestConfigUniq(t *testing.T) {
	r := &Config{
		Name:      "redis",
		Namespace: "test",
	}

	if r.Uniq() != r.uniqName {
		t.Error("method uniq() doesnt work")
	}

	if !r.Uniq().Validate() {
		t.Error("problem with validate uniqname")
	}
}
