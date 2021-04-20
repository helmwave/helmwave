// +build ignore integration

package integration

import (
	"github.com/helmwave/helmwave/pkg/helmwave"
	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/kubedog"
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/helmwave/helmwave/pkg/yml"
	"github.com/stretchr/testify/suite"
	"helm.sh/helm/v3/pkg/action"
	helm "helm.sh/helm/v3/pkg/cli"
	helmRelease "helm.sh/helm/v3/pkg/release"
	"os"
	"testing"
)

type SyncTestSuite struct {
	suite.Suite
	app *helmwave.Config
}

func (s *SyncTestSuite) SetupTest() {
	s.app = &helmwave.Config{
		Helm: helm.New(),
		Tpl: template.Tpl{
			From: "../fixtures/helmwave.yml.tpl",
			To:   "helmwave.yml",
		},
		Yml:      yml.Config{},
		PlanPath: PlanPath,
		Logger: &helmwave.Log{
			Level:  "DEBUG",
			Format: "flat",
			Color:  false,
		},
		Kubedog: &kubedog.Config{},
	}

	err := s.app.InitLogger()
	s.Require().NoError(err)

	err = s.app.Tpl.Render()
	s.Require().NoError(err)

	err = yml.Read(s.app.Tpl.To, &s.app.Yml)
	s.Require().NoError(err)
}

func (s *SyncTestSuite) TestSync() {
	opts := &yml.SavePlanOptions{}
	opts.File(s.app.PlanPath + PlanFile).Dir(s.app.PlanPath)
	opts.PlanValues().PlanRepos().PlanValues()
	opts.Tags(s.app.Tags.Value())

	err := s.app.Yml.Plan(opts, s.app.Helm)
	s.Require().NoError(err)

	_ = os.Setenv("HELM_NAMESPACE", "test-nginx")

	err = s.app.Yml.Sync(s.app.PlanPath+PlanFile, s.app.Helm)
	s.Require().NoError(err)

	cfg, err := helper.ActionCfg("test-nginx", s.app.Helm)
	s.Require().NoError(err)
	s.Require().NotNil(cfg)

	list := action.NewList(cfg)
	list.All = true
	s.Require().NotNil(list)

	releases, err := list.Run()
	s.Require().NoError(err)
	s.Require().NotZero(len(releases))

	found := false
	for _, release := range releases {
		s.Require().NotNil(release)
		if release.Name != "nginx" {
			continue
		}

		found = true

		s.Require().NotNil(release.Info)
		s.Require().Equal(helmRelease.StatusDeployed, release.Info.Status)
	}

	s.Require().True(found, "Release not found")
}

func TestSyncTestSuite(t *testing.T) {
	suite.Run(t, new(SyncTestSuite))
}
