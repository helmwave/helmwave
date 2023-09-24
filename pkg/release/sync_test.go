//go:build integration

package release_test

import (
	"context"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/helmwave/go-fsimpl/filefs"
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

type SyncTestSuite struct {
	suite.Suite
}

func (s *SyncTestSuite) SetupSuite() {
	var rs repo.Configs
	str := `
- name: prometheus-community
  url: https://prometheus-community.github.io/helm-charts
- name: bitnami
  url: https://charts.bitnami.com/bitnami
`
	err := yaml.Unmarshal([]byte(str), &rs)

	s.Require().NoError(err)

	s.Require().NoError(plan.SyncRepositories(context.Background(), rs))
}

func (s *SyncTestSuite) TestInstallUpgrade() {
	rel := release.NewConfig()
	rel.NamespaceF = strings.ToLower(strings.ReplaceAll(s.T().Name(), "/", ""))
	rel.CreateNamespace = true
	rel.Wait = false
	rel.ChartF.Name = "bitnami/nginx"
	rel.ValuesF = append(rel.ValuesF, release.ValuesReference{
		Src: filepath.Join(tests.Root, "06_values.yaml"),
	})

	wd, _ := os.Getwd()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: wd})
	r, err := rel.Sync(context.Background(), baseFS)
	s.Require().NoError(err)
	s.Require().NotNil(r)

	r, err = rel.Sync(context.Background(), baseFS)
	s.Require().NoError(err)
	s.Require().NotNil(r)
}

func (s *SyncTestSuite) TestInvalidValues() {
	rel := release.NewConfig()
	rel.NamespaceF = strings.ToLower(strings.ReplaceAll(s.T().Name(), "/", ""))
	rel.CreateNamespace = true
	rel.Wait = false
	rel.ChartF.Name = "bitnami/nginx"
	rel.ValuesF = append(rel.ValuesF, release.ValuesReference{})

	wd, _ := os.Getwd()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: wd})
	r, err := rel.Sync(context.Background(), baseFS)
	s.Require().Error(err)
	s.Require().Nil(r)
}

func (s *SyncTestSuite) TestSyncWithoutCRD() {
	rel := release.NewConfig()
	rel.NamespaceF = strings.ToLower(strings.ReplaceAll(s.T().Name(), "/", ""))
	rel.CreateNamespace = true
	rel.Wait = false
	rel.ChartF.Name = "prometheus-community/kube-prometheus-stack"

	rel.DryRun(true)

	wd, _ := os.Getwd()
	baseFS, _ := filefs.New(&url.URL{Scheme: "file", Path: wd})
	r, err := rel.Sync(context.Background(), baseFS)
	s.Require().NoError(err)
	s.Require().NotNil(r)
}

func TestSyncTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(SyncTestSuite))
}
