//go:build integration

package release_test

import (
	"context"
	"strings"
	"testing"

	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/release"
	"github.com/helmwave/helmwave/pkg/repo"
	"github.com/helmwave/helmwave/tests"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

type GetTestSuite struct {
	suite.Suite

	ctx context.Context
}

func TestGetTestSuite(t *testing.T) { //nolint:paralleltest // uses helm repository.yaml flock
	// t.Parallel()
	suite.Run(t, new(GetTestSuite))
}

func (ts *GetTestSuite) SetupSuite() {
	ts.ctx = tests.GetContext(ts.T())

	var rs repo.Configs
	str := `
- name: bitnami
  url: https://charts.bitnami.com/bitnami
`
	err := yaml.Unmarshal([]byte(str), &rs)

	ts.Require().NoError(err)
	ts.Require().Len(rs, 1)

	ts.Require().NoError(plan.SyncRepositories(ts.ctx, rs))
}

func (ts *GetTestSuite) TestGetNotInstalled() {
	rel := release.NewConfig()
	rel.NamespaceF = strings.ToLower(strings.ReplaceAll(ts.T().Name(), "/", ""))
	rel.CreateNamespace = true
	rel.Wait = false
	rel.ChartF.Name = "bitnami/nginx"

	r, err := rel.Get(0)
	ts.Require().Error(err)
	ts.Require().Nil(r)

	_, err = rel.GetValues()
	ts.Require().Error(err)
}

func (ts *GetTestSuite) TestGet() {
	rel := release.NewConfig()
	rel.NamespaceF = strings.ToLower(strings.ReplaceAll(ts.T().Name(), "/", ""))
	rel.CreateNamespace = true
	rel.Wait = false
	rel.ChartF.Name = "bitnami/nginx"

	r1, err := rel.Sync(ts.ctx, false)
	ts.Require().NoError(err)
	ts.Require().NotNil(r1)

	r2, err := rel.Get(0)
	ts.Require().NoError(err)
	ts.Require().NotNil(r2)

	_, err = rel.GetValues()
	ts.Require().NoError(err)
}
