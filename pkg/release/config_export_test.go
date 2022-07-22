// This package exports some fields for tests
package release

import (
	"math/rand"
	"strconv"
)

func NewConfig() *config { //nolint:revive
	return &config{
		NameF:      "blabla" + strconv.Itoa(rand.Int()),
		NamespaceF: "blabla",
	}
}
