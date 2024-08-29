package action

import "github.com/urfave/cli/v2"

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

// func flagLabels(v *cli.StringSlice) cli.Flag {
//	return &cli.StringSliceFlag{
//		Name:        "labels",
//		Aliases:     []string{"l"},
//		Usage:       "build releases by label: -l app=nginx",
//		Category:    "SELECTION",
//		EnvVars:     EnvVars("LABELS"),
//		Destination: v,
//	}
//}
