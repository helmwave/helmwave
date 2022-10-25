// This package exports some fields for tests
package release

import (
	"math/rand"
	"strconv"
)

func NewConfig() *config {
	return &config{
		NameF:      "blabla" + strconv.Itoa(rand.Int()),
		NamespaceF: "blabla",
	}
}

func (rel *config) IsDryRun() bool {
	return rel.dryRun
}

func (rel *config) BuildAfterUnmarshal(releases ...*config) {
	rel.buildAfterUnmarshal(releases)
}
