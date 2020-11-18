package release

import (
	"fmt"
	"github.com/zhilyaev/helmwave/pkg/template"
	"os"
	"sort"
)

func contains(t string, a []string) bool {
	i := sort.SearchStrings(a, t)
	return i < len(a) && a[i] == t
}

func (rel *Config) In(a []Config) bool {
	for _, r := range a {
		if rel == &r {
			return true
		}
	}
	return false
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

func (rel *Config) RenderValues(debug bool) {
	rel.PlanValues()

	for i, v := range rel.Values {
		p := v + ".plan"
		err := template.Tpl2yml(v, p, struct{ Release *Config }{rel}, debug)
		if err != nil {
			fmt.Println(err)
		}

		rel.Values[i] = p
	}

}
