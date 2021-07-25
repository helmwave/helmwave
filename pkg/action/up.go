package action

import (
	"time"

	"github.com/helmwave/helmwave/pkg/kubedog"
	"github.com/helmwave/helmwave/pkg/plan"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type Up struct {
	dog            *kubedog.Config
	kubedogEnabled bool

	autoBuild bool
	build     *Build
}

func (i *Up) Run() error {
	if i.autoBuild {
		if err := i.build.Run(); err != nil {
			return err
		}
	}

	p := plan.New(i.build.plandir)
	if err := p.Import(); err != nil {
		return err
	}

	p.PrettyPlan()

	if i.kubedogEnabled {
		log.Warn("🐶 kubedog is enable")
		return p.ApplyWithKubedog(i.dog)
	}

	return p.Apply()
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
