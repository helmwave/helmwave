// +build ignore unit

package plan

import (
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/repo"
	"testing"
)

func TestExportValuesEmpty(t *testing.T) {
	p := New(Dir)

	p.body = &planBody{
		Releases:     []*release.Config{},
		Repositories: []*repo.Config{},
	}

	err := p.exportValues()
	if err != nil {
		t.Error(err)
	}
}

func TestExportValuesOneRelease(t *testing.T) {
	p := New(Dir)

	p.body = &planBody{
		Releases: []*release.Config{
			{
				Name:   "bitnami",
				Values: []release.ValuesReference{},
			},
			{
				Name:   "redis",
				Values: []release.ValuesReference{},
			},
		},
	}

	err := p.exportValues()
	if err != nil {
		t.Error(err)
	}

}
