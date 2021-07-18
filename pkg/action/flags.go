package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
)

func flagPlandir(v *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "plandir",
		Value:       plan.Dir,
		Usage:       "Path to plandir",
		EnvVars:     []string{"HELMWAVE_PLANDIR"},
		Destination: v,
	}
}

func flagTags(v *cli.StringSlice) *cli.StringSliceFlag {
	return &cli.StringSliceFlag{
		Name:        "tags",
		Aliases:     []string{"t"},
		Usage:       "It allows you choose releases for sync. Example: -t tag1 -t tag3,tag4",
		EnvVars:     []string{"HELMWAVE_TAGS"},
		Destination: v,
	}
}

func flagMatchAllTags(v *bool) *cli.BoolFlag {
	return &cli.BoolFlag{
		Name:        "match-all-tags",
		Usage:       "Match all provided tags",
		Value:       false,
		EnvVars:     []string{"HELMWAVE_MATCH_ALL_TAGS"},
		Destination: v,
	}
}

// func flagParallel(v *bool) *cli.BoolFlag {
// 	return &cli.BoolFlag{
// 		Name:        "parallel",
// 		Usage:       "It allows you call `helm` in parallel mode ",
// 		Value:       true,
// 		EnvVars:     []string{"HELMWAVE_PARALLEL"},
// 		Destination: v,
// 	}
// }

func flagFile(v *string) *cli.StringFlag {
	return &cli.StringFlag{
		Name:        "file",
		Aliases:     []string{"f"},
		Value:       plan.Body,
		Usage:       "Main yml file",
		EnvVars:     []string{"HELMWAVE_YAML"},
		Destination: v,
	}
}

func flagDiffWide(v *int) *cli.IntFlag {
	return &cli.IntFlag{
		Name:        "diff-wide",
		Value:       5,
		Usage:       "Show line around change",
		EnvVars:     []string{"HELMWAVE_DIFF_WIDE"},
		Destination: v,
	}
}
