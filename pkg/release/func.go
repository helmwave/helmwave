package release

import (
	"fmt"
	"github.com/zhilyaev/helmwave/pkg/template"
	"sort"
)

func contains(t string, a []string) bool {
	i := sort.SearchStrings(a, t)
	return i < len(a) && a[i] == t
}

func RemoveIndex(s []Config, index int) []Config {
	return append(s[:index], s[index+1:]...)
}

func (rel *Config) In(a []Config) bool {
	for _, r := range a {
		if rel == &r {
			return true
		}
	}
	return false
}

func (rel *Config) RenderValues(debug bool) {
	for i, v := range rel.Values {
		p := v + ".plan"
		err := template.Tpl2yml(v, p, debug)
		if err != nil {
			fmt.Println(err)
		}

		rel.Values[i] = p
	}

}
