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
		if _, err := os.Stat(v); err != nil {
			if os.IsNotExist(err) {
				rel.Values = append(rel.Values[:i], rel.Values[i+1:]...)
				continue
			} else {
				fmt.Println(err)
			}
		}

		p := v + ".plan"
		err := template.Tpl2yml(v, p, struct{ Release *Config }{rel}, debug)
		if err != nil {
			fmt.Println(err)
		}

		rel.Values[i] = p
	}

}
