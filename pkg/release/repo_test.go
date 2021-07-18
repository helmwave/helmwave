package release

import (
	"errors"
	"testing"
)

func TestConfig_Repo(t *testing.T) {
	const bitnami = "bitnami"
	r := &Config{Chart: Chart{
		Name: bitnami + "/redis",
	}}

	if r.Repo() != bitnami {
		t.Error(errors.New("get repo failed"))
	}
}
