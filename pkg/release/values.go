package release

import (
	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/template"
	log "github.com/sirupsen/logrus"
	"os"
)

type ValuesReference struct {
	Src   string
	Local string
}

func (v *ValuesReference) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshal(&v.Src)
}

// ValuesMap Todo: Parallel
func (rel *Config) ValuesMap(dir string) (Map map[string]string) {
	for i, v := range rel.Values {
		stat, err := os.Stat(v)
		local := err == nil && !stat.IsDir()

		dst := dir + ".values/" + string(rel.Uniq()) + "/" + string(rune(i)) + ".yml"

		var src string
		if local {
			src = v
		} else if helper.IsUrl(v) {
			err = helper.Download(dst, v)
			if err != nil {
				log.Warn(v, " skipping: ", err)
				continue
			}

			src = dst
		} else {
			log.Warn("bad values path: ", v)
			continue
		}

		err = template.Tpl2yml(src, dst, struct{ Release *Config }{rel})
		if err != nil {
			log.Error(err)
			continue
		}

		// Create symlink for facilities
		if local {
			symlink := dir + v
			err = os.Symlink(dst, symlink)
			if err != nil {
				log.Fatal(err)
			}

			Map[v] = symlink
		} else {
			Map[v] = dst
		}

	}

	return Map

}

func map2Slice(Map map[string]string, Slice []string) {
	i := 0
	for _, v := range Map {
		Slice[i] = v
		i++
	}
}
