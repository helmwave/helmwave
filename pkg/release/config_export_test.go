// This package exports some fields for tests
package release

import (
	"math/rand"
	"strconv"

	"github.com/helmwave/helmwave/pkg/pubsub"
	"github.com/helmwave/helmwave/pkg/release/uniqname"
)

func NewConfig() *config { //nolint:revive
	return &config{
		NameF:      "blabla" + strconv.Itoa(rand.Int()),
		NamespaceF: "blabla",
	}
}

func (rel *config) GetDependencies() map[uniqname.UniqName]<-chan pubsub.ReleaseStatus {
	return rel.dependencies
}

func (rel *config) WaitForDependencies() error {
	return rel.waitForDependencies()
}
