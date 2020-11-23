package helmwave

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/wayt/parallel"
	"github.com/zhilyaev/helmwave/pkg/release"
	"github.com/zhilyaev/helmwave/pkg/template"
	"github.com/zhilyaev/helmwave/pkg/yml"
	helm "helm.sh/helm/v3/pkg/cli"
	"os"
)

func (c *Config) Render(ctx *cli.Context) error {
	return template.Tpl2yml(c.Tpl.File, c.Yml.File, nil)
}

func (c *Config) Planfile(ctx *cli.Context) error {
	err := c.Render(ctx)
	if err != nil {
		return err
	}

	log.Info("🛠 Your planfile is ", c.Plan.File)
	c.ReadHelmWaveYml()
	c.Plan.Body.Project = c.Yml.Body.Project
	c.Plan.Body.Version = c.Yml.Body.Version

	// Releases
	c.PlanReleases()
	c.RenderValues()
	names := make([]string, 0)
	for _, v := range c.Plan.Body.Releases {
		names = append(names, v.Name)
	}
	log.Infof("🛠 -> 🛥 %+v", names)

	// Repos
	c.PlanRepos()
	names = make([]string, 0)
	for _, v := range c.Plan.Body.Repositories {
		names = append(names, v.Name)
	}
	log.Infof("🛠 -> 🗄 %+v", names)

	return yml.Save(c.Plan.File, c.Plan.Body)
}

//func (c *Config) LogHelper(format string, v ...interface{}) {
//	format = fmt.Sprintf("🐞 %s\n", format)
//	log.Output(2, fmt.Sprintf(format, v...))
//}

func (c *Config) SyncRepos(ctx *cli.Context) error {
	err := c.Planfile(ctx)
	if err != nil {
		return err
	}

	log.Info("🗄 Sync repositories")
	for _, r := range c.Plan.Body.Repositories {
		r.Sync(c.Helm)
	}
	return nil
}

func (c *Config) SyncReleases(ctx *cli.Context) error {
	err := c.SyncRepos(ctx)
	if err != nil {
		return err
	}

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
			log.Fatal(err)
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
			log.Errorf("%q was not deploy to %q", rel.Name, rel.Options.Namespace)
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
