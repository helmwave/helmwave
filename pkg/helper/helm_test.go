package helper

import (
	"testing"
)

func TestHelmNS(t *testing.T) {
	h1, err := NewHelm("my")
	if err != nil {
		t.Error(err)
	}

	if h1.Namespace() != "my" {
		t.Error("helm custom namespace is failed")
	}
}
