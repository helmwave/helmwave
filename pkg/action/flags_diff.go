package action

import "github.com/urfave/cli/v2"

// flagDiffMode pass val to urfave flag.
func flagDiffMode(v *string) cli.Flag {
	return &cli.StringFlag{
		Name:        "diff-mode",
		Value:       "live",
		Category:    "DIFF",
		Usage:       "you can set: [ live | local | none ]. Set `offline_kube_version` to use [ local | none ]",
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
