package action

import (
	"github.com/helmwave/helmwave/pkg/plan"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"sort"
	"strings"
)

type Build struct {
	plandir, yml string
	tags         cli.StringSlice
	matchAll     bool
}

func (i *Build) Run() error {
	newPlan := plan.New(i.plandir)
	err := newPlan.Build(i.yml, i.normalizeTags(), i.matchAll)
	if err != nil {
		return err
	}

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
			flagPlandir(&i.plandir),
			flagTags(&i.tags),
			flagMatchAllTags(&i.matchAll),
			flagFile(&i.yml),
		},
		Action: toCtx(i.Run),
	}
}

// normalizeTags normalizes and splits comma-separated tag list.
// ["c", " b ", "a "] -> ["a", "b", "c"].
func (i *Build) normalizeTags() []string {
	m := make([]string, len(i.tags.Value()))
	for i, t := range i.tags.Value() {
		m[i] = strings.TrimSpace(t)
	}
	sort.Strings(m)

	return m
}
