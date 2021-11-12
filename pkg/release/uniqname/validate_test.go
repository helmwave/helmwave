//go:build ignore || unit

package uniqname

import "testing"

func TestUniqname(t *testing.T) {
	good := []string{
		"my@test",
		"my-release@test-1",
	}

	for _, s := range good {
		if !UniqName(s).Validate() {
			t.Error("false positive: " + s)
		}
	}

	bad := []string{
		"my-release",
		"my",
		"my@",
		"my@-",
		"my@ ",
		"@name",
		"@",
		"@-",
		"-@-",
	}

	for _, s := range bad {
		if UniqName(s).Validate() {
			t.Error("false negative: " + s)
		}
	}
}
