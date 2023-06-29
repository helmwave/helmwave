package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/urfave/cli/v2"
)

// flagPlandir pass val to urfave flag.
func flagPlandir(v *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "plandir",
		Aliases:     []string{"p"},
		Value:       plan.Dir,
		Usage:       "path to plandir",
		EnvVars:     []string{"HELMWAVE_PLANDIR", "HELMWAVE_PLAN"},
		Destination: v,
	}
}

// flagTags pass val to urfave flag.
func flagTags(v *cli.StringSlice) *cli.StringSliceFlag {
	return &cli.StringSliceFlag{
		Name:        "tags",
		Aliases:     []string{"t"},
		Usage:       "build releases by tags: -t tag1 -t tag3,tag4",
		EnvVars:     []string{"HELMWAVE_TAGS"},
		Destination: v,
	}
}

// flagTemplateEngine pass val to urfave flag.
func flagMatchAllTags(v *bool) *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:        "match-all-tags",
		Usage:       "match all provided tags",
		Value:       false,
		EnvVars:     []string{"HELMWAVE_MATCH_ALL_TAGS"},
		Destination: v,
	}
}

// flagYmlFile pass val to urfave flag.
func flagYmlFile(v *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "file",
		Aliases:     []string{"f"},
		Value:       plan.Body,
		Usage:       "main yml file",
		EnvVars:     []string{"HELMWAVE_YAML", "HELMWAVE_YML"},
		Destination: v,
	}
}

// flagTplFile pass val to urfave flag.
func flagTplFile(v *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "tpl",
		Value:       "helmwave.yml.tpl",
		Usage:       "main tpl file",
		EnvVars:     []string{"HELMWAVE_TPL"},
		Destination: v,
	}
}

// flagDiffMode pass val to urfave flag.
func flagDiffMode(v *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "diff-mode",
		Value:       "live",
		Usage:       "you can set: [ live | local | none ]",
		EnvVars:     []string{"HELMWAVE_DIFF_MODE"},
		Destination: v,
	}
}

// flagDiffWide pass val to urfave flag.
func flagDiffWide(v *int) *cli.IntFlag {
	return &cli.IntFlag{
		Name:        "wide",
		Value:       5,
		Usage:       "show line around changes",
		EnvVars:     []string{"HELMWAVE_DIFF_WIDE"},
		Destination: v,
	}
}

// flagDiffShowSecret pass val to urfave flag.
func flagDiffShowSecret(v *bool) *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:        "show-secret",
		Value:       true,
		Usage:       "show secret in diff",
		EnvVars:     []string{"HELMWAVE_DIFF_SHOW_SECRET"},
		Destination: v,
	}
}

// flagTemplateEngine pass val to urfave flag.
func flagTemplateEngine(v *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "templater",
		Value:       template.TemplaterSprig,
		Usage:       "select template engine: [ sprig | gomplate ]",
		EnvVars:     []string{"HELMWAVE_TEMPLATER", "HELMWAVE_TEMPLATE_ENGINE"},
		Destination: v,
	}
}

// flagAutoBuild pass val to urfave flag.
func flagAutoBuild(v *bool) *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:        "build",
		Usage:       "auto build",
		Value:       false,
		EnvVars:     []string{"HELMWAVE_AUTO_BUILD"},
		Destination: v,
	}
}

func flagDiffThreeWayMerge(v *bool) *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:        "3-way-merge",
		Usage:       "show 3-way merge diff",
		Value:       false,
		EnvVars:     []string{"HELMWAVE_DIFF_3_WAY_MERGE"},
		Destination: v,
	}
}

// flagDiffMode pass val to urfave flag.
func flagChartsCacheDir(v *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "charts-cache-dir",
		Value:       "",
		Usage:       "enable caching of helm charts in specified directory",
		EnvVars:     []string{"HELMWAVE_CHARTS_CACHE"},
		Destination: v,
	}
}

// flagSkipUnchanged skip unchanged releases.
func flagSkipUnchanged(v *bool) *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:        "skip-unchanged",
		Usage:       "skip unchanged releases",
		Value:       false,
		EnvVars:     []string{"HELMWAVE_SKIP_UNCHANGED"},
		Destination: v,
	}
}

// flagGraphWidth pass val to an urfave flag.
func flagGraphWidth(v *int) *cli.IntFlag {
	return &cli.IntFlag{
		Name: "graph-width",
		Usage: "set ceil width: " +
			"1 – disable graph; " +
			"0 – full names; " +
			"N>1 – show only N symbols; " +
			"N<0 – drop N symbols from end.",
		Value:       0,
		EnvVars:     []string{"HELMWAVE_GRAPH_WIDTH"},
		Destination: v,
	}
}
