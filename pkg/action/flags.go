package action

import (
	"fmt"
	"strings"
	"time"

	"github.com/helmwave/helmwave/pkg/kubedog"
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/helmwave/helmwave/pkg/template"
	"github.com/urfave/cli/v2"
)

const ROOT_PREFIX = "HELMWAVE_"

// EnvVars helper function for HELMWAVE environment variables
func EnvVars(names ...string) []string {
	a := make([]string, 0, len(names))
	for _, name := range names {
		s := strings.ToUpper(ROOT_PREFIX + name)
		a = append(a, s)
	}

	return a
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

// flagTags pass val to urfave flag.
func flagTags(v *cli.StringSlice) cli.Flag {
	return &cli.StringSliceFlag{
		Name:        "tags",
		Aliases:     []string{"t"},
		Usage:       "build releases by tags: -t tag1 -t tag3,tag4",
		Category:    "SELECTION",
		EnvVars:     EnvVars("TAGS"),
		Destination: v,
	}
}

// flagTemplateEngine pass val to urfave flag.
func flagMatchAllTags(v *bool) cli.Flag {
	return &cli.BoolFlag{
		Name:        "match-all-tags",
		Aliases:     []string{"tt"},
		Usage:       "match all provided tags",
		Value:       false,
		Category:    "SELECTION",
		EnvVars:     EnvVars("MATCH_ALL_TAGS"),
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

// flagDiffMode pass val to urfave flag.
func flagDiffMode(v *string) cli.Flag {
	return &cli.StringFlag{
		Name:        "diff-mode",
		Value:       "live",
		Category:    "DIFF",
		Usage:       "you can set: [ live | local | none ]",
		EnvVars:     EnvVars("DIFF_MODE"),
		Destination: v,
	}
}

// flagDiffWide pass val to urfave flag.
func flagDiffWide(v *int) cli.Flag {
	return &cli.IntFlag{
		Name:        "wide",
		Value:       5,
		Category:    "DIFF",
		Usage:       "show line around changes",
		EnvVars:     EnvVars("DIFF_WIDE"),
		Destination: v,
	}
}

// flagDiffShowSecret pass val to urfave flag.
func flagDiffShowSecret(v *bool) cli.Flag {
	return &cli.BoolFlag{
		Name:        "show-secret",
		Value:       true,
		Category:    "DIFF",
		Usage:       "show secret in diff",
		EnvVars:     EnvVars("DIFF_SHOW_SECRET"),
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

func flagDiffThreeWayMerge(v *bool) cli.Flag {
	return &cli.BoolFlag{
		Name:        "3-way-merge",
		Usage:       "show 3-way merge diff",
		Value:       false,
		Category:    "DIFF",
		EnvVars:     EnvVars("DIFF_3_WAY_MERGE"),
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

func flagsKubedog(dog *kubedog.Config) []cli.Flag {
	return []cli.Flag{
		&cli.BoolFlag{
			Name:        "kubedog",
			Usage:       "enable/disable kubedog",
			Value:       false,
			Category:    "KUBEDOG",
			EnvVars:     EnvVars("KUBEDOG_ENABLED", "KUBEDOG"),
			Destination: &dog.Enabled,
		},
		&cli.DurationFlag{
			Name:        "kubedog-status-interval",
			Usage:       "interval of kubedog status messages: set -1s to stop showing status progress",
			Value:       5 * time.Second,
			Category:    "KUBEDOG",
			EnvVars:     EnvVars("KUBEDOG_STATUS_INTERVAL"),
			Destination: &dog.StatusInterval,
		},
		&cli.DurationFlag{
			Name:        "kubedog-start-delay",
			Usage:       "delay kubedog start, don't make it too late",
			Value:       time.Second,
			Category:    "KUBEDOG",
			EnvVars:     EnvVars("KUBEDOG_START_DELAY"),
			Destination: &dog.StartDelay,
		},
		&cli.DurationFlag{
			Name:        "kubedog-timeout",
			Usage:       "timeout of kubedog multitrackers",
			Value:       5 * time.Minute,
			Category:    "KUBEDOG",
			EnvVars:     EnvVars("KUBEDOG_TIMEOUT"),
			Destination: &dog.Timeout,
		},
		&cli.IntFlag{
			Name:        "kubedog-log-width",
			Usage:       "set kubedog max log line width",
			Value:       140,
			Category:    "KUBEDOG",
			EnvVars:     EnvVars("KUBEDOG_LOG_WIDTH"),
			Destination: &dog.LogWidth,
		},
		&cli.BoolFlag{
			Name:        "kubedog-track-all",
			Usage:       "track almost all resources, experimental",
			Value:       false,
			Category:    "KUBEDOG",
			EnvVars:     EnvVars("KUBEDOG_TRACK_ALL"),
			Destination: &dog.TrackGeneric,
		},
	}
}
