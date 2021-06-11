package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"sort"
	"strings"
)

type Build struct {
	plandir  string
	tags     cli.StringSlice
	matchAll bool
}

func (i *Build) Run() error {
	//tags := normalizeTags(i.tags)
	newPlan := plan.New(i.plandir)
	newPlan.Build()

	oldPlan := plan.New(i.plandir)
	if oldPlan.IsExist() {
		if err := oldPlan.Import(); err != nil {
			return err
		}

		// Show difference
		changelog, err := newPlan.Diff(oldPlan)
		log.Info(changelog)
		if err != nil {
			return err
		}
	}

	return newPlan.Export()
}

func (i *Build) Cmd() *cli.Command {
	return &cli.Command{
		Name:  "build",
		Usage: "Build a plandir",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "plandir",
				Value:       ".helmwave/",
				Usage:       "Path to plandir",
				EnvVars:     []string{"HELMWAVE_PLANDIR"},
				Destination: &i.plandir,
			},
			&cli.StringSliceFlag{
				Name:        "tags",
				Aliases:     []string{"t"},
				Usage:       "It allows you choose releases for sync. Example: -t tag1 -t tag3,tag4",
				EnvVars:     []string{"HELMWAVE_TAGS"},
				Destination: &i.tags,
			},
			&cli.BoolFlag{
				Name:        "match-all-tags",
				Usage:       "Match all provided tags",
				Value:       false,
				EnvVars:     []string{"HELMWAVE_MATCH_ALL_TAGS"},
				Destination: &i.matchAll,
			},
		},
		Action: toCtx(i.Run),
	}
}

// normalizeTags normalizes and splits comma-separated tag list.
// ["c", " b ", "a "] -> ["a", "b", "c"].
func normalizeTags(tags []string) []string {
	m := make([]string, len(tags))
	for i, t := range tags {
		m[i] = strings.TrimSpace(t)
	}
	sort.Strings(m)

	return m
}
