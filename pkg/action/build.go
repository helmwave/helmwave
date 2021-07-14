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
	diffWide     int
}

func (i *Build) Run() error {
	newPlan := plan.New(i.plandir)
	err := newPlan.Build(i.yml, i.normalizeTags(), i.matchAll)
	if err != nil {
		return err
	}

	// Show current plan
	newPlan.PrettyPlan()

	oldPlan := plan.New(i.plandir)
	if oldPlan.IsExist() {
		log.Info("🆚 diff with previous plan")
		if err = oldPlan.Import(); err != nil {
			return err
		}

		// Показать Diff
		newPlan.Diff(oldPlan, i.diffWide)
	}

	err = newPlan.Export()
	if err != nil {
		return err
	}

	log.WithField(
		"deploy it with next command",
		"helmwave deploy",
	).Info("🏗 Planfile is ready!")

	return nil
}

func (i *Build) Cmd() *cli.Command {
	return &cli.Command{
		Name:    "build",
		Usage:   "🏗 Build a plan",
		Aliases: []string{"plan"},
		Flags: []cli.Flag{
			flagPlandir(&i.plandir),
			flagTags(&i.tags),
			flagMatchAllTags(&i.matchAll),
			flagFile(&i.yml),
			flagDiffWide(&i.diffWide),
		},
		Action: toCtx(i.Run),
	}
}

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
