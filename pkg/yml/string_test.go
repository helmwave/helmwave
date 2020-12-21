package yml

import (
	"testing"
)

func TestString(t *testing.T) {
	b := &Config{
		Project: "my-project",
		Version: "0.7.0",
	}

	s := String(b)
	const c = "project: my-project\nversion: 0.7.0\nrepositories: []\nreleases: []\n"
	if s != c {
		t.Error("Failed yml.String")
	}
}
