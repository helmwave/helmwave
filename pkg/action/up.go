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
	plandir        string
	kubedogEnabled bool
}

func (i *Up) Run() error {
	p := plan.New(i.plandir)
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
	i.dog = &kubedog.Config{}
	return &cli.Command{
		Name:  "up",
		Usage: "🚢 Apply your plan",
		Flags: []cli.Flag{
			flagPlandir(&i.plandir),
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
		},
		Action: toCtx(i.Run),
	}
}
