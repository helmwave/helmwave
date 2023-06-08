package action

import (
	"context"
	"sort"
	"strings"

	"github.com/helmwave/helmwave/pkg/cache"
	"github.com/helmwave/helmwave/pkg/plan"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

// Build is struct for running 'build' CLI command.
type Build struct {
	yml            *Yml
	diff           *Diff
	options        plan.BuildOptions
	plandir        string
	diffMode       string
	chartsCacheDir string
	tags           cli.StringSlice
	autoYml        bool
	skipUnchanged  bool
}

// Run is main function for 'build' CLI command.
func (i *Build) Run(ctx context.Context) (err error) {
	if i.autoYml {
		err = i.yml.Run(ctx)
		if err != nil {
			return err
		}
	}

	err = cache.ChartsCache.Init(i.chartsCacheDir)
	if err != nil {
		return err
	}

	newPlan := plan.New(i.plandir)

	i.options.Tags = i.normalizeTags()
	i.options.Yml = i.yml.file
	i.options.Templater = i.yml.templater

	err = newPlan.Build(ctx, i.options)
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
			if err := oldPlan.Import(ctx); err != nil {
				return err
			}

			newPlan.DiffPlan(oldPlan, i.diff.ShowSecret, i.diff.Wide)
		}

	case DiffModeLive:
		log.Info("🆚 Diff manifests in the kubernetes cluster")
		newPlan.DiffLive(ctx, i.diff.ShowSecret, i.diff.Wide, i.diff.ThreeWayMerge)
	case DiffModeNone:
		log.Info("🆚 Skip diffing")
	default:
		log.Warnf("I don't know what is %q diff mode. I am skiping diff.", i.diffMode)
	}

	err = newPlan.Export(ctx, i.skipUnchanged)
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
		flagMatchAllTags(&i.options.MatchAll),
		flagGraphWidth(&i.options.GraphWidth),
		flagSkipUnchanged(&i.skipUnchanged),
		flagDiffMode(&i.diffMode),
		flagChartsCacheDir(&i.chartsCacheDir),

		&cli.BoolFlag{
			Name:        "yml",
			Usage:       "auto helmwave.yml.tpl --> helmwave.yml",
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
