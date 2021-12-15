package action

import (
	"sort"
	"strings"

	"github.com/helmwave/helmwave/pkg/plan"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

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

var (
	diffModeLive  = "live"
	diffModeLocal = "local"
)

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

	switch i.diffMode {
	case diffModeLocal:
		oldPlan := plan.New(i.plandir)
		if oldPlan.IsExist() {
			log.Info("🆚 Diff with previous local plan")
			if err := oldPlan.Import(); err != nil {
				return err
			}

			newPlan.DiffPlan(oldPlan, i.diff.ShowSecret, i.diff.Wide)
		}

	case diffModeLive:
		log.Info("🆚 Diff manifests in the kubernetes cluster")
		newPlan.DiffLive(i.diff.ShowSecret, i.diff.Wide)
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
