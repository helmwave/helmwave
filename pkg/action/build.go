package action

import (
	"context"
	"sort"
	"strings"

	"github.com/helmwave/helmwave/pkg/plan"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// Build is struct for running 'build' CLI command.
type Build struct {
	yml      *Yml
	diff     *Diff
	plandir  string
	diffMode string
	tags     cli.StringSlice
	matchAll bool
	autoYml  bool

	// diffLive *DiffLive
	// diffLocal *DiffLocalPlan
}

const (
	// DiffModeLive is a subcommand name for diffing manifests in plan with actually running manifests in k8s.
	DiffModeLive = "live"

	// DiffModeLocal is a subcommand name for diffing manifests in two plans.
	DiffModeLocal = "local"
)

// Run is main function for 'build' CLI command.
func (i *Build) Run() (err error) {
	if i.autoYml {
		err = i.yml.Run()
		if err != nil {
			return err
		}
	}

	newPlan := plan.New(i.plandir)
	ctx := context.TODO()
	err = newPlan.Build(ctx, i.yml.file, i.normalizeTags(), i.matchAll, i.yml.templater)
	if err != nil {
		return err
	}

	// Show current plan
	newPlan.Logger().Info("🏗 Plan")

	switch i.diffMode {
	case DiffModeLocal:
		oldPlan := plan.New(i.plandir)
		if oldPlan.IsExist() {
			log.Info("🆚 Diff with previous local plan")
			if err := oldPlan.Import(); err != nil {
				return err
			}

			newPlan.DiffPlan(oldPlan, i.diff.ShowSecret, i.diff.Wide)
		}

	case DiffModeLive:
		log.Info("🆚 Diff manifests in the kubernetes cluster")
		newPlan.DiffLive(ctx, i.diff.ShowSecret, i.diff.Wide)
	default:
		log.Warnf("I dont know what is %q. I am skiping diff.", i.diffMode)
	}

	err = newPlan.Export()
	if err != nil {
		return err
	}

	log.WithField(
		"deploy it with next command",
		"helmwave up --plandir "+i.plandir,
	).Info("🏗 Planfile is ready!")

	return nil
}

// Cmd returns 'build' *cli.Command.
func (i *Build) Cmd() *cli.Command {
	return &cli.Command{
		Name:   "build",
		Usage:  "🏗 Build a plan",
		Flags:  i.flags(),
		Action: toCtx(i.Run),
	}
}

// flags return flag set of CLI urfave.
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

// normalizeTags is wrapper for normalizeTagList.
func (i *Build) normalizeTags() []string {
	return normalizeTagList(i.tags.Value())
}

// normalizeTags normalizes and splits comma-separated tag list.
// ["c", " b ", "a "] -> ["a", "b", "c"].
func normalizeTagList(tags []string) []string {
	m := make([]string, len(tags))
	for i, t := range tags {
		m[i] = strings.TrimSpace(t)
	}
	sort.Strings(m)

	return m
}
