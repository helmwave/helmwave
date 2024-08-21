package action

import (
	"fmt"
	"github.com/helmwave/helmwave/pkg/cache"
	logSetup "github.com/helmwave/helmwave/pkg/log"
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/urfave/cli/v2"
	"strings"
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

// GlobalFlags is a set of global flags
func GlobalFlags() (r []cli.Flag) {
	r = []cli.Flag{
		flagCancel(),
	}

	r = append(r, cache.Default.Flags()...)
	r = append(r, logSetup.Default.Flags()...)

	return r
}

// flagCancel is flag for canceling process on SigINT or SigTERM
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
		Category:    "BUILD",
		Usage:       "path to plandir",
		EnvVars:     EnvVars("PLANDIR", "PLAN"),
		Destination: v,
	}
}

// flagYmlFile pass val to urfave flag.
func flagYmlFile(v *string) cli.Flag {
	return &cli.PathFlag{
		Name:        "file",
		Category:    "YML",
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
		Category:    "YML",
		Value:       "helmwave.yml.tpl",
		Usage:       "main tpl file",
		EnvVars:     EnvVars("TPL"),
		Destination: v,
	}
}

// flagTemplateEngine pass val to urfave flag.
func flagTemplateEngine(v *string) cli.Flag {
	return &cli.StringFlag{
		Name:        "templater",
		Category:    "YML",
		Value:       template.TemplaterSprig,
		Usage:       fmt.Sprintf("select template engine: [ %s | %s ]", template.TemplaterSprig, template.TemplaterGomplate),
		EnvVars:     EnvVars("TEMPLATER", "TEMPLATE_ENGINE"),
		Destination: v,
	}
}

// flagAutoBuild pass val to urfave flag.
func flagAutoBuild(v *bool) cli.Flag {
	return &cli.BoolFlag{
		Name:        "build",
		Usage:       "auto build",
		Value:       false,
		Category:    "BUILD",
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
		Category:    "BUILD",
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
		Category:    "BUILD",
		EnvVars:     EnvVars("GRAPH_WIDTH"),
		Destination: v,
	}
}
