package release_test

import (
	"strings"
	"testing"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
	helmRepo "helm.sh/helm/v3/pkg/repo"
)

type GetTestSuite struct {
	suite.Suite
}

func (s *GetTestSuite) SetupSuite() {
	var rs rt
	str := `
- name: bitnami
  url: https://charts.bitnami.com/bitnami
`
	err := yaml.Unmarshal([]byte(str), &rs)

	s.Require().NoError(err)
	s.Require().Len(rs, 1)

	r := rs[0]

	var f *helmRepo.File
	// Create if not exits
	if !helper.IsExists(helper.Helm.RepositoryConfig) {
		f = helmRepo.NewFile()

		_, err = helper.CreateFile(helper.Helm.RepositoryConfig)
		s.Require().NoError(err)
	} else {
		f, err = helmRepo.LoadFile(helper.Helm.RepositoryConfig)
		s.Require().NoError(err)
	}

	err = r.Install(helper.Helm, f)
	s.Require().NoError(err)
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
