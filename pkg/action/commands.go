package action

import (
	"time"

	"github.com/helmwave/helmwave/pkg/kubedog"
	"github.com/urfave/cli/v2"
)

func (i *Build) Cmd() *cli.Command {
	return &cli.Command{
		Name:   "build",
		Usage:  "🏗 Build a plan",
		Flags:  i.flags(),
		Action: toCtx(i.Run),
	}
}

func (i *Build) flags() []cli.Flag {
	// Init sub-structures
	i.yml = &Yml{}
	i.diff = &Diff{}

	self := []cli.Flag{
		flagPlandir(&i.plandir),
		flagTags(&i.tags),
		flagMatchAllTags(&i.matchAll),
		flagDiffMode(&i.diffMode),

		&cli.BoolFlag{
			Name:        "yml",
			Usage:       "Auto helmwave.yml.tpl --> helmwave.yml",
			Value:       false,
			EnvVars:     []string{"HELMWAVE_AUTO_YML", "HELMWAVE_AUTO_YAML"},
			Destination: &i.autoYml,
		},
	}

	self = append(self, i.diff.flags()...)
	self = append(self, i.yml.flags()...)

	return self
}

func (d *Diff) Cmd() *cli.Command {
	plan := DiffLocalPlan{diff: d}
	live := DiffLive{diff: d}

	return &cli.Command{
		Name:    "diff",
		Usage:   "🆚 Show Differences",
		Aliases: []string{"vs"},
		Flags:   d.flags(),
		Subcommands: []*cli.Command{
			plan.Cmd(),
			live.Cmd(),
		},
	}
}

func (d *Diff) flags() []cli.Flag {
	return []cli.Flag{
		flagDiffWide(&d.Wide),
		flagDiffShowSecret(&d.ShowSecret),
	}
}

func (d *DiffLive) Cmd() *cli.Command {
	return &cli.Command{
		Name:  "live",
		Usage: "plan 🆚 live",
		Flags: []cli.Flag{
			flagPlandir(&d.plandir),
		},
		Action: toCtx(d.Run),
	}
}

func (d *DiffLocalPlan) Cmd() *cli.Command {
	return &cli.Command{
		Name:   "plan",
		Usage:  "plan1  🆚  plan2",
		Flags:  d.flags(),
		Action: toCtx(d.Run),
	}
}

func (d *DiffLocalPlan) flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "plandir1",
			Value:       ".helmwave/",
			Usage:       "Path to plandir1",
			EnvVars:     []string{"HELMWAVE_PLANDIR_1", "HELMWAVE_PLANDIR"},
			Destination: &d.plandir1,
		},
		&cli.StringFlag{
			Name:        "plandir2",
			Value:       ".helmwave/",
			Usage:       "Path to plandir2",
			EnvVars:     []string{"HELMWAVE_PLANDIR_2"},
			Destination: &d.plandir2,
		},
	}
}

func (i *Down) Cmd() *cli.Command {
	return &cli.Command{
		Name:  "down",
		Usage: "🔪 Delete all",
		Flags: []cli.Flag{
			flagPlandir(&i.plandir),
		},
		Action: toCtx(i.Run),
	}
}

func (l *List) Cmd() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"ls"},
		Usage:   "👀 List of deployed releases",
		Flags: []cli.Flag{
			flagPlandir(&l.plandir),
		},
		Action: toCtx(l.Run),
	}
}

func (i *Rollback) Cmd() *cli.Command {
	return &cli.Command{
		Name:  "rollback",
		Usage: "⏮  Rollback your plan",
		Flags: []cli.Flag{
			flagPlandir(&i.plandir),
		},
		Action: toCtx(i.Run),
	}
}

func (l *Status) Cmd() *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "👁️ Status of deployed releases",
		Flags: []cli.Flag{
			flagPlandir(&l.plandir),
		},
		Action: toCtx(l.Run),
	}
}

func (i *Up) Cmd() *cli.Command {
	return &cli.Command{
		Name:   "up",
		Usage:  "🚢 Apply your plan",
		Flags:  i.flags(),
		Action: toCtx(i.Run),
	}
}

func (i *Up) flags() []cli.Flag {
	// Init sub-structures
	i.dog = &kubedog.Config{}
	i.build = &Build{}

	self := []cli.Flag{
		&cli.BoolFlag{
			Name:        "build",
			Usage:       "auto build",
			Value:       false,
			EnvVars:     []string{"HELMWAVE_AUTO_BUILD"},
			Destination: &i.autoBuild,
		},
		&cli.BoolFlag{
			Name:        "kubedog",
			Usage:       "Enable/Disable kubedog",
			Value:       false,
			EnvVars:     []string{"HELMWAVE_KUBEDOG_ENABLED", "HELMWAVE_KUBEDOG"},
			Destination: &i.kubedogEnabled,
		},
		&cli.DurationFlag{
			Name:        "kubedog-status-interval",
			Usage:       "Interval of kubedog status messages",
			Value:       5 * time.Second,
			EnvVars:     []string{"HELMWAVE_KUBEDOG_STATUS_INTERVAL"},
			Destination: &i.dog.StatusInterval,
		},
		&cli.DurationFlag{
			Name:        "kubedog-start-delay",
			Usage:       "Delay kubedog start, don't make it too late",
			Value:       time.Second,
			EnvVars:     []string{"HELMWAVE_KUBEDOG_START_DELAY"},
			Destination: &i.dog.StartDelay,
		},
		&cli.DurationFlag{
			Name:        "kubedog-timeout",
			Usage:       "Timout of kubedog multitrackers",
			Value:       5 * time.Minute,
			EnvVars:     []string{"HELMWAVE_KUBEDOG_TIMEOUT"},
			Destination: &i.dog.Timeout,
		},
	}

	return append(self, i.build.flags()...)
}

func (l *Validate) Cmd() *cli.Command {
	return &cli.Command{
		Name:  "validate",
		Usage: "🛂 Validate your plan",
		Flags: []cli.Flag{
			flagPlandir(&l.plandir),
		},
		Action: toCtx(l.Run),
	}
}

func (i *Yml) Cmd() *cli.Command {
	return &cli.Command{
		Name:   "yml",
		Usage:  "📄 Render helmwave.yml.tpl -> helmwave.yml",
		Flags:  i.flags(),
		Action: toCtx(i.Run),
	}
}

func (i *Yml) flags() []cli.Flag {
	return []cli.Flag{
		flagTplFile(&i.tpl),
		flagYmlFile(&i.file),
	}
}
