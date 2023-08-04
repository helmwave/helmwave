//go:build ignore || integration

package release_test

import (
	"context"
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
	"path/filepath"
	"strings"
	"testing"
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
		Dst: filepath.Join(tests.Root, "06_values.yaml"),
	})

	r, err := rel.Sync(context.Background())
	s.Require().NoError(err)
	s.Require().NotNil(r)

	r, err = rel.Sync(context.Background())
	s.Require().NoError(err)
	s.Require().NotNil(r)
}

func (s *SyncTestSuite) TestSyncWithoutCRD() {
	rel := release.NewConfig()
	rel.NamespaceF = strings.ToLower(strings.ReplaceAll(s.T().Name(), "/", ""))
	rel.CreateNamespace = true
	rel.Wait = false
	rel.ChartF.Name = "prometheus-community/kube-prometheus-stack"

	rel.DryRun(true)

	r, err := rel.Sync(context.Background())
	s.Require().NoError(err)
	s.Require().NotNil(r)
}

func TestSyncTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(SyncTestSuite))
}
