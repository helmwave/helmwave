package helmwave

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/wayt/parallel"
	"github.com/zhilyaev/helmwave/pkg/release"
)

func (c *Config) SyncPlan() error {
	err := c.SyncPlanRepos()
	if err != nil {
		return err
	}

	return c.SyncPlanReleases()
}

func (c *Config) SyncPlanRepos() error {
	log.Info("🗄 Install repositories")
	for _, r := range c.Plan.Body.Repositories {
		err := r.Install(c.Helm)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) SyncPlanReleases() error {
	log.Info("🛥 Install releases")
	var fails []*release.Config

	if c.Parallel {
		g := &parallel.Group{}
		log.Debug("🐞 Run in parallel mode")
		for i, _ := range c.Plan.Body.Releases {
			g.Go(c.SyncRelease, &c.Plan.Body.Releases[i], &fails)
		}
		err := g.Wait()
		if err != nil {
			return err
		}
	} else {
		for _, r := range c.Plan.Body.Releases {
			c.SyncRelease(&r, &fails)
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

func (c *Config) SyncRelease(r *release.Config, fails *[]*release.Config) {
	err := r.Sync(c.PlanDir + ".manifest")
	if err != nil {
		log.Error("❌ ", err)
		*fails = append(*fails, r)
	}
}
