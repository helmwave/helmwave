// This package exports some fields for tests
package release

import (
	"math/rand"
	"strconv"
)

func NewConfig() *Release {
	return &Release{
		NameF:      "blabla" + strconv.Itoa(rand.Int()),
		NamespaceF: "blabla",
	}
}

func (rel *Release) GetDryRun() bool {
	return rel.dryRun
}

func (rel *Release) BuildAfterUnmarshal() {
	rel.buildAfterUnmarshal()
}
