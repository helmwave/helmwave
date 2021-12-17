//go:build ignore || unit

package plan

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type ExportTestSuite struct {
	suite.Suite
}

func (s *ExportTestSuite) TestValuesEmpty() {
	p := New(Dir)

	p.body = &planBody{
		Releases:     releaseConfigs{},
		Repositories: repoConfigs{},
	}

	err := p.exportValues()
	s.Require().NoError(err)
}

/*
// TODO: fix release.Config usage
func (s *ExportTestSuite) TestValuesOneRelease() {
	p := New(Dir)

	p.body = &planBody{
		Releases: releaseConfigs{
			&release.Config{
				Name:   "bitnami",
				Values: []release.ValuesReference{},
			},
			&release.Config{
				Name:   "redis",
				Values: []release.ValuesReference{},
			},
		},
	}

	err := p.exportValues()
	s.Require().NoError(err)
}
*/

func TestExportTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(ExportTestSuite))
}
