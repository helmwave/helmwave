package release

import (
	"sort"
	"strings"

	"github.com/helmwave/helmwave/pkg/helper"
	log "github.com/sirupsen/logrus"
)

func Plan(tags []string, releases []*Config, deps bool) (plan []*Config) {
	if len(tags) == 0 {
		return releases
	}

	m := normalizeTagList(tags)

	releasesMap := make(map[string]*Config)
	if deps {
		for _, r := range releases {
			releasesMap[r.UniqName()] = r
		}
	}

	for _, r := range releases {
		if checkTagInclusion(m, r.Tags, deps) {
			plan = addToPlan(plan, r, releasesMap, deps)
		}
	}

	return plan
}

func addToPlan(plan []*Config, release *Config, releases map[string]*Config, deps bool) []*Config {
	if release.In(plan) {
		return plan
	}

	r := append(plan, release)

	if deps {
		for _, depName := range release.DependsOn {
			if dep, ok := releases[depName]; ok {
				r = addToPlan(r, dep, releases, true)
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
func checkTagInclusion(targetTags []string, releaseTags []string, matchAll bool) bool {
	if len(targetTags) == 0 {
		return true
	}

	for _, t := range targetTags {
		contains := helper.Contains(t, releaseTags)
		if matchAll && !contains {
			return false
		}
		if !matchAll && contains {
			return true
		}
	}

	return matchAll
}

// filterValuesFiles filters non-existent values files.
func (rel *Config) filterValuesFiles() {
	for i := len(rel.Values) - 1; i >= 0; i-- {
		err := rel.Values[i].Download()
		if err != nil {
			log.Errorf("Failed to find %s, skipping: %v", rel.Values[i].GetPath(), err)
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
