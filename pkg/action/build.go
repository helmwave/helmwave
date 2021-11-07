package action

import (
	"sort"
	"strings"

	"github.com/helmwave/helmwave/pkg/plan"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type Build struct {
	yml            *Yml
	plandir        string
	tags           cli.StringSlice
	matchAll       bool
	autoYml        bool
	diffWide       int
	diffShowSecret bool
}

func (i *Build) Run() error {
	if i.autoYml {
		if err := i.yml.Run(); err != nil {
			return err
		}
	}

	newPlan := plan.New(i.plandir)
	err := newPlan.Build(i.yml.file, i.normalizeTags(), i.matchAll)
	if err != nil {
		return err
	}

	// Show current plan
	newPlan.PrettyPlan()

	oldPlan := plan.New(i.plandir)
	if oldPlan.IsExist() {
		log.Info("🆚 Diff with previous plan")
		if err := oldPlan.Import(); err != nil {
			return err
		}

		// Diff
		newPlan.Diff(oldPlan, i.diffWide, i.diffShowSecret)
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

func (i *Build) Cmd() *cli.Command {
	return &cli.Command{
		Name:   "build",
		Usage:  "🏗 Build a plan",
		Flags:  i.flags(),
		Action: toCtx(i.Run),
	}
}

func (i *Build) flags() []cli.Flag {
	// Init sub-structures
	i.yml = &Yml{}

	self := []cli.Flag{
		flagPlandir(&i.plandir),
		flagTags(&i.tags),
		flagMatchAllTags(&i.matchAll),
		flagDiffWide(&i.diffWide),
		flagDiffShowSecret(&i.diffShowSecret),

		&cli.BoolFlag{
			Name:        "yml",
			Usage:       "Auto helmwave.yml.tpl --> helmwave.yml",
			Value:       false,
			EnvVars:     []string{"HELMWAVE_AUTO_YML", "HELMWAVE_AUTO_YAML"},
			Destination: &i.autoYml,
		},
	}

	return append(self, i.yml.flags()...)
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
