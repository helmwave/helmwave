package action

import (
	"fmt"
	"slices"
	"strings"

	"github.com/helmwave/helmwave/pkg/cache"
	logSetup "github.com/helmwave/helmwave/pkg/log"
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/urfave/cli/v2"
)

const ROOT_PREFIX = "HELMWAVE_"

// EnvVars helper function for HELMWAVE environment variables.
func EnvVars(names ...string) []string {
	a := make([]string, 0, len(names))
	for _, name := range names {
		s := strings.ToUpper(ROOT_PREFIX + name)
		a = append(a, s)
	}

	return a
}

// GlobalFlags is a set of global flags.
func GlobalFlags() []cli.Flag {
	r := []cli.Flag{
		flagCancel(),
		&cli.IntFlag{
			Name:    "parallel-limit",
			Usage:   "limit amount of parallel releases",
			EnvVars: EnvVars("PARALLEL_LIMIT"),
			Value:   0,
		},
	}

	return slices.Concat(r, cache.Default.Flags(), logSetup.Default.Flags())
}

// flagCancel is flag for canceling process on SigINT or SigTERM.
func flagCancel() cli.Flag {
	return &cli.BoolFlag{
		Name:    "handle-signal",
		Usage:   "cancel helm on SigINT,SigTERM",
		Value:   false,
		EnvVars: EnvVars("HANDLE_SIGNAL"),
	}
}

// flagPlandir pass val to urfave flag.
func flagPlandir(v *string) cli.Flag {
	return &cli.PathFlag{
		Name:        "plandir",
		Aliases:     []string{"p"},
		Value:       plan.Dir,
		Category:    Step1,
		Usage:       "path to plandir",
		EnvVars:     EnvVars("PLANDIR", "PLAN"),
		Destination: v,
	}
}

// flagYmlFile pass val to urfave flag.
func flagYmlFile(v *string) cli.Flag {
	return &cli.PathFlag{
		Name:        "file",
		Category:    CategoryYML,
		Aliases:     []string{"f"},
		Value:       plan.Body,
		Usage:       "main yml file",
		EnvVars:     EnvVars("YAML", "YML"),
		Destination: v,
	}
}

// flagTplFile pass val to urfave flag.
func flagTplFile(v *string) cli.Flag {
	return &cli.PathFlag{
		Name:        "tpl",
		Category:    CategoryYML,
		Value:       "helmwave.yml.tpl",
		Usage:       "main tpl file",
		EnvVars:     EnvVars("TPL"),
		Destination: v,
	}
}

// flagTemplateEngine pass val to urfave flag.
func flagTemplateEngine() cli.Flag {
	return &cli.StringFlag{
		Name:     "templater",
		Category: CategoryYML,
		Value:    template.TemplaterSprig,
		Usage:    fmt.Sprintf("select template engine: [ %s | %s ]", template.TemplaterSprig, template.TemplaterGomplate),
		EnvVars:  EnvVars("TEMPLATER", "TEMPLATE_ENGINE"),
	}
}

// flagYmlTemplateEngine pass val to urfave flag for yml templating.
func flagYmlTemplateEngine(v *string) cli.Flag {
	return &cli.StringFlag{
		Name:     "yml-templater",
		Category: CategoryYML,
		Usage: fmt.Sprintf(
			"select template engine for rendering helmwave.yml: [ %s | %s ]",
			template.TemplaterSprig, template.TemplaterGomplate,
		),
		EnvVars:     EnvVars("YML_TEMPLATER", "YML_TEMPLATE_ENGINE"),
		Destination: v,
	}
}

// flagBuildTemplateEngine pass val to urfave flag for values templating.
func flagBuildTemplateEngine(v *string) cli.Flag {
	return &cli.StringFlag{
		Name:     "build-templater",
		Category: CategoryYML,
		Usage: fmt.Sprintf(
			"select template engine for rendering values: [ %s | %s ]",
			template.TemplaterSprig, template.TemplaterGomplate,
		),
		EnvVars:     EnvVars("BUILD_TEMPLATER", "BUILD_TEMPLATE_ENGINE"),
		Destination: v,
	}
}

// flagAutoBuild pass val to urfave flag.
func flagAutoBuild(v *bool) cli.Flag {
	return &cli.BoolFlag{
		Name:        "build",
		Usage:       "auto build",
		Value:       false,
		Category:    Step1,
		EnvVars:     EnvVars("AUTO_BUILD"),
		Destination: v,
	}
}

// flagSkipUnchanged skip unchanged releases.
func flagSkipUnchanged(v *bool) cli.Flag {
	return &cli.BoolFlag{
		Name:        "skip-unchanged",
		Usage:       "skip unchanged releases",
		Value:       false,
		Category:    Step1,
		EnvVars:     EnvVars("SKIP_UNCHANGED"),
		Destination: v,
	}
}

// flagGraphWidth pass val to an urfave flag.
func flagGraphWidth(v *int) cli.Flag {
	return &cli.IntFlag{
		Name: "graph-width",
		Usage: "set ceil width: " +
			"1 – disable graph; " +
			"0 – full names; " +
			"N>1 – show only N symbols; " +
			"N<0 – drop N symbols from end.",
		Value:       0,
		Category:    Step1,
		EnvVars:     EnvVars("GRAPH_WIDTH"),
		Destination: v,
	}
}
