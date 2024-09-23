//go:build integration

package release_test

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/helmwave/helmwave/pkg/fileref"
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

type SyncTestSuite struct {
	suite.Suite

	ctx context.Context
}

func TestSyncTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(SyncTestSuite))
}

func (ts *SyncTestSuite) SetupSuite() {
	ts.ctx = tests.GetContext(ts.T())

	var rs repo.Configs
	str := `
- name: prometheus-community
  url: https://prometheus-community.github.io/helm-charts
- name: bitnami
  url: https://charts.bitnami.com/bitnami
`
	err := yaml.Unmarshal([]byte(str), &rs)

	ts.Require().NoError(err)

	ts.Require().NoError(plan.SyncRepositories(ts.ctx, rs))
}

func (ts *SyncTestSuite) TestInstallUpgrade() {
	rel := release.NewConfig()
	rel.NamespaceF = strings.ToLower(strings.ReplaceAll(ts.T().Name(), "/", ""))
	rel.CreateNamespace = true
	rel.Wait = false
	rel.ChartF.Name = "bitnami/nginx"
	rel.ValuesF = append(rel.ValuesF, fileref.Config{
		Dst: filepath.Join(tests.Root, "06_values.yaml"),
	})

	r, err := rel.Sync(ts.ctx, false)
	ts.Require().NoError(err)
	ts.Require().NotNil(r)

	r, err = rel.Sync(ts.ctx, false)
	ts.Require().NoError(err)
	ts.Require().NotNil(r)
}

func (ts *SyncTestSuite) TestInvalidValues() {
	rel := release.NewConfig()
	rel.NamespaceF = strings.ToLower(strings.ReplaceAll(ts.T().Name(), "/", ""))
	rel.CreateNamespace = true
	rel.Wait = false
	rel.ChartF.Name = "bitnami/nginx"
	rel.ValuesF = append(rel.ValuesF, fileref.Config{})

	r, err := rel.Sync(ts.ctx, false)
	ts.Require().Error(err)
	ts.Require().Nil(r)
}

func (ts *SyncTestSuite) TestSyncWithoutCRD() {
	rel := release.NewConfig()
	rel.NamespaceF = strings.ToLower(strings.ReplaceAll(ts.T().Name(), "/", ""))
	rel.CreateNamespace = true
	rel.Wait = false
	rel.ChartF.Name = "prometheus-community/kube-prometheus-stack"

	rel.DryRun(true)

	r, err := rel.Sync(ts.ctx, false)
	ts.Require().NoError(err)
	ts.Require().NotNil(r)
}
