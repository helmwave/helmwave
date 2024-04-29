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
	options       plan.BuildOptions
	remoteSource  string
	plandir       string
	diffMode      string
	tags          cli.StringSlice
	autoYml       bool
	skipUnchanged bool
}

// Run is the main function for 'build' CLI command.
//
//nolint:funlen,gocognit,cyclop
func (i *Build) Run(ctx context.Context) (err error) {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	if i.remoteSource != "" {
		remoteSource, err := url.Parse(i.remoteSource)
		if err != nil {
			return fmt.Errorf("failed to parse remote source: %w", err)
		}

		downloadPath := cache.DefaultConfig.GetRemoteSourcePath(remoteSource)
		err = getter.Get(
			downloadPath,
			i.remoteSource,
			getter.WithContext(ctx),
			getter.WithDetectors(getter.Detectors),
			getter.WithGetters(getter.Getters),
			getter.WithDecompressors(getter.Decompressors),
		)
		if err != nil {
			return fmt.Errorf("failed to download remote source: %w", err)
		}

		// we need to move plandir back to where it should be
		defer func() {
			srcPlandir := filepath.Join(downloadPath, i.plandir)
			destPlandir := filepath.Join(wd, i.plandir)
			err := os.RemoveAll(destPlandir)
			if err != nil {
				log.WithError(err).Error("failed to clean plandir")
			}
			err = helper.MoveFile(srcPlandir, destPlandir)
			if err != nil {
				log.WithError(err).Error("failed to move plandir")
			}
		}()

		err = os.Chdir(downloadPath)
		if err != nil {
			return fmt.Errorf("failed to chdir to downloaded remote source: %w", err)
		}
	}

	if i.autoYml {
		err = i.yml.Run(ctx)
		if err != nil {
			return err
		}
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

			newPlan.DiffPlan(oldPlan, i.diff.Options)
		}

	case DiffModeLive:
		log.Info("🆚 Diff manifests in the kubernetes cluster")
		newPlan.DiffLive(ctx, i.diff.Options, i.diff.ThreeWayMerge)
	case DiffModeNone:
		log.Info("🆚 Skip diffing")
	default:
		log.Warnf("🆚❔Unknown %q diff mode, skipping", i.diffMode)
	}

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

		&cli.BoolFlag{
			Name:        "yml",
			Usage:       "auto helmwave.yml.tpl --> helmwave.yml",
			Value:       false,
			Category:    "YML",
			EnvVars:     []string{"HELMWAVE_AUTO_YML", "HELMWAVE_AUTO_YAML"},
			Destination: &i.autoYml,
		},
		&cli.StringFlag{
			Name:        "remote-source",
			Usage:       "go-getter URL to download build sources",
			Value:       "",
			Category:    "BUILD",
			EnvVars:     []string{"HELMWAVE_REMOTE_SOURCE"},
			Destination: &i.remoteSource,
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
