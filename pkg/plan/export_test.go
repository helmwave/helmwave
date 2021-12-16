//go:build ignore || unit

package plan

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/stretchr/testify/suite"
)

type ExportTestSuite struct {
	suite.Suite
}

func (s *ExportTestSuite) TestValuesEmpty() {
	p := New(Dir)

	p.body = &planBody{
		Releases:     []*release.Config{},
		Repositories: []*repo.Config{},
	}

	err := p.exportValues()
	s.Require().NoError(err)
}

func (s *ExportTestSuite) TestValuesOneRelease() {
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
	s.Require().NoError(err)
}

func TestExportTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ExportTestSuite))
}
