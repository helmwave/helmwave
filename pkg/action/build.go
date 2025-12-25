package action

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/go-getter"
	"github.com/helmwave/helmwave/pkg/cache"
	"github.com/helmwave/helmwave/pkg/clictx"
	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/plan"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var _ Action = (*Build)(nil)

// Build is a struct for running 'build' CLI command.
type Build struct {
	yml           *Yml
	diff          *Diff
	remoteSource  string
	plandir       string
	diffMode      string
	tags          cli.StringSlice
	options       plan.BuildOptions
	autoYml       bool
	skipUnchanged bool
}

// Run is the main function for 'build' CLI command.
func (i *Build) Run(ctx context.Context) (err error) {
	cliCtx := clictx.GetCLIFromContext(ctx)
	if cliCtx == nil {
		return fmt.Errorf("failed to get CLI context")
	}

	// Download Remote source
	if i.remoteSource != "" {
		wd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current directory: %w", err)
		}
		log.Tracef("Work dir: %s", wd)

		downloadPath, err := i.downloadRemoteSrc(ctx)
		if err != nil {
			return err
		}

		defer i.cleanRemoteSrc(downloadPath, wd)
	}

	if i.autoYml {
		if helper.IsExists(i.yml.tpl) {
			err = i.yml.Run(ctx)
			if err != nil {
				return err
			}
		} else {
			log.Warnf("You've enabled auto yml, but I can't find %q. I'm skipping yml phase.", i.yml.tpl)
		}
	}

	newPlan := plan.New(i.plandir)

	i.options.Tags = i.normalizeTags()
	i.options.Yml = i.yml.file
	if i.options.Templater == "" {
		i.options.Templater = cliCtx.String("templater")
	}

	err = newPlan.Build(ctx, i.options)
	if err != nil {
		return err
	}

	// Show current plan
	newPlan.Logger().Info("🏗 Plan")

	// Diff
	err = i.diffing(ctx, newPlan)
	if err != nil {
		return err
	}

	// Save new plan
	err = newPlan.Export(ctx, i.skipUnchanged)
	if err != nil {
		return err
	}

	log.Info("🏗 Planfile is ready!")

	return nil
}

// Cmd returns 'build' *cli.Command.
func (i *Build) Cmd() *cli.Command {
	return &cli.Command{
		Name:     "build",
		Category: Step1,
		Usage:    "🏗 build a plan",
		Flags:    i.flags(),
		Before: func(q *cli.Context) error {
			i.diff.FixFields()

			return nil
		},
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
		flagYmlTemplateEngine(&i.yml.templater),
		flagBuildTemplateEngine(&i.options.Templater),

		&cli.BoolFlag{
			Name:        "yml",
			Usage:       "auto helmwave.yml.tpl --> helmwave.yml",
			Value:       false,
			Category:    "YML",
			EnvVars:     EnvVars("AUTO_YML", "AUTO_YAML"),
			Destination: &i.autoYml,
		},
		&cli.StringFlag{
			Name:        "remote-source",
			Usage:       "go-getter URL to download build sources",
			Value:       "",
			Category:    "BUILD",
			EnvVars:     EnvVars("REMOTE_SOURCE"),
			Destination: &i.remoteSource,
		},
		&cli.BoolFlag{
			Name:        "dependencies",
			Usage:       "evaluate releases dependencies and add them to the plan even if they don't match provided tags",
			Value:       true,
			Category:    Step1,
			EnvVars:     EnvVars("DEPENDENCIES_ENABLED", "DEPENDENCIES"),
			Destination: &i.options.EnableDependencies,
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
	m := helper.SlicesMap(tags, strings.TrimSpace)
	sort.Strings(m)

	return m
}

func (i *Build) cleanRemoteSrc(downloadPath, wd string) {
	srcPlandir := filepath.Join(downloadPath, i.plandir)
	destPlandir := filepath.Join(wd, i.plandir)

	err := helper.MoveFile(srcPlandir, destPlandir)
	if err != nil {
		log.WithError(err).Error("failed to move plandir")
	}
}

func (i *Build) downloadRemoteSrc(ctx context.Context) (path string, err error) {
	src, err := url.Parse(i.remoteSource)
	if err != nil {
		return "", fmt.Errorf("failed to parse remote source: %w", err)
	}

	path = cache.Default.GetRemoteSourcePath(src)
	err = getter.Get(
		path,
		i.remoteSource,
		getter.WithContext(ctx),
		getter.WithDetectors(getter.Detectors),
		getter.WithGetters(getter.Getters),
		getter.WithDecompressors(getter.Decompressors),
	)
	if err != nil {
		return "", fmt.Errorf("failed to download remote source: %w", err)
	}

	err = os.Chdir(path)
	if err != nil {
		return "", fmt.Errorf("failed to chdir to downloaded remote source: %w", err)
	}

	return
}

func (i *Build) diffing(ctx context.Context, p *plan.Plan) error {
	switch i.diffMode {
	case DiffModeLocal:
		oldPlan := plan.New(i.plandir)
		if oldPlan.IsExist() {
			log.Info("🆚 Diff with previous local plan")
			if err := oldPlan.Import(ctx); err != nil {
				return err
			}

			p.DiffPlan(oldPlan, i.diff.Options)
		}
	case DiffModeLive:
		log.Info("🆚 Diff manifests in the kubernetes cluster")
		p.DiffLive(ctx, i.diff.Options, i.diff.ThreeWayMerge)
	case DiffModeNone:
		log.Info("🆚 Skip diffing")
	default:
		log.Warnf("🆚❔Unknown %q diff mode, skipping", i.diffMode)
	}

	return nil
}
