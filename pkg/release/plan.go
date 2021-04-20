package release

import (
	"github.com/helmwave/helmwave/pkg/feature"
	"github.com/helmwave/helmwave/pkg/helper"
	log "github.com/sirupsen/logrus"
	"os"
	"sort"
	"strings"
)

func Plan(tags []string, releases []*Config) (plan []*Config) {
	if len(tags) == 0 {
		return releases
	}

	m := normalizeTagList(tags)

	releasesMap := make(map[string]*Config)
	if feature.Dependencies {
		for _, r := range releases {
			releasesMap[r.UniqName()] = r
		}
	}

	for _, r := range releases {
		if checkTagInclusion(m, r.Tags) {
			plan = addToPlan(plan, r, releasesMap)
		}
	}

	return plan
}

func addToPlan(plan []*Config, release *Config, releases map[string]*Config) []*Config {
	if release.In(plan) {
		return plan
	}

	r := append(plan, release)

	if feature.PlanDependencies {
		for _, depName := range release.DependsOn {
			if dep, ok := releases[depName]; ok {
				r = addToPlan(r, dep, releases)
			} else {
				log.Warnf("cannot find dependency %s in available releases, skipping it", depName)
			}
		}
	}

	return r
}

// normalizeTagList normalizes and splits comma-separated tag list.
// ["c", " b ", "a "] -> ["a", "b", "c"].
func normalizeTagList(tags []string) []string {
	m := make([]string, len(tags))
	for i, t := range tags {
		m[i] = strings.TrimSpace(t)
	}
	sort.Strings(m)

	return m
}

// checkTagInclusion checks where any of release tags are included in target tags.
func checkTagInclusion(targetTags []string, releaseTags []string) bool {
	for _, t := range targetTags {
		if !helper.Contains(t, releaseTags) {
			return false
		}
	}

	return true
}

// filterValuesFiles filters non-existent values files.
func (rel *Config) filterValuesFiles() {
	for i := len(rel.Values) - 1; i >= 0; i-- {
		stat, err := os.Stat(rel.Values[i])
		if os.IsNotExist(err) || stat.IsDir() {
			rel.Values = append(rel.Values[:i], rel.Values[i+1:]...)
		}
	}
}

func PlanValues(releases []*Config, dir string) error {
	for i, rel := range releases {
		err := rel.RenderValues(dir)
		if err != nil {
			return err
		}

		releases[i].Values = rel.Values
		log.WithFields(log.Fields{
			"release":   rel.Name,
			"namespace": rel.Options.Namespace,
			"values":    releases[i].Values,
		}).Debug("🐞 Render Values")
	}

	return nil
}
