//go:build ignore || unit

package helper

import "testing"

func TestContains(t *testing.T) {
	b := Contains("c", []string{
		"a",
		"b",
		"c",
		"d",
		"c",
	})

	if !b {
		t.Error("bad contains: True-False")
	}

	b = Contains("12", []string{
		"a",
		"b",
		"c",
		"d",
		"c",
	})

	if b {
		t.Error("bad contains: False-True")
	}
}
