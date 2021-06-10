package action

import (
	"github.com/helmwave/helmwave/pkg/kubedog"
	"github.com/helmwave/helmwave/pkg/plan"
	"github.com/urfave/cli/v2"
	"time"
)

type Install struct {
	plandir  string
	parallel bool
	kubedog  bool
	dog      *kubedog.Config
}

func (i *Install) Run() error {
	p := plan.New(i.plandir)
	return p.Apply()
}

func (i *Install) Cmd() *cli.Command {
	i.dog = &kubedog.Config{}
	return &cli.Command{
		Name:    "install",
		Aliases: []string{"apply", "sync", "deploy"},
		Usage:   "Apply your plan",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "plandir",
				Value:       ".helmwave/",
				Usage:       "Path to plandir",
				EnvVars:     []string{"HELMWAVE_PLANDIR"},
				Destination: &i.plandir,
			},
			&cli.BoolFlag{
				Name:        "parallel",
				Usage:       "It allows you call `helm install` in parallel mode ",
				Value:       true,
				EnvVars:     []string{"HELMWAVE_PARALLEL"},
				Destination: &i.parallel,
			},

			&cli.BoolFlag{
				Name:        "kubedog",
				Usage:       "Enable/Disable kubedog",
				Value:       false,
				EnvVars:     []string{"HELMWAVE_KUBEDOG"},
				Destination: &i.kubedog,
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
		},
		Action: toCtx(i.Run),
	}
}
