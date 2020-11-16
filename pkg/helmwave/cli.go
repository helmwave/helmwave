package helmwave

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/wayt/parallel"
	"github.com/zhilyaev/helmwave/pkg/release"
	"github.com/zhilyaev/helmwave/pkg/template"
	"github.com/zhilyaev/helmwave/pkg/yml"
	helm "helm.sh/helm/v3/pkg/cli"
	"os"
)

func (c *Config) Render(ctx *cli.Context) error {
	if c.Debug == false {
		fmt.Println("📄 Render", c.Tpl.File, "->", c.Yml.File)
	}
	return template.Tpl2yml(c.Tpl.File, c.Yml.File, nil, c.Debug)
}

func (c *Config) Planfile(ctx *cli.Context) error {
	err := c.Render(ctx)
	if err != nil {
		return err
	}

	fmt.Println("🛠 Your planfile is", c.Plan.File)
	c.ReadHelmWaveYml()
	c.Plan.Body.Project = c.Yml.Body.Project
	c.Plan.Body.Version = c.Yml.Body.Version
	c.PlanReleases()
	c.RenderValues()

	fmt.Print("🛠 -> 🛥 [ ")
	for _, rel := range c.Plan.Body.Releases {
		fmt.Print(rel.Name, " ")
	}
	fmt.Println("]")
	c.PlanRepos()

	fmt.Print("🛠 -> 🗄  [ ")
	for _, rep := range c.Plan.Body.Repositories {
		fmt.Print(rep.Name, " ")
	}
	fmt.Println("]")
	return yml.Save(c.Plan.File, c.Plan.Body)
}

func (c *Config) SyncRepos(ctx *cli.Context) error {
	err := c.Planfile(ctx)
	if err != nil {
		return err
	}

	fmt.Println("🗄 Sync repositories")
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

	fmt.Println("🛥 Sync releases")

	if c.Parallel {
		g := &parallel.Group{}
		for i, _ := range c.Plan.Body.Releases {
			g.Go(c.DoRelease, &c.Plan.Body.Releases[i])
		}
		return g.Wait()
	} else {
		for _, r := range c.Plan.Body.Releases {
			c.DoRelease(&r)
		}
	}

	return nil
}

func (c *Config) DoRelease(r *release.Config) {
	fmt.Printf("🛥 %s -> %s\n", r.Name, r.Options.Namespace)

	// I hate Private
	_ = os.Setenv("HELM_NAMESPACE", r.Options.Namespace)
	settings := helm.New()
	cfg, err := c.ActionCfg(r.Options.Namespace, settings)
	if err != nil {
		fmt.Println("❌", err)
	}

	err = r.Sync(cfg, settings)
	if err != nil {
		fmt.Println("❌", err)
	} else {
		fmt.Printf("✅ %s -> %s\n", r.Name, r.Options.Namespace)
	}
}
