package release

import (
	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/template"
	log "github.com/sirupsen/logrus"
	"os"
)

func (rel *Config) V(dir string) (err error) {
	tmp := os.TempDir()
	v := rel.VDownload(tmp)

	v, err = rel.VTemplate(v, dir)
	if err != nil {
		return err
	}

	rel.Values = v
	return nil
}

func (rel *Config) VDownload(dir string) (values []string) {
	for i, url := range rel.Values {
		stat, err := os.Stat(url)
		local := err == nil && !stat.IsDir()
		if local {
			values = append(values, url)
		} else {
			// Download
			file := dir + string(rune(i)) + ".yml"
			err := helper.Download(file, url)
			if err != nil {
				log.Warn(url, " skipping: ", err)
			} else {
				values = append(values, file)
			}
		}
	}

	return values
}

func (rel *Config) VTemplate(values []string, dir string) (vals []string, err error) {
	for _, path := range values {
		dst := dir + path + "." + string(rel.Uniq()) + ".plan.yml"
		err = template.Tpl2yml(path, dst, struct{ Release *Config }{rel})
		if err != nil {
			return nil, err
		}

		vals = append(vals, dst)
	}

	return vals, nil
}
