package release

import (
	"github.com/zhilyaev/helmwave/pkg/helper"
	"os"
	"sort"
	"strings"
)

func Plan(tags []string, releases []Config) (plan []Config) {
	if len(tags) == 0 {
		return releases
	}

	for _, t := range tags {
		// "c, b , a " -> "c,b,a"
		t := strings.Replace(t, " ", "", -1)
		// "c,b,a" -> ["c", "b", "a"]
		m := strings.Split(t, ",")

		for _, r := range releases {
			sort.Strings(r.Tags)
			if len(m) > 1 {
				// ["c", "b", "a"] -> ["a", "b", "c"]
				sort.Strings(m)
				sort.Strings(r.Tags)

				// ["a", "b", "c"] -> "a,b,c"
				s1 := strings.Join(m, ",")
				s2 := strings.Join(r.Tags, ",")

				// "myTag,myTag2" == "myTag,myTag2"
				if s1 == s2 && !r.In(plan) {
					plan = append(plan, r)
				}

			} else {
				// if myTag in [myTag2, myTag, myTag1]
				if helper.Contains(t, r.Tags) && !r.In(plan) {
					plan = append(plan, r)
				}
			}
		}
	}

	return plan
}

func (rel *Config) PlanValues() {

	for i := len(rel.Values) - 1; i >= 0; i-- {
		if _, err := os.Stat(rel.Values[i]); err != nil {
			if os.IsNotExist(err) {
				rel.Values = append(rel.Values[:i], rel.Values[i+1:]...)
			}
		}
	}

}
