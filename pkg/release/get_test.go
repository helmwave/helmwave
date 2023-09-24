//go:build integration

package release_test

import (
	"context"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/helmwave/go-fsimpl"
	"github.com/helmwave/go-fsimpl/filefs"
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

type GetTestSuite struct {
	suite.Suite
}

func (s *GetTestSuite) SetupSuite() {
	var rs repo.Configs
	str := `
- name: bitnami
  url: https://charts.bitnami.com/bitnami
`
	err := yaml.Unmarshal([]byte(str), &rs)

	s.Require().NoError(err)
	s.Require().Len(rs, 1)

	s.Require().NoError(plan.SyncRepositories(context.Background(), rs))
}

func (s *GetTestSuite) TestGetNotInstalled() {
	rel := release.NewConfig()
	rel.NamespaceF = strings.ToLower(strings.ReplaceAll(s.T().Name(), "/", ""))
	rel.CreateNamespace = true
	rel.Wait = false
	rel.ChartF.Name = "bitnami/nginx"

	r, err := rel.Get(0)
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

	wd, _ := os.Getwd()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: wd})
	r1, err := rel.Sync(context.Background(), baseFS.(fsimpl.CurrentPathFS))
	s.Require().NoError(err)
	s.Require().NotNil(r1)

	r2, err := rel.Get(0)
	s.Require().NoError(err)
	s.Require().NotNil(r2)

	_, err = rel.GetValues()
	s.Require().NoError(err)
}

func TestGetTestSuite(t *testing.T) { //nolint:paralleltest // uses helm repository.yaml flock
	// t.Parallel()
	suite.Run(t, new(GetTestSuite))
}
