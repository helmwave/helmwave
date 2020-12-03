package helmwave

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/wayt/parallel"
	"github.com/zhilyaev/helmwave/pkg/release"
	helm "helm.sh/helm/v3/pkg/cli"
	"os"
)

func (c *Config) SyncPlan() error {
	err := c.SyncPlanRepos()
	if err != nil {
		return err
	}

	return c.SyncPlanReleases()
}

func (c *Config) SyncPlanRepos() error {
	log.Info("🗄 Sync repositories")
	for _, r := range c.Plan.Body.Repositories {
		err := r.Sync(c.Helm)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) SyncPlanReleases() error {
	log.Info("🛥 Sync releases")
	var fails []*release.Config

	if c.Parallel {
		g := &parallel.Group{}
		log.Debug("🐞 Run in parallel mode")
		for i, _ := range c.Plan.Body.Releases {
			g.Go(c.DoRelease, &c.Plan.Body.Releases[i], &fails)
		}
		err := g.Wait()
		if err != nil {
			return err
		}
	} else {
		for _, r := range c.Plan.Body.Releases {
			c.DoRelease(&r, &fails)
		}
	}

	n := len(c.Plan.Body.Releases)
	k := len(fails)

	log.Infof("Success %d / %d", n-k, n)
	if k > 0 {
		for _, rel := range fails {
			log.Errorf("%q was not deployed to %q", rel.Name, rel.Options.Namespace)
		}

		return errors.New("deploy failed")
	}
	return nil
}

func (c *Config) DoRelease(r *release.Config, fails *[]*release.Config) {
	log.Infof("🛥 %s -> %s\n", r.Name, r.Options.Namespace)

	// I hate Private
	_ = os.Setenv("HELM_NAMESPACE", r.Options.Namespace)
	settings := helm.New()
	cfg, err := c.ActionCfg(r.Options.Namespace, settings)
	if err != nil {
		log.Fatal("❌ ", err)
	}

	err = r.Sync(cfg, settings)
	if err != nil {
		log.Error("❌ ", err)
		*fails = append(*fails, r)

	} else {
		log.Infof("✅ %s -> %s\n", r.Name, r.Options.Namespace)
	}
}
