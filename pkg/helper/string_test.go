// +build ignore unit

package helper

import (
	"github.com/helmwave/helmwave/pkg/yml"
	"testing"
)

func TestString(t *testing.T) {
	b := &yml.Config{
		Project: "my-project",
		Version: "0.7.0",
	}

	s := String(b)
	const c = "project: my-project\nversion: 0.7.0\nrepositories: []\nreleases: []\n"
	if s != c {
		t.Error("Failed yml.String")
	}
}
