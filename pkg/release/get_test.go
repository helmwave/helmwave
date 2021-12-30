package release_test

import (
	"strings"
	"testing"

	"github.com/helmwave/helmwave/pkg/release"
	"github.com/stretchr/testify/suite"
)

type GetTestSuite struct {
	suite.Suite
}

func (s *GetTestSuite) TestGetNotInstalled() {
	rel := release.NewConfig()
	rel.NamespaceF = strings.ToLower(strings.ReplaceAll(s.T().Name(), "/", ""))
	rel.CreateNamespace = true
	rel.Wait = false
	rel.ChartF.Name = "bitnami/nginx"

	r, err := rel.Get()
	s.Require().Error(err)
	s.Require().Nil(r)

	_, err = rel.GetValues()
	s.Require().Error(err)
}

func (s *GetTestSuite) TestGet() {
	rel := release.NewConfig()
	rel.NamespaceF = strings.ToLower(strings.ReplaceAll(s.T().Name(), "/", ""))
	rel.CreateNamespace = true
	rel.Wait = false
	rel.ChartF.Name = "bitnami/nginx"

	r1, err := rel.Sync()
	s.Require().NoError(err)
	s.Require().NotNil(r1)

	r2, err := rel.Get()
	s.Require().NoError(err)
	s.Require().NotNil(r2)

	_, err = rel.GetValues()
	s.Require().NoError(err)
}

func TestGetTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(GetTestSuite))
}
