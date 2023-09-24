package action

import (
	"context"
	"io/fs"
	"net/url"
	"sort"
	"strings"

	"github.com/helmwave/go-fsimpl/filefs"
	"github.com/helmwave/helmwave/pkg/cache"
	"github.com/helmwave/helmwave/pkg/plan"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var _ Action = (*Build)(nil)

// Build is a struct for running 'build' CLI command.
type Build struct {
	yml            *Yml
	diff           *Diff
	options        plan.BuildOptions
	planFS         fs.StatFS
	diffMode       string
	chartsCacheDir string
	tags           cli.StringSlice
	autoYml        bool
	skipUnchanged  bool
}

// Run is the main function for 'build' CLI command.
func (i *Build) Run(ctx context.Context) error {
	if i.autoYml {
		err := i.yml.Run(ctx)
		if err != nil {
			return err
		}
	}

	err := cache.ChartsCache.Init(i.planFS, i.chartsCacheDir)
	if err != nil {
		return err
	}

	newPlan := plan.New()

	i.options.Tags = i.normalizeTags()
	i.options.Yml = i.yml.destFS
	i.options.Templater = i.yml.templater

	workdirFS, err := filefs.New(&url.URL{Scheme: "file", Path: "."})
	if err != nil {
		return err
	}

	err = newPlan.Build(ctx, workdirFS, i.planFS, i.options)
	if err != nil {
		return err
	}

	// Show current plan
	newPlan.Logger().Info("🏗 Plan")

	switch i.diffMode {
	case DiffModeLocal:
		oldPlan := plan.New()
		if oldPlan.IsExist(i.planFS) {
			log.Info("🆚 Diff with previous local plan")
			if err := oldPlan.Import(ctx, i.planFS); err != nil {
				return err
			}

			newPlan.DiffPlan(oldPlan, i.diff.ShowSecret, i.diff.Wide)
		}

	case DiffModeLive:
		log.Info("🆚 Diff manifests in the kubernetes cluster")
		newPlan.DiffLive(ctx, i.planFS, i.diff.ShowSecret, i.diff.Wide, i.diff.ThreeWayMerge)
	case DiffModeNone:
		log.Info("🆚 Skip diffing")
	default:
		log.Warnf("🆚❔Unknown %q diff mode, skipping", i.diffMode)
	}

	err = newPlan.Export(ctx, i.planFS, i.skipUnchanged)
	if err != nil {
		return err
	}

	log.Info("🏗 Planfile is ready!")

	return nil
}

// Cmd returns 'build' *cli.Command.
func (i *Build) Cmd() *cli.Command {
	return &cli.Command{
		Name:   "build",
		Usage:  "🏗 build a plan",
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
		flagTags(&i.tags),
		flagMatchAllTags(&i.options.MatchAll),
		flagGraphWidth(&i.options.GraphWidth),
		flagSkipUnchanged(&i.skipUnchanged),
		flagDiffMode(&i.diffMode),
		flagChartsCacheDir(&i.chartsCacheDir),
		flagPlandirLocation(&i.planFS),

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
