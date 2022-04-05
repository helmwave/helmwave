package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

// flagPlandir pass val to urfave flag
func flagPlandir(v *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "plandir",
		Aliases:     []string{"p"},
		Value:       plan.Dir,
		Usage:       "Path to plandir",
		EnvVars:     []string{"HELMWAVE_PLANDIR", "HELMWAVE_PLAN"},
		Destination: v,
	}
}

// flagTags pass val to urfave flag
func flagTags(v *cli.StringSlice) *cli.StringSliceFlag {
	return &cli.StringSliceFlag{
		Name:        "tags",
		Aliases:     []string{"t"},
		Usage:       "It allows you choose releases for sync. Example: -t tag1 -t tag3,tag4",
		EnvVars:     []string{"HELMWAVE_TAGS"},
		Destination: v,
	}
}

// flagTemplateEngine pass val to urfave flag
func flagMatchAllTags(v *bool) *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:        "match-all-tags",
		Usage:       "Match all provided tags",
		Value:       false,
		EnvVars:     []string{"HELMWAVE_MATCH_ALL_TAGS"},
		Destination: v,
	}
}

// flagYmlFile pass val to urfave flag
func flagYmlFile(v *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "file",
		Aliases:     []string{"f"},
		Value:       plan.Body,
		Usage:       "Main yml file",
		EnvVars:     []string{"HELMWAVE_YAML", "HELMWAVE_YML"},
		Destination: v,
	}
}

// flagTplFile pass val to urfave flag
func flagTplFile(v *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "tpl",
		Value:       "helmwave.yml.tpl",
		Usage:       "Main tpl file",
		EnvVars:     []string{"HELMWAVE_TPL"},
		Destination: v,
	}
}

// flagDiffMode pass val to urfave flag
func flagDiffMode(v *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "diff-mode",
		Value:       "live",
		Usage:       "You can set: [ live | local ]",
		EnvVars:     []string{"HELMWAVE_DIFF_MODE"},
		Destination: v,
	}
}

// flagDiffWide pass val to urfave flag
func flagDiffWide(v *int) *cli.IntFlag {
	return &cli.IntFlag{
		Name:        "wide",
		Value:       5,
		Usage:       "Show line around change",
		EnvVars:     []string{"HELMWAVE_DIFF_WIDE"},
		Destination: v,
	}
}

// flagDiffShowSecret pass val to urfave flag
func flagDiffShowSecret(v *bool) *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:        "show-secret",
		Value:       true,
		Usage:       "Show secret in diff",
		EnvVars:     []string{"HELMWAVE_DIFF_SHOW_SECRET"},
		Destination: v,
	}
}

// flagTemplateEngine pass val to urfave flag
func flagTemplateEngine(v *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "templater",
		Value:       "sprig",
		Usage:       "Select template engine: sprig or gomplate",
		EnvVars:     []string{"HELMWAVE_TEMPLATER", "HELMWAVE_TEMPLATE_ENGINE"},
		Destination: v,
	}
}

// flagAutoBuild pass val to urfave flag
func flagAutoBuild(v *bool) *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:        "build",
		Usage:       "auto build",
		Value:       false,
		EnvVars:     []string{"HELMWAVE_AUTO_BUILD"},
		Destination: v,
	}
}
